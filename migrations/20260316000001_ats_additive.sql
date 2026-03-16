-- =============================================================================
-- ATS Recruitment — Consolidated Additive Migration
-- Merges: 20260316000001_ats_m01_m05_additive.sql + 20260316000001_ats_m06_m15_batch2.sql
-- Covers: M01-M15 additive schema changes on top of initial migration
-- Conflict resolution: Batch2 schemas used for job_embeddings, blacklist_logs, email_templates
-- Generated: 20260316000001
-- ADDITIVE ONLY — no existing tables dropped or renamed.
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- PREREQUISITE EXTENSIONS
-- ---------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ---------------------------------------------------------------------------
-- ALTER TABLE: users — locale, theme, updated_at, temp_interviewer support
-- Source: Batch1 (locale, theme, updated_at) + Batch2 (role check, invited_by, is_temp)
-- ---------------------------------------------------------------------------
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_role_check;
ALTER TABLE users
    ADD CONSTRAINT users_role_check
        CHECK (role IN ('admin', 'hr', 'supervisor', 'temp_interviewer'));

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS locale           TEXT NOT NULL DEFAULT 'th'
                                              CHECK (locale IN ('th', 'en')),
    ADD COLUMN IF NOT EXISTS theme            TEXT NOT NULL DEFAULT 'light'
                                              CHECK (theme IN ('light', 'dark')),
    ADD COLUMN IF NOT EXISTS updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS invited_by       UUID REFERENCES users (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS invite_expires_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS is_temp          BOOLEAN NOT NULL DEFAULT false;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email    ON users(email);

-- ---------------------------------------------------------------------------
-- ALTER TABLE: jobs — expand status to include 'scheduled', add scheduler columns
-- Source: Batch1 (publish_at, close_at, updated_at) + Batch2 (status check)
-- ---------------------------------------------------------------------------
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;
ALTER TABLE jobs
    ADD CONSTRAINT jobs_status_check
        CHECK (status IN ('draft', 'scheduled', 'open', 'closed'));

ALTER TABLE jobs
    ADD COLUMN IF NOT EXISTS publish_at  TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS close_at    TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS updated_at  TIMESTAMPTZ NOT NULL DEFAULT now();

-- Index for scheduler goroutine polling (ADR-003)
CREATE INDEX IF NOT EXISTS idx_jobs_status_publish_at ON jobs(status, publish_at);
CREATE INDEX IF NOT EXISTS idx_jobs_status_close_at   ON jobs(status, close_at);

-- ---------------------------------------------------------------------------
-- ALTER TABLE: candidates — add is_starred, blacklist_reason, updated_at
-- Source: Batch1 (is_starred, updated_at) + Batch2 (blacklist_reason, starred index)
-- ---------------------------------------------------------------------------
ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS is_starred      BOOLEAN     NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS blacklist_reason TEXT,
    ADD COLUMN IF NOT EXISTS updated_at      TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE INDEX IF NOT EXISTS idx_candidates_starred ON candidates (is_starred) WHERE is_starred = true;

-- ---------------------------------------------------------------------------
-- ALTER TABLE: application_documents — add OCR status tracking + similarity score
-- Source: Batch1
-- ---------------------------------------------------------------------------
ALTER TABLE application_documents
    ADD COLUMN IF NOT EXISTS ocr_status       VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (ocr_status IN ('pending', 'processing', 'done', 'failed')),
    ADD COLUMN IF NOT EXISTS similarity_score FLOAT;

-- Index for OCR worker queue polling (ADR-004)
CREATE INDEX IF NOT EXISTS idx_app_docs_ocr_status ON application_documents(ocr_status);

-- ---------------------------------------------------------------------------
-- ALTER TABLE: applications — add unique constraint (duplicate application block, M03)
-- Source: Batch1
-- ---------------------------------------------------------------------------
CREATE UNIQUE INDEX IF NOT EXISTS idx_applications_candidate_job
    ON applications(candidate_id, job_id);

-- ===========================================================================
-- NEW TABLES
-- ===========================================================================

-- ---------------------------------------------------------------------------
-- position_marks — M06: HR marks up to 3 suitable positions per applicant
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS position_marks (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    job_id         UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    marked_by      UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id, job_id)
);

CREATE INDEX IF NOT EXISTS idx_position_marks_application ON position_marks (application_id);
CREATE INDEX IF NOT EXISTS idx_position_marks_job         ON position_marks (job_id);

-- ---------------------------------------------------------------------------
-- job_embeddings — M06/M14: JD + hard-skill vector per job (async post-publish)
-- RESOLUTION: Using Batch2 JSONB schema (more flexible, pgvector can be added later)
-- Source: Batch2 (replaces Batch1 pgvector version)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_embeddings (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id       UUID NOT NULL UNIQUE REFERENCES jobs (id) ON DELETE CASCADE,
    jd_vector    JSONB,                        -- float[] from FastAPI /embed
    skill_vector JSONB,                        -- aggregated skill vectors
    jd_weight    NUMERIC(5,4) NOT NULL DEFAULT 0.5000,
    skill_weight NUMERIC(5,4) NOT NULL DEFAULT 0.5000,
    computed_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT job_embeddings_weights_sum CHECK (
        ABS((jd_weight + skill_weight) - 1.0) < 0.0001
    )
);

-- ---------------------------------------------------------------------------
-- candidate_embeddings — M06/M10: resume text vector per candidate
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS candidate_embeddings (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id  UUID NOT NULL UNIQUE REFERENCES candidates (id) ON DELETE CASCADE,
    resume_vector JSONB,                       -- float[] from FastAPI /embed
    computed_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- blacklist_logs — M11: audit trail for blacklist/unblacklist events
-- RESOLUTION: Using Batch2 schema (candidate_id + action enum)
-- Source: Batch2 (replaces Batch1 email/phone version)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS blacklist_logs (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES candidates (id) ON DELETE CASCADE,
    action       TEXT NOT NULL CHECK (action IN ('blacklist', 'unblacklist')),
    reason       TEXT NOT NULL,
    performed_by UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_blacklist_logs_candidate ON blacklist_logs (candidate_id);

-- ---------------------------------------------------------------------------
-- pipeline_actions — M06/M07/M08: one row per HR action taken on a group
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS pipeline_actions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id       UUID NOT NULL REFERENCES jobs (id) ON DELETE RESTRICT,
    action_type  TEXT NOT NULL CHECK (action_type IN ('test', 'interview', 'hire', 'drop', 'transfer')),
    round_number INT NOT NULL DEFAULT 1,
    phase        TEXT NOT NULL DEFAULT 'phase1' CHECK (phase IN ('phase1', 'phase2')),
    status       TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'grading', 'closed')),
    created_by   UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pipeline_actions_job_id ON pipeline_actions (job_id);
CREATE INDEX IF NOT EXISTS idx_pipeline_actions_type   ON pipeline_actions (action_type);

-- ---------------------------------------------------------------------------
-- test_library — M07/M14: HR/Supervisor managed test bank
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_library (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title          TEXT NOT NULL,
    scope          TEXT NOT NULL CHECK (scope IN ('general', 'technical')),
    department_id  UUID REFERENCES departments (id) ON DELETE SET NULL,
    time_limit_min INT NOT NULL DEFAULT 60 CHECK (time_limit_min > 0),
    created_by     UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_test_library_scope ON test_library (scope);
CREATE INDEX IF NOT EXISTS idx_test_library_dept  ON test_library (department_id);

-- ---------------------------------------------------------------------------
-- test_questions — M07: questions belonging to a test library entry
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_id     UUID NOT NULL REFERENCES test_library (id) ON DELETE CASCADE,
    question    TEXT NOT NULL,
    type        TEXT NOT NULL CHECK (type IN ('mcq', 'essay', 'short_answer', 'true_false')),
    options     JSONB,
    score       NUMERIC(6,2) NOT NULL DEFAULT 1.00,
    order_index INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_test_questions_test_id ON test_questions (test_id);

-- ---------------------------------------------------------------------------
-- test_rounds — M07: one round instance per pipeline_action of type='test'
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_rounds (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pipeline_action_id UUID NOT NULL UNIQUE REFERENCES pipeline_actions (id) ON DELETE CASCADE,
    test_id            UUID NOT NULL REFERENCES test_library (id) ON DELETE RESTRICT,
    shuffle_seed       BIGINT,
    phase1_cutoff      INT,
    total_score        NUMERIC(6,2),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- test_assignments — M07: which applicants are assigned to which test round
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_assignments (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_round_id  UUID NOT NULL REFERENCES test_rounds (id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    status         TEXT NOT NULL DEFAULT 'invited'
                       CHECK (status IN ('invited', 'in_progress', 'submitted', 'graded')),
    invited_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (test_round_id, application_id)
);

CREATE INDEX IF NOT EXISTS idx_test_assignments_round ON test_assignments (test_round_id);
CREATE INDEX IF NOT EXISTS idx_test_assignments_app   ON test_assignments (application_id);

-- ---------------------------------------------------------------------------
-- test_sessions — M07: anti-cheat per applicant per test round
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_sessions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_assignment_id UUID NOT NULL UNIQUE REFERENCES test_assignments (id) ON DELETE CASCADE,
    started_at         TIMESTAMPTZ,
    submitted_at       TIMESTAMPTZ,
    tab_switch_count   INT NOT NULL DEFAULT 0,
    copy_paste_count   INT NOT NULL DEFAULT 0,
    last_heartbeat_at  TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- test_submissions — M07: submitted answers per applicant per test round
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_submissions (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_assignment_id UUID NOT NULL UNIQUE REFERENCES test_assignments (id) ON DELETE CASCADE,
    answer_data        JSONB NOT NULL DEFAULT '[]'::JSONB,
    submitted_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- test_grades — M07: per-question scores (MCQ auto, essay manual)
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_grades (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES test_submissions (id) ON DELETE CASCADE,
    question_id   UUID NOT NULL REFERENCES test_questions (id) ON DELETE CASCADE,
    score         NUMERIC(6,2) NOT NULL DEFAULT 0,
    is_correct    BOOLEAN,
    grader_note   TEXT,
    graded_by     UUID REFERENCES users (id) ON DELETE SET NULL,
    graded_at     TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (submission_id, question_id)
);

CREATE INDEX IF NOT EXISTS idx_test_grades_submission ON test_grades (submission_id);

-- ---------------------------------------------------------------------------
-- interview_rounds — M08: one round instance per pipeline_action of type='interview'
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS interview_rounds (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pipeline_action_id UUID NOT NULL UNIQUE REFERENCES pipeline_actions (id) ON DELETE CASCADE,
    phase1_cutoff      INT,
    status             TEXT NOT NULL DEFAULT 'open'
                           CHECK (status IN ('open', 'scoring', 'closed')),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- scoring_criteria — M08: evaluation criteria per interview round
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS scoring_criteria (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_round_id UUID NOT NULL REFERENCES interview_rounds (id) ON DELETE CASCADE,
    name               TEXT NOT NULL,
    description        TEXT,
    max_score          NUMERIC(6,2) NOT NULL DEFAULT 10.00,
    order_index        INT NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_scoring_criteria_round ON scoring_criteria (interview_round_id);

-- ---------------------------------------------------------------------------
-- interview_assignments — M08: applicants assigned to interview round
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS interview_assignments (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_round_id UUID NOT NULL REFERENCES interview_rounds (id) ON DELETE CASCADE,
    application_id     UUID NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    status             TEXT NOT NULL DEFAULT 'invited'
                           CHECK (status IN ('invited', 'scored', 'passed', 'dropped')),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (interview_round_id, application_id)
);

CREATE INDEX IF NOT EXISTS idx_interview_assignments_round ON interview_assignments (interview_round_id);
CREATE INDEX IF NOT EXISTS idx_interview_assignments_app   ON interview_assignments (application_id);

-- ---------------------------------------------------------------------------
-- interview_evaluations — M08: per-interviewer scores per criterion per applicant
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS interview_evaluations (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_assignment_id UUID NOT NULL REFERENCES interview_assignments (id) ON DELETE CASCADE,
    criterion_id            UUID NOT NULL REFERENCES scoring_criteria (id) ON DELETE CASCADE,
    interviewer_id          UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    score                   NUMERIC(6,2),
    note                    TEXT,
    submitted_at            TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (interview_assignment_id, criterion_id, interviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_interview_evals_assignment  ON interview_evaluations (interview_assignment_id);
CREATE INDEX IF NOT EXISTS idx_interview_evals_interviewer ON interview_evaluations (interviewer_id);

-- ---------------------------------------------------------------------------
-- invite_links — M08/M10: 256-bit token for temp interviewer access
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS invite_links (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    interview_round_id UUID NOT NULL REFERENCES interview_rounds (id) ON DELETE CASCADE,
    token              TEXT NOT NULL UNIQUE,
    invited_email      TEXT NOT NULL,
    invited_user_id    UUID REFERENCES users (id) ON DELETE SET NULL,
    status             TEXT NOT NULL DEFAULT 'active'
                           CHECK (status IN ('active', 'used', 'revoked', 'expired')),
    expires_at         TIMESTAMPTZ NOT NULL,
    revoked_by         UUID REFERENCES users (id) ON DELETE SET NULL,
    revoked_at         TIMESTAMPTZ,
    created_by         UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_invite_links_token  ON invite_links (token);
CREATE INDEX IF NOT EXISTS idx_invite_links_round  ON invite_links (interview_round_id);
CREATE INDEX IF NOT EXISTS idx_invite_links_email  ON invite_links (invited_email);
CREATE INDEX IF NOT EXISTS idx_invite_links_status ON invite_links (status) WHERE status = 'active';

-- ---------------------------------------------------------------------------
-- transfer_history — M09: transfer chain for an application
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS transfer_history (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    from_job_id    UUID NOT NULL REFERENCES jobs (id) ON DELETE RESTRICT,
    to_job_id      UUID NOT NULL REFERENCES jobs (id) ON DELETE RESTRICT,
    transferred_by UUID REFERENCES users (id) ON DELETE SET NULL,
    reason         TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_transfer_history_application ON transfer_history (application_id);
CREATE INDEX IF NOT EXISTS idx_transfer_history_from_job    ON transfer_history (from_job_id);
CREATE INDEX IF NOT EXISTS idx_transfer_history_to_job      ON transfer_history (to_job_id);

-- ---------------------------------------------------------------------------
-- offers — M12: hire offer records linked to an application
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS offers (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id   UUID NOT NULL UNIQUE REFERENCES applications (id) ON DELETE CASCADE,
    offer_letter_url TEXT,
    status           TEXT NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'accepted', 'declined', 'withdrawn')),
    offered_by       UUID REFERENCES users (id) ON DELETE SET NULL,
    offered_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_offers_application ON offers (application_id);
CREATE INDEX IF NOT EXISTS idx_offers_status      ON offers (status);

-- ---------------------------------------------------------------------------
-- requisitions — M12: supervisor-submitted requisitions for new headcount
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS requisitions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    department_id       UUID REFERENCES departments (id) ON DELETE SET NULL,
    employment_type_id  UUID REFERENCES employment_types (id) ON DELETE SET NULL,
    experience_level_id UUID REFERENCES experience_levels (id) ON DELETE SET NULL,
    description         TEXT,
    requirements        TEXT,
    slots               INT NOT NULL DEFAULT 1 CHECK (slots > 0),
    status              TEXT NOT NULL DEFAULT 'draft'
                            CHECK (status IN ('draft', 'submitted', 'approved', 'rejected', 'converted')),
    submitted_by        UUID REFERENCES users (id) ON DELETE SET NULL,
    reviewed_by         UUID REFERENCES users (id) ON DELETE SET NULL,
    review_note         TEXT,
    converted_job_id    UUID REFERENCES jobs (id) ON DELETE SET NULL,
    submitted_at        TIMESTAMPTZ,
    reviewed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_requisitions_status    ON requisitions (status);
CREATE INDEX IF NOT EXISTS idx_requisitions_dept      ON requisitions (department_id);
CREATE INDEX IF NOT EXISTS idx_requisitions_submitter ON requisitions (submitted_by);

-- ---------------------------------------------------------------------------
-- email_templates — M14: SYS (system-managed) + TPL (HR-editable) templates
-- RESOLUTION: Using Batch2 schema (code + type) instead of Batch1 (slug + locale)
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS email_templates (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code      TEXT NOT NULL UNIQUE,
    name      TEXT NOT NULL,
    type      TEXT NOT NULL CHECK (type IN ('system', 'custom')),
    subject   TEXT NOT NULL,
    body_html TEXT NOT NULL,
    variables JSONB NOT NULL DEFAULT '[]'::JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- jd_templates — M14: HR-managed job description templates
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS jd_templates (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title        TEXT NOT NULL,
    description  TEXT NOT NULL,
    requirements TEXT,
    created_by   UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- notifications — M13: in-app notification records per user
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    payload    JSONB NOT NULL DEFAULT '{}'::JSONB,
    is_read    BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_unread  ON notifications (user_id, is_read) WHERE is_read = false;
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications (created_at DESC);

-- ---------------------------------------------------------------------------
-- notification_preferences — M13: per-user per-event channel toggles
-- Source: Batch2
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS notification_preferences (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    in_app     BOOLEAN NOT NULL DEFAULT true,
    email      BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, event_type)
);

CREATE INDEX IF NOT EXISTS idx_notif_prefs_user ON notification_preferences (user_id);

-- ===========================================================================
-- SEED DATA
-- ===========================================================================

-- System email templates (SYS type — not editable by HR)
INSERT INTO email_templates (code, name, type, subject, body_html, variables) VALUES
    ('SYS_OTP',   'OTP Password Reset',        'system',
     'รหัส OTP สำหรับรีเซ็ตรหัสผ่าน',
     '<p>รหัส OTP ของคุณคือ: <strong>{{otp_code}}</strong> (หมดอายุใน 10 นาที)</p>',
     '["{{otp_code}}"]'::JSONB),
    ('TPL01',     'Application Received',       'custom',
     'ได้รับใบสมัครของคุณแล้ว — {{job_title}}',
     '<p>เรียน {{applicant_name}}, ขอบคุณสำหรับการสมัครตำแหน่ง {{job_title}}</p>',
     '["{{applicant_name}}", "{{job_title}}"]'::JSONB),
    ('TPL02',     'Application Status Update',  'custom',
     'อัปเดตสถานะใบสมัคร — {{job_title}}',
     '<p>เรียน {{applicant_name}}, สถานะใบสมัครของคุณได้รับการอัปเดต</p>',
     '["{{applicant_name}}", "{{job_title}}", "{{status}}"]'::JSONB),
    ('TPL03',     'Test Invitation',            'custom',
     'เชิญเข้าสอบออนไลน์ — {{job_title}}',
     '<p>เรียน {{applicant_name}}, คุณได้รับเชิญให้เข้าสอบออนไลน์ กรุณาเข้าสอบภายใน {{deadline}}</p>',
     '["{{applicant_name}}", "{{job_title}}", "{{test_link}}", "{{deadline}}"]'::JSONB),
    ('TPL04',     'Drop Notification',          'custom',
     'ผลการพิจารณาใบสมัคร — {{job_title}}',
     '<p>เรียน {{applicant_name}}, ขอขอบคุณสำหรับความสนใจ แต่เราไม่สามารถดำเนินการต่อได้ในขณะนี้</p>',
     '["{{applicant_name}}", "{{job_title}}"]'::JSONB),
    ('TPL05',     'Interview Invitation',       'custom',
     'เชิญเข้าสัมภาษณ์ — {{job_title}}',
     '<p>เรียน {{applicant_name}}, คุณได้รับเชิญสัมภาษณ์ในวันที่ {{interview_date}}</p>',
     '["{{applicant_name}}", "{{job_title}}", "{{interview_date}}", "{{interview_location}}"]'::JSONB),
    ('TPL06',     'Offer Letter',               'custom',
     'ข้อเสนองาน — {{job_title}}',
     '<p>เรียน {{applicant_name}}, เรามีความยินดีที่จะเสนองานในตำแหน่ง {{job_title}}</p>',
     '["{{applicant_name}}", "{{job_title}}", "{{offer_link}}"]'::JSONB)
ON CONFLICT (code) DO NOTHING;

COMMIT;

-- ===========================================================================
-- DOWN (for rollback reference — DO NOT run in production without review)
-- ===========================================================================
-- ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_check;
-- ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'hr', 'supervisor'));
-- ALTER TABLE users DROP COLUMN IF EXISTS locale, DROP COLUMN IF EXISTS theme, DROP COLUMN IF EXISTS updated_at;
-- ALTER TABLE users DROP COLUMN IF EXISTS invited_by, DROP COLUMN IF EXISTS invite_expires_at, DROP COLUMN IF EXISTS is_temp;
-- ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_status_check;
-- ALTER TABLE jobs ADD CONSTRAINT jobs_status_check CHECK (status IN ('draft', 'open', 'closed'));
-- ALTER TABLE jobs DROP COLUMN IF EXISTS publish_at, DROP COLUMN IF EXISTS close_at, DROP COLUMN IF EXISTS updated_at;
-- ALTER TABLE candidates DROP COLUMN IF EXISTS is_starred, DROP COLUMN IF EXISTS blacklist_reason, DROP COLUMN IF EXISTS updated_at;
-- ALTER TABLE application_documents DROP COLUMN IF EXISTS ocr_status, DROP COLUMN IF EXISTS similarity_score;
-- DROP TABLE IF EXISTS notification_preferences CASCADE;
-- DROP TABLE IF EXISTS notifications CASCADE;
-- DROP TABLE IF EXISTS jd_templates CASCADE;
-- DROP TABLE IF EXISTS email_templates CASCADE;
-- DROP TABLE IF EXISTS requisitions CASCADE;
-- DROP TABLE IF EXISTS offers CASCADE;
-- DROP TABLE IF EXISTS transfer_history CASCADE;
-- DROP TABLE IF EXISTS invite_links CASCADE;
-- DROP TABLE IF EXISTS interview_evaluations CASCADE;
-- DROP TABLE IF EXISTS interview_assignments CASCADE;
-- DROP TABLE IF EXISTS scoring_criteria CASCADE;
-- DROP TABLE IF EXISTS interview_rounds CASCADE;
-- DROP TABLE IF EXISTS test_grades CASCADE;
-- DROP TABLE IF EXISTS test_submissions CASCADE;
-- DROP TABLE IF EXISTS test_sessions CASCADE;
-- DROP TABLE IF EXISTS test_assignments CASCADE;
-- DROP TABLE IF EXISTS test_rounds CASCADE;
-- DROP TABLE IF EXISTS test_questions CASCADE;
-- DROP TABLE IF EXISTS test_library CASCADE;
-- DROP TABLE IF EXISTS pipeline_actions CASCADE;
-- DROP TABLE IF EXISTS blacklist_logs CASCADE;
-- DROP TABLE IF EXISTS candidate_embeddings CASCADE;
-- DROP TABLE IF EXISTS job_embeddings CASCADE;
-- DROP TABLE IF EXISTS position_marks CASCADE;
