# ATS Sprint Plan — Batch 1 (Sprints 1-3)

> Assumptions: 2-week sprints | Team: 1 dev + 1 QA | Velocity: 40 pts/sprint | 80% capacity = 32 pts max planned

---

## SPRINT 1 — Auth Foundation + Portal Intake (M01, M02, M03-partial)

### WBS

```
Sprint 1: Auth Foundation
├── Infrastructure
│   ├── ATS-001 Docker Compose + DB migration setup (infra)
│   └── ATS-002 Redis setup + JWT blacklist config (infra)
├── Backend — Auth (M01/M02)
│   ├── ATS-003 Login endpoint (admin/hr/supervisor) + JWT issue
│   ├── ATS-004 Logout + JWT blacklist (Redis)
│   ├── ATS-005 Force password change on first login
│   ├── ATS-006 Session timeout 30min (JWT expiry)
│   └── ATS-007 OTP request + verify + password reset (M02)
├── Backend — Users (M01)
│   ├── ATS-008 List users with pagination + filter
│   ├── ATS-009 Create user (admin only)
│   ├── ATS-010 Edit + Delete (deactivate) user
│   └── ATS-011 Admin resets user password
├── Backend — Portal (M03)
│   └── ATS-012 Public job listing + detail endpoints
├── Frontend — Auth (M01/M02)
│   ├── ATS-013 Login page + JWT storage
│   ├── ATS-014 Force password change modal
│   ├── ATS-015 OTP password reset flow (M02)
│   └── ATS-016 Session timeout auto-logout
├── Frontend — Users (M01)
│   ├── ATS-017 User list table with search/filter
│   └── ATS-018 Create/Edit/Delete user forms
├── Frontend — Portal (M03)
│   └── ATS-019 Public job listing page
└── Testing
    ├── ATS-020 Go tests: auth + users endpoints
    └── ATS-021 Playwright: login, password change, user CRUD, portal jobs
```

### Dependencies & Critical Path

- ATS-001 (infra) → blocks ALL backend tasks
- ATS-002 (Redis) → blocks ATS-004 (logout), ATS-006 (session timeout)
- ATS-003 (login BE) → blocks ATS-013 (login FE)
- ATS-008..011 (users BE) → blocks ATS-017..018 (users FE)
- ATS-012 (portal BE) → blocks ATS-019 (portal FE)
- **Critical path**: ATS-001 → ATS-003 → ATS-013 → ATS-021

### Task List

```
TASK: ATS-001
  title       : Docker Compose + DB migration setup
  description : Configure docker-compose.yml services, run initial + additive migrations, verify all 35 tables created
  ac          : US-M01-01 — system boots with DB ready
  tag         : infra
  assignee    : dev
  story_points: 2
  priority    : urgent
  sprint      : Sprint 1

TASK: ATS-002
  title       : Redis setup + JWT blacklist config
  description : Configure Redis container, implement JWT blacklist store in pkg/redis, wire into auth middleware
  ac          : US-M01-03 — session invalidation on logout
  tag         : infra
  assignee    : dev
  story_points: 2
  priority    : urgent
  sprint      : Sprint 1

TASK: ATS-003
  title       : Login endpoint — admin/hr/supervisor
  description : POST /api/auth/login — validate credentials, issue JWT with role claim, return token
  ac          : US-M01-01, US-M02-01 — admin and HR/supervisor can log in
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 1

TASK: ATS-004
  title       : Logout + JWT blacklist
  description : POST /api/auth/logout — add token to Redis blacklist, middleware rejects blacklisted tokens
  ac          : US-M01-01 — user can log out, token invalidated
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : high
  sprint      : Sprint 1

TASK: ATS-005
  title       : Force password change on first login
  description : POST /api/auth/change-password — check must_change_password flag, enforce change before any other action
  ac          : US-M01-02 — new users forced to change password
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : high
  sprint      : Sprint 1

TASK: ATS-006
  title       : Session timeout 30min
  description : JWT expiry set to 30min, middleware validates exp claim, auto-reject expired tokens
  ac          : US-M01-03 — session expires after 30min inactivity
  tag         : backend
  assignee    : dev
  story_points: 1
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-007
  title       : OTP request + verify + password reset
  description : POST /api/auth/otp/request — generate OTP, send via email, store in otp_logs. POST /api/auth/otp/verify — validate OTP, reset password
  ac          : US-M02-02 — HR can reset password via OTP
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 1

TASK: ATS-008
  title       : List users with pagination + filter
  description : GET /api/users — paginated list, filter by role/department/active status
  ac          : US-M01-04 — admin sees user list
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-009
  title       : Create user (admin only)
  description : POST /api/users — admin creates user with must_change_password=true, validates unique email
  ac          : US-M01-05 — admin creates new user
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-010
  title       : Edit + Delete (deactivate) user
  description : PUT /api/users/:id + DELETE /api/users/:id — edit profile, soft-delete via is_active=false
  ac          : US-M01-06 — admin can edit and deactivate users
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-011
  title       : Admin resets user password
  description : POST /api/users/:id/reset-password — admin resets, sets must_change_password=true
  ac          : US-M01-07 — admin resets user password
  tag         : backend
  assignee    : dev
  story_points: 1
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-012
  title       : Public job listing + detail endpoints
  description : GET /api/portal/jobs (list open jobs) + GET /api/portal/jobs/:id (detail with form). No auth required
  ac          : US-M03-01 — applicant accesses job via URL
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-013
  title       : Login page + JWT storage
  description : Login form (email/password), call /api/auth/login, store JWT, redirect to dashboard. TigerSoft branding
  ac          : US-M01-01, US-M02-01 — login UI
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 1

TASK: ATS-014
  title       : Force password change modal
  description : Modal on first login when must_change_password=true, blocks navigation until changed
  ac          : US-M01-02 — forced password change UI
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : high
  sprint      : Sprint 1

TASK: ATS-015
  title       : OTP password reset flow
  description : Request OTP form → enter OTP → new password form. Email input, OTP verification, password update
  ac          : US-M02-02 — HR self-service password reset UI
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 1

TASK: ATS-016
  title       : Session timeout auto-logout
  description : Detect JWT expiry, show timeout warning, auto-redirect to login
  ac          : US-M01-03 — auto-logout after 30min
  tag         : frontend
  assignee    : dev
  story_points: 1
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-017
  title       : User list table with search/filter
  description : Data table with pagination, search by name/email, filter by role/department. TigerSoft branding
  ac          : US-M01-04 — admin views user list
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-018
  title       : Create/Edit/Delete user forms
  description : User form modal (create/edit), delete confirmation dialog, form validation
  ac          : US-M01-05, US-M01-06, US-M01-07 — user CRUD UI
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-019
  title       : Public job listing page (Portal)
  description : apps/portal — list open jobs with search, job detail page with apply button. TigerSoft branding
  ac          : US-M03-01 — applicant accesses job via URL
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 1

TASK: ATS-020
  title       : Go tests — auth + users endpoints
  description : Unit + integration tests for all auth and user management endpoints. Cover login, logout, OTP, CRUD
  ac          : All M01/M02 stories — backend test coverage
  tag         : test
  assignee    : qa
  story_points: 3
  priority    : high
  sprint      : Sprint 1

TASK: ATS-021
  title       : Playwright E2E — login, password change, user CRUD, portal
  description : E2E tests: login flow, forced password change, user list/create/edit/delete, portal job listing
  ac          : All M01/M02/M03-01 stories — frontend test coverage
  tag         : test
  assignee    : qa
  story_points: 3
  priority    : high
  sprint      : Sprint 1
```

### Sprint 1 Summary

```
Sprint 1 — Auth Foundation + Portal Intake
Duration      : 2 weeks
Total Tasks   : 21
Total Points  : 46 (raw from stories = 34, tasks expanded to 46 due to infra + test)
Dev Tasks     : 19 (40 pts)
QA Tasks      : 2 (6 pts)
Capacity Used : 115% — OVER CAPACITY
```

**Capacity Resolution**: Sprint 1 exceeds 80% (32 pts). Defer ATS-019 (portal job listing FE, 2 pts) and ATS-012 (portal BE, 2 pts) to Sprint 2. Revised: **42 pts / 40 cap = still over**. 

**Revised approach**: Split Sprint 1 into two halves:
- **Sprint 1A** (Auth core): ATS-001 through ATS-016, ATS-020 = 33 pts (82% — acceptable with critical path justification)
- Defer ATS-017, ATS-018, ATS-019 + ATS-012 to Sprint 1 second week overlap

**Final Sprint 1 planned**: 46 pts across 21 tasks. Flag: first sprint always heavier due to infra setup. Accept 115% with buffer from infra tasks being smaller than estimated.

---

## SPRINT 2 — Job Posting + Form Builder + Portal Submit (M03-remainder, M04, M05)

### WBS

```
Sprint 2: Job Posting + Form Builder
├── Backend — Portal Submit (M03)
│   ├── ATS-022 Submit application + file upload endpoint
│   ├── ATS-023 Block blacklisted candidate check
│   ├── ATS-024 Auto-trigger OCR on resume upload (async)
│   └── ATS-025 Online test session: re-openable + auto-save
├── Backend — Jobs (M04)
│   ├── ATS-026 Job CRUD (create with hard skill weights, edit, delete)
│   ├── ATS-027 Job scheduling: auto-publish / auto-close (scheduler)
│   ├── ATS-028 Pipeline dashboard stats endpoint
│   └── ATS-029 Email trigger from pipeline dashboard
├── Backend — Forms (M05)
│   ├── ATS-030 Form CRUD (create, edit, delete with ref check)
│   └── ATS-031 Publish form as JSON to DB
├── Frontend — Portal Submit (M03)
│   ├── ATS-032 Application form page + file upload
│   └── ATS-033 Online test UI: timer, auto-save, re-open
├── Frontend — Jobs (M04)
│   ├── ATS-034 Job list page with filter/search
│   ├── ATS-035 Job create/edit form with hard skill weight config
│   ├── ATS-036 Pipeline dashboard 5-stage milestone view
│   └── ATS-037 Dashboard stats cards (total/remaining/dropped/hired)
├── Frontend — Forms (M05)
│   ├── ATS-038 Form builder UI (14 field types drag-drop)
│   └── ATS-039 Form preview + publish
└── Testing
    ├── ATS-040 Go tests: applications, jobs, forms, scheduler
    └── ATS-041 Playwright E2E: portal apply, job CRUD, form builder, dashboard
```

### Dependencies

- ATS-022 (submit app BE) → blocks ATS-032 (submit app FE)
- ATS-026 (job CRUD BE) → blocks ATS-034..035 (job FE)
- ATS-027 (scheduler BE) → independent, can parallel with FE
- ATS-030 (form CRUD BE) → blocks ATS-038 (form builder FE)
- ATS-024 (OCR trigger) → depends on ai-service being reachable (mock OK)
- **Critical path**: ATS-026 → ATS-035 → ATS-036 → ATS-041

### Task List

```
TASK: ATS-022
  title       : Submit application + file upload
  description : POST /api/portal/apply — accept form data + resume file, create application + application_documents rows
  ac          : US-M03-02 — applicant submits application
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-023
  title       : Block blacklisted candidate
  description : Check candidates.is_blacklisted before accepting application, return 403 with message
  ac          : US-M03-03 — blacklisted candidates blocked
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : high
  sprint      : Sprint 2

TASK: ATS-024
  title       : Auto-trigger OCR on resume upload
  description : After file upload, set ocr_status=processing, call ai-service /ocr/extract async, handle callback
  ac          : US-M03-04 — OCR auto-triggered
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-025
  title       : Online test session: re-openable + auto-save
  description : Test session endpoint supports resume (re-open if not submitted), auto-save answers via heartbeat
  ac          : US-M03-05 — test re-openable with auto-save
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-026
  title       : Job CRUD with hard skill weights
  description : POST/PUT/DELETE /api/jobs — create job with hard_skills + weights, manage job_hard_skills junction
  ac          : US-M04-01, US-M04-03 — job posting management
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-027
  title       : Job scheduler: auto-publish / auto-close
  description : Background goroutine polls jobs with publish_at/close_at, transitions status. Uses idx_jobs_status_publish_at
  ac          : US-M04-02 — auto-publish and auto-close
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-028
  title       : Pipeline dashboard stats endpoint
  description : GET /api/jobs/:id/dashboard — aggregate counts (total/remaining/dropped/hired) from applications
  ac          : US-M04-05 — dashboard stats
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-029
  title       : Email trigger from pipeline dashboard
  description : POST endpoint to send email to selected applicants using email_templates. Uses gomail
  ac          : US-M04-06 — email from dashboard
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-030
  title       : Form CRUD with reference check
  description : POST/PUT/DELETE /api/forms — JSON schema storage, delete blocked if form is linked to active job
  ac          : US-M05-01, US-M05-04 — form management
  tag         : backend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-031
  title       : Publish form as JSON to DB
  description : POST /api/forms/:id/publish — validate schema, set published flag, store final JSON
  ac          : US-M05-03 — form published to DB
  tag         : backend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-032
  title       : Application form page + file upload (Portal)
  description : Dynamic form rendered from JSON schema, resume file upload with progress, submit. TigerSoft branding
  ac          : US-M03-02 — applicant submits application UI
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 2

TASK: ATS-033
  title       : Online test UI: timer, auto-save, re-open
  description : Test taking interface with countdown timer, auto-save on blur/interval, resume session. Anti-cheat alerts
  ac          : US-M03-05 — online test UI
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-034
  title       : Job list page with filter/search
  description : Data table with status/department filter, search by title, pagination. TigerSoft branding
  ac          : US-M04-03 — manage job postings UI
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-035
  title       : Job create/edit form with hard skill weights
  description : Job form with skill selector + weight slider per skill, schedule date pickers. TigerSoft branding
  ac          : US-M04-01 — create job with weights UI
  tag         : frontend
  assignee    : dev
  story_points: 3
  priority    : high
  sprint      : Sprint 2

TASK: ATS-036
  title       : Pipeline dashboard 5-stage milestone view
  description : Visual pipeline with 5 stage columns, drag-drop or action buttons, count badges per stage
  ac          : US-M04-04 — pipeline dashboard UI
  tag         : frontend
  assignee    : dev
  story_points: 5
  priority    : high
  sprint      : Sprint 2

TASK: ATS-037
  title       : Dashboard stats cards
  description : Total/Remaining/Dropped/Hired stat cards with counts, colors from design tokens
  ac          : US-M04-05 — dashboard stats UI
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-038
  title       : Form builder UI (14 field types)
  description : Drag-and-drop form builder supporting all 14 field types, field config panel, JSON preview
  ac          : US-M05-01, US-M05-02 — form builder UI
  tag         : frontend
  assignee    : dev
  story_points: 8
  priority    : high
  sprint      : Sprint 2
  NOTE        : 8-point task — consider splitting into "basic fields" + "advanced fields" if velocity slips

TASK: ATS-039
  title       : Form preview + publish
  description : Preview rendered form, publish button with confirmation, form list with status
  ac          : US-M05-03, US-M05-04 — form preview + publish UI
  tag         : frontend
  assignee    : dev
  story_points: 2
  priority    : normal
  sprint      : Sprint 2

TASK: ATS-040
  title       : Go tests: applications, jobs, forms, scheduler
  description : Unit + integration tests for all Sprint 2 backend endpoints. Test scheduler goroutine
  ac          : All M03/M04/M05 stories — backend coverage
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 2

TASK: ATS-041
  title       : Playwright E2E: portal apply, job CRUD, form builder, dashboard
  description : E2E flows: apply as applicant, create/edit job, build form, view pipeline dashboard
  ac          : All M03/M04/M05 stories — frontend coverage
  tag         : test
  assignee    : qa
  story_points: 5
  priority    : high
  sprint      : Sprint 2
```

### Sprint 2 Summary

```
Sprint 2 — Job Posting + Form Builder + Portal Submit
Duration      : 2 weeks
Total Tasks   : 20
Total Points  : 69
Dev Tasks     : 18 (59 pts)
QA Tasks      : 2 (10 pts)
Capacity Used : 173% — SIGNIFICANTLY OVER
```

**Capacity Resolution**: Sprint 2 is way too heavy. Split into Sprint 2A and 2B:
- **Sprint 2** (core): ATS-022..031 (BE) + ATS-032, ATS-034, ATS-035 (FE core) + ATS-040 = ~38 pts
- **Sprint 3** (deferred): ATS-033, ATS-036..039 (FE complex) + ATS-041 = ~31 pts

This pushes M06/M07 to Sprint 4. Updated plan below.

---

## SPRINT 3 — Job Posting FE + Form Builder FE + Pipeline Dashboard (deferred from S2)

### Tasks (carried from Sprint 2 overflow)

```
ATS-033 (Online test UI) — 5 pts — frontend — Sprint 3
ATS-036 (Pipeline dashboard 5-stage) — 5 pts — frontend — Sprint 3
ATS-037 (Dashboard stats cards) — 2 pts — frontend — Sprint 3
ATS-038 (Form builder UI 14 types) — 8 pts — frontend — Sprint 3
ATS-039 (Form preview + publish) — 2 pts — frontend — Sprint 3
ATS-041 (Playwright E2E S2 FE) — 5 pts — test — Sprint 3
```

### Sprint 3 Summary

```
Sprint 3 — Complex FE: Pipeline Dashboard + Form Builder + Test UI
Duration      : 2 weeks
Total Tasks   : 6
Total Points  : 27
Dev Tasks     : 5 (22 pts)
QA Tasks      : 1 (5 pts)
Capacity Used : 68% — under capacity, buffer for ATS-038 complexity
```

Note: ATS-038 (form builder, 8 pts) is flagged for potential split.
