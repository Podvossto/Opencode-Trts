-- =============================================================================
-- ATS Recruitment — Initial Schema Migration
-- Feature: M01 Auth + M02 User MGMT + M03 Portal + M04 Job Posting + M05 Form Builder
-- Generated: 20260316000000
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- departments — org chart hierarchy for HR/Supervisor/Admin users and jobs
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS departments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_departments_name ON departments (name);

-- ---------------------------------------------------------------------------
-- employment_types — Full-time, Part-time, Contract, etc.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS employment_types (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- experience_levels — Entry, Mid, Senior, Lead, etc.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS experience_levels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- hard_skills — skills dictionary for JD matching
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS hard_skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_hard_skills_name ON hard_skills (name);

-- ---------------------------------------------------------------------------
-- users — HR / Supervisor / Admin accounts
-- All new users created by Admin default to must_change_password = true (SRS FRA03)
-- Admin does NOT have OTP self-reset (SRS FRA07)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                TEXT NOT NULL UNIQUE,
    password_hash        TEXT NOT NULL,
    first_name           TEXT NOT NULL,
    last_name            TEXT NOT NULL,
    role                 TEXT NOT NULL CHECK (role IN ('admin', 'hr', 'supervisor')),
    department_id        UUID REFERENCES departments (id) ON DELETE SET NULL,
    is_active            BOOLEAN NOT NULL DEFAULT true,
    must_change_password BOOLEAN NOT NULL DEFAULT true,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email    ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_role     ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_dept     ON users (department_id);

-- ---------------------------------------------------------------------------
-- otp_logs — OTP records for HR/Supervisor password reset (SRS FRA05)
-- One valid OTP per user at a time; expires_at enforced at application layer
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS otp_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    otp_code    TEXT NOT NULL,
    is_used     BOOLEAN NOT NULL DEFAULT false,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_otp_logs_user_id ON otp_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_otp_logs_expires ON otp_logs (expires_at);

-- ---------------------------------------------------------------------------
-- application_forms — Form builder templates (schema stored as JSONB)
-- ADR-005: JSONB allows dynamic field definitions without schema migration per form change
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS application_forms (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    form_schema JSONB NOT NULL DEFAULT '[]'::JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- GIN index for full-text JSON queries on form_schema if needed
CREATE INDEX IF NOT EXISTS idx_application_forms_schema ON application_forms USING GIN (form_schema);

-- ---------------------------------------------------------------------------
-- jobs — Job posting lifecycle: draft → open → closed
-- publish_at / close_at drive the Go scheduler goroutine (ADR-004)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS jobs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    department_id       UUID NOT NULL REFERENCES departments (id) ON DELETE RESTRICT,
    employment_type_id  UUID NOT NULL REFERENCES employment_types (id) ON DELETE RESTRICT,
    experience_level_id UUID NOT NULL REFERENCES experience_levels (id) ON DELETE RESTRICT,
    form_id             UUID NOT NULL REFERENCES application_forms (id) ON DELETE RESTRICT,
    description         TEXT NOT NULL,
    requirements        TEXT,
    slots               INT NOT NULL DEFAULT 1 CHECK (slots > 0),
    status              TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'closed')),
    publish_at          TIMESTAMPTZ,
    close_at            TIMESTAMPTZ,
    created_by          UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_jobs_status      ON jobs (status);
CREATE INDEX IF NOT EXISTS idx_jobs_publish_at  ON jobs (publish_at) WHERE publish_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_jobs_close_at    ON jobs (close_at)   WHERE close_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_jobs_dept        ON jobs (department_id);

-- ---------------------------------------------------------------------------
-- job_hard_skills — many-to-many: job ↔ hard_skill with weight ratio 1-5
-- Weight is a ratio, not a percentage — does NOT need to sum to 100 (SRS §M04)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS job_hard_skills (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id     UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    skill_id   UUID NOT NULL REFERENCES hard_skills (id) ON DELETE CASCADE,
    weight     INT NOT NULL DEFAULT 1 CHECK (weight >= 1 AND weight <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (job_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_job_hard_skills_job_id ON job_hard_skills (job_id);

-- ---------------------------------------------------------------------------
-- candidates — Applicant identity records (portal-facing)
-- is_blacklisted enforced at submission time by email match (SRS §M03)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS candidates (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email          TEXT NOT NULL UNIQUE,
    first_name     TEXT NOT NULL,
    last_name      TEXT NOT NULL,
    phone          TEXT,
    is_blacklisted BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_candidates_email        ON candidates (email);
CREATE INDEX IF NOT EXISTS idx_candidates_blacklisted  ON candidates (is_blacklisted) WHERE is_blacklisted = true;

-- ---------------------------------------------------------------------------
-- applications — Application submissions
-- ocr_status tracks async OCR pipeline state (ADR-003)
-- UNIQUE (candidate_id, job_id) enforces no-duplicate-application rule (SRS §M03)
-- auto-close does NOT affect in-progress applications (SRS FRH10b)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS applications (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id     UUID NOT NULL REFERENCES candidates (id) ON DELETE CASCADE,
    job_id           UUID NOT NULL REFERENCES jobs (id) ON DELETE RESTRICT,
    status           TEXT NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending','in_review','testing','interviewing','hired','dropped','transferred')),
    form_data        JSONB NOT NULL DEFAULT '{}'::JSONB,
    ocr_status       TEXT NOT NULL DEFAULT 'pending' CHECK (ocr_status IN ('pending','done','failed')),
    similarity_score NUMERIC(5,4),               -- 0.0000 to 1.0000
    submitted_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (candidate_id, job_id)
);

CREATE INDEX IF NOT EXISTS idx_applications_job_id     ON applications (job_id);
CREATE INDEX IF NOT EXISTS idx_applications_candidate  ON applications (candidate_id);
CREATE INDEX IF NOT EXISTS idx_applications_status     ON applications (status);
CREATE INDEX IF NOT EXISTS idx_applications_ocr_status ON applications (ocr_status) WHERE ocr_status = 'pending';
CREATE INDEX IF NOT EXISTS idx_applications_submitted  ON applications (submitted_at DESC);

-- GIN index for form_data JSONB queries
CREATE INDEX IF NOT EXISTS idx_applications_form_data ON applications USING GIN (form_data);

-- ---------------------------------------------------------------------------
-- application_documents — Uploaded files linked to applications (resume, portfolio, etc.)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS application_documents (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id UUID NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    file_type      TEXT NOT NULL,           -- e.g. 'resume', 'portfolio', 'certificate'
    file_url       TEXT NOT NULL,           -- storage URL (S3, GCS, local)
    file_name      TEXT NOT NULL,
    file_size_kb   INT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_app_docs_application_id ON application_documents (application_id);

-- ---------------------------------------------------------------------------
-- test_autosave — Partial form autosave for portal applications
-- Stores intermediate state before final submission
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS test_autosave (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id         UUID NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    candidate_email TEXT NOT NULL,
    form_data      JSONB NOT NULL DEFAULT '{}'::JSONB,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (job_id, candidate_email)
);

CREATE INDEX IF NOT EXISTS idx_autosave_job_email ON test_autosave (job_id, candidate_email);

-- ---------------------------------------------------------------------------
-- Seed: reference data
-- ---------------------------------------------------------------------------

INSERT INTO employment_types (name) VALUES
    ('Full-time'),
    ('Part-time'),
    ('Contract'),
    ('Internship')
ON CONFLICT (name) DO NOTHING;

INSERT INTO experience_levels (name) VALUES
    ('Entry Level'),
    ('Mid Level'),
    ('Senior Level'),
    ('Lead / Principal'),
    ('Executive')
ON CONFLICT (name) DO NOTHING;

COMMIT;

-- Down migration (run manually to rollback):
-- DROP TABLE IF EXISTS test_autosave CASCADE;
-- DROP TABLE IF EXISTS application_documents CASCADE;
-- DROP TABLE IF EXISTS applications CASCADE;
-- DROP TABLE IF EXISTS candidates CASCADE;
-- DROP TABLE IF EXISTS job_hard_skills CASCADE;
-- DROP TABLE IF EXISTS jobs CASCADE;
-- DROP TABLE IF EXISTS application_forms CASCADE;
-- DROP TABLE IF EXISTS otp_logs CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS hard_skills CASCADE;
-- DROP TABLE IF EXISTS experience_levels CASCADE;
-- DROP TABLE IF EXISTS employment_types CASCADE;
-- DROP TABLE IF EXISTS departments CASCADE;
