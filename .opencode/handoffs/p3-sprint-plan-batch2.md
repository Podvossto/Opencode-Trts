# ATS Sprint Plan — Batch 2 (Sprints 4-7)

---

## SPRINT 4 — Resume Review + Test Management (M06, M07)

### Task List

```
TASK: ATS-042
  title       : Resume review API — 3-panel data + score calculation
  description : GET application detail (resume, OCR text, job info), POST score with weight popup. position_marks CRUD (max 3)
  ac          : US-M06-01, US-M06-02, US-M06-03
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-043
  title       : Resume review action buttons + similarity auto-compute
  description : POST pipeline action (pass/drop/transfer). On new job publish, compute similarity for starred candidates via ai-service
  ac          : US-M06-04, US-M06-05
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-044
  title       : Test library CRUD + test round creation
  description : GET/POST/PUT/DELETE test_library. POST test_rounds — create round, assign applicants, link to pipeline_action
  ac          : US-M07-06, US-M07-01, US-M07-02
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-045
  title       : Online test delivery with anti-cheat + grading
  description : Test session endpoint with heartbeat, tab_switch/copy_paste tracking. MCQ auto-grade, manual grading inbox
  ac          : US-M07-03, US-M07-04
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-046
  title       : Test phase 2 selection + leaderboard
  description : GET leaderboard (ranked by total score), POST phase2 selection with cutoff. Close round
  ac          : US-M07-05
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 4

TASK: ATS-047
  title       : Resume review 3-panel UI
  description : Left: candidate info, Center: resume/OCR viewer, Right: scoring panel + weight popup + position marks. TigerSoft branding
  ac          : US-M06-01, US-M06-02, US-M06-03, US-M06-04
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-048
  title       : Test management UI — library + round creation + grading
  description : Test library table, create/edit test, assign applicants to round, grading inbox, leaderboard view
  ac          : US-M07-01..06
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 4

TASK: ATS-049
  title       : Go tests + Playwright E2E: resume review, test management
  description : Backend tests for resume review + test APIs. E2E: resume review flow, test creation, grading, leaderboard
  ac          : All M06/M07 stories
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 4
```

### Sprint 4 Summary

```
Sprint 4 — Resume Review + Test Management (M06, M07)
Duration      : 2 weeks
Total Tasks   : 8
Total Points  : 38
Dev Tasks     : 7 (33 pts)
QA Tasks      : 1 (5 pts)
Capacity Used : 95% — slightly over, acceptable given parallel BE/FE work
```

### Dependencies
- ATS-042, ATS-043 (resume review BE) → blocks ATS-047 (resume review FE)
- ATS-044, ATS-045 (test BE) → blocks ATS-048 (test FE)
- ATS-043 depends on ai-service /embed + /match endpoints being available (mock OK)
- **Critical path**: ATS-042 → ATS-047 → ATS-049

---

## SPRINT 5 — Interview + Transfer (M08, M09)

### Task List

```
TASK: ATS-050
  title       : Interview round CRUD + scoring criteria
  description : POST interview_rounds (linked to pipeline_action), POST scoring_criteria per round
  ac          : US-M08-01, US-M08-02
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 5

TASK: ATS-051
  title       : Assign interviewers + invite links
  description : POST assign interviewers (supervisor + temp). Create invite_links with 256-bit token, manage status (active/used/revoked/expired)
  ac          : US-M08-03, US-M08-04
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 5

TASK: ATS-052
  title       : Interview grading + leaderboard + email invite
  description : POST evaluation (per-interviewer per-criterion), GET leaderboard, phase2 selection. Send TPL05 interview invite email
  ac          : US-M08-05, US-M08-06, US-M08-07
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 5

TASK: ATS-053
  title       : Transfer API + history + validations
  description : POST /api/transfer — transfer applicant to new job, carry history. Warn if target closed. Block test from wrong position
  ac          : US-M09-01, US-M09-02, US-M09-03, US-M09-04
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 5

TASK: ATS-054
  title       : Interview management UI
  description : Create round, define criteria, assign interviewers, manage invite links, grading interface, leaderboard. TigerSoft branding
  ac          : US-M08-01..07
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 5

TASK: ATS-055
  title       : Transfer UI + history view
  description : Transfer action button (from pipeline/interview), target job selector with closed-warning, transfer history timeline
  ac          : US-M09-01..04
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 5

TASK: ATS-056
  title       : Go tests + Playwright E2E: interview, transfer
  description : Backend tests for interview + transfer APIs. E2E: create interview round, grade, transfer flow
  ac          : All M08/M09 stories
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 5
```

### Sprint 5 Summary

```
Sprint 5 — Interview + Transfer (M08, M09)
Duration      : 2 weeks
Total Tasks   : 7
Total Points  : 31
Dev Tasks     : 6 (26 pts)
QA Tasks      : 1 (5 pts)
Capacity Used : 78% — within capacity
```

### Dependencies
- ATS-050..052 (interview BE) → blocks ATS-054 (interview FE)
- ATS-053 (transfer BE) → blocks ATS-055 (transfer FE)
- ATS-051 depends on auth middleware supporting temp_interviewer role
- **Critical path**: ATS-050 → ATS-052 → ATS-054 → ATS-056

---

## SPRINT 6 — Talent Pool + Candidates + Requisitions (M10, M11, M12)

### Task List

```
TASK: ATS-057
  title       : Talent pool auto-compute + blacklist removal
  description : On new job publish, compute similarity for all starred candidates via ai-service. Remove from pool on blacklist
  ac          : US-M10-01, US-M10-02
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 6

TASK: ATS-058
  title       : Candidate list + detail + OCR trigger
  description : GET /api/candidates with search/filter, GET detail with full history. Auto-trigger OCR on resume upload
  ac          : US-M11-01, US-M11-02, US-M11-03
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 6

TASK: ATS-059
  title       : Blacklist/unblacklist with reason + audit log
  description : POST blacklist (mandatory reason), POST unblacklist (mandatory reason). Write to blacklist_logs. History preserved
  ac          : US-M11-04, US-M11-05
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 6

TASK: ATS-060
  title       : Requisition CRUD + convert to job
  description : POST/PUT requisitions, workflow (draft→submitted→approved→converted). 1-click convert to job posting pre-fill
  ac          : US-M12-01, US-M12-02, US-M12-03
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 6

TASK: ATS-061
  title       : Candidate management UI
  description : Candidate list table with search/filter, detail page with timeline/docs, blacklist/unblacklist dialogs. TigerSoft branding
  ac          : US-M11-01..05
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 6

TASK: ATS-062
  title       : Requisition workflow UI
  description : Requisition form (supervisor), review panel (HR), convert-to-job button with pre-fill preview. TigerSoft branding
  ac          : US-M12-01..03
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 6

TASK: ATS-063
  title       : Go tests + Playwright E2E: candidates, requisitions, talent pool
  description : Backend tests for candidate/requisition/talent pool APIs. E2E: candidate search, blacklist, requisition flow
  ac          : All M10/M11/M12 stories
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 6
```

### Sprint 6 Summary

```
Sprint 6 — Talent Pool + Candidates + Requisitions (M10, M11, M12)
Duration      : 2 weeks
Total Tasks   : 7
Total Points  : 29
Dev Tasks     : 6 (24 pts)
QA Tasks      : 1 (5 pts)
Capacity Used : 73% — within capacity, buffer for candidate detail complexity
```

---

## SPRINT 7 — Notifications + Setup + i18n/Theme (M13, M14, M15)

### Task List

```
TASK: ATS-064
  title       : Notification API + 10 event triggers
  description : GET /api/notifications, PUT mark-read. Trigger creation on all 10 defined events (new application, status change, etc.)
  ac          : US-M13-01, US-M13-03
  tag         : backend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 7

TASK: ATS-065
  title       : Notification preferences API
  description : GET/PUT /api/notifications/preferences — per-user per-event toggle for in-app vs email channel
  ac          : US-M13-02
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 7

TASK: ATS-066
  title       : Email template management API
  description : GET/PUT /api/setup/email-templates — list system+custom templates, edit custom templates (body_html, subject, variables)
  ac          : US-M14-01
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 7

TASK: ATS-067
  title       : JD template CRUD + master data management
  description : CRUD for jd_templates. CRUD for master data (departments, employment_types, experience_levels, hard_skills)
  ac          : US-M14-02, US-M14-03
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 7

TASK: ATS-068
  title       : JD + skill embedding trigger on job create/edit
  description : After job create/edit, async call ai-service /embed for JD text + aggregated hard skills, store in job_embeddings
  ac          : US-M14-04
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 7

TASK: ATS-069
  title       : Notification bell UI + preferences page
  description : Bell icon with unread badge (header), dropdown list, mark-read. Preferences page with toggles. TigerSoft branding
  ac          : US-M13-01, US-M13-02
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 7

TASK: ATS-070
  title       : Setup pages — email templates, JD templates, master data
  description : Template editor with variable preview, JD template list/edit, master data CRUD tables. TigerSoft branding
  ac          : US-M14-01..03
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 7

TASK: ATS-071
  title       : i18n — Thai/English toggle + translation keys
  description : Configure next-intl, create th.json + en.json, replace all hardcoded text with translation keys, locale toggle in header
  ac          : US-M15-01, US-M15-02
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 7

TASK: ATS-072
  title       : Dark/Light theme + WCAG AA compliance
  description : CSS variables for light/dark, theme toggle in header, verify WCAG 2.1 AA contrast for both themes
  ac          : US-M15-03, US-M15-04
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 7

TASK: ATS-073
  title       : Go tests + Playwright E2E: notifications, setup, i18n, theme
  description : Backend tests for notification/setup APIs. E2E: notification bell, template editor, locale toggle, theme toggle, contrast check
  ac          : All M13/M14/M15 stories
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 7
```

### Sprint 7 Summary

```
Sprint 7 — Notifications + Setup + i18n/Theme (M13, M14, M15)
Duration      : 2 weeks
Total Tasks   : 10
Total Points  : 36
Dev Tasks     : 9 (31 pts)
QA Tasks      : 1 (5 pts)
Capacity Used : 90% — slightly over, acceptable for final sprint
```

---

## GLOBAL RISK REGISTER

| # | Risk | Prob | Impact | Score | Mitigation | Owner |
|---|------|------|--------|-------|------------|-------|
| R1 | API contract not finalized before parallel FE/BE work | Med | High | H | API contracts defined in P2 — FE can mock from contract. Review after BE done | SA/Dev |
| R2 | AI service latency affecting OCR/embedding integration | Med | Med | M | Async processing, mock ai-service in dev/test. Timeouts + retry | Dev |
| R3 | Scope creep from mid-sprint stakeholder changes | Med | High | H | All stories locked in P1, changes go through PO gate. Sprint backlog frozen | PO |
| R4 | Key developer unavailability | Low | High | M | All tasks documented with AC, any dev can pick up. Pair programming on critical path | PM |
| R5 | Environment/infra instability blocking testing | Low | High | M | Docker Compose with health checks, seed data in migrations. QA has local env | Dev |
| R6 | Form builder complexity (ATS-038, 8 pts) | High | Med | H | Flag for split into basic+advanced. Start early in Sprint 3. Prototype first | Dev |
| R7 | Anti-cheat false positives in online tests | Med | Med | M | Configurable thresholds, manual review option, test in staging | QA |
| R8 | i18n retrofit across all existing pages (Sprint 7) | Med | Med | M | Use translation keys from start where possible. Sprint 7 is final pass | Dev |

---

## FULL PLAN SUMMARY

| Sprint | Modules | Tasks | Points | Capacity |
|--------|---------|-------|--------|----------|
| Sprint 1 | M01, M02, M03-partial | 21 | 46 | 115% (infra overhead) |
| Sprint 2 | M03-remainder, M04, M05 (BE) | 20 | 38* | 95% |
| Sprint 3 | M04, M05 (FE complex) | 6 | 27 | 68% |
| Sprint 4 | M06, M07 | 8 | 38 | 95% |
| Sprint 5 | M08, M09 | 7 | 31 | 78% |
| Sprint 6 | M10, M11, M12 | 7 | 29 | 73% |
| Sprint 7 | M13, M14, M15 | 10 | 36 | 90% |
| **TOTAL** | M01–M15 | **79** | **245** | 7 sprints (14 weeks) |

*Sprint 2 reduced after deferring complex FE to Sprint 3
