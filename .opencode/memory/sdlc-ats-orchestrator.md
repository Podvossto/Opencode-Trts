# SDLC Orchestrator — Persistent Memory

> Keep under 200 lines. Update as you learn. Remove stale entries.

## Notion Configuration (Confirmed IDs)

| Property | Value |
|---|---|
| Notion API Token | `ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA` |
| Parent Page (ATS Recruitment System) | `3220531c-537a-80aa-8333-c42da4ce362f` |
| Team & Agents Subpage | `3250531c-537a-8104-b718-d3e52f8d296e` |
| Sprint Planning Subpage | `3250531c-537a-81de-b835-f2153c00351c` |
| Architecture & API Subpage | `3250531c-537a-8184-b587-c6d54543c22e` |
| Activity Log Subpage | `3250531c-537a-8118-a1d0-c14de74dec51` |
| Team Directory DB | `3250531c-537a-813a-b729-f21f311568f5` |
| Sprint Roadmap DB | `3250531c-537a-811d-a96d-d757c9080f9a` |
| Sprint 1 Kanban DB | `3250531c-537a-813b-abfd-f0fdb4d05df1` |
| API Registry DB | `3250531c-537a-8125-85f0-f9db90e9b9e9` |
| Activity Log DB | `3250531c-537a-812d-be1b-e3a658c3f682` |

Notion Mode: **LIVE** (workspace fully rebuilt 2026-03-16 — grouped subpage structure)

## Phase 1 Status

| Batch | Modules | Stories | Points | Status |
|---|---|---|---|---|
| Batch 1 | M01–M05 | 24 | 86 | ✅ COMPLETE (2026-03-16) |
| Batch 2 | M06–M15 | 43 | 143 | ✅ COMPLETE (2026-03-16) |
| **TOTAL** | M01–M15 | **67** | **229** | ✅ P1 DONE |

## Phase 2 Status

| Output | Status |
|---|---|
| Architecture Overview (5 services, Docker Compose) | ✅ COMPLETE |
| Database Schema (35 tables, 2 migrations) | ✅ COMPLETE |
| API Contract (65+ endpoints) | ✅ COMPLETE |
| Folder Structure & Scaffolds (15 Go domains, all FE features) | ✅ COMPLETE |
| ADRs (9 decisions) | ✅ COMPLETE |
| Non-Functional Assessment | ✅ COMPLETE |
| Design Tokens (TigerSoft CI) | ✅ COMPLETE |
| **P2 Gate** | **✅ APPROVED (2026-03-16)** |

Phase Report: `.opencode/handoffs/p2-sa-phase-report.txt`

## Phase 3 Status

| Output | Status |
|---|---|
| Sprint Plan (7 sprints, 79 tasks, 245 pts, 14 weeks) | ✅ COMPLETE |
| Task breakdown ATS-001 through ATS-073 | ✅ COMPLETE |
| Risk Register (8 risks) | ✅ COMPLETE |
| Sprint Plan Batch 1 (Sprints 1-3) | ✅ COMPLETE |
| Sprint Plan Batch 2 (Sprints 4-7) | ✅ COMPLETE |
| Phase Report | ✅ COMPLETE |
| **P3 Gate** | **✅ APPROVED (2026-03-16)** |

Phase Report: `.opencode/handoffs/p3-pm-phase-report.txt`

## Sprint & Feature Mappings (Finalized)

| Sprint | Modules | Tasks | Points | Focus |
|---|---|---|---|---|
| Sprint 1 | M01, M02, M03-partial | 21 (ATS-001..021) | 46 | Auth + Infra foundation |
| Sprint 2 | M03-rest, M04, M05 (BE) | 20 (ATS-022..041*) | 38 | Job/Form/Portal backend |
| Sprint 3 | M04, M05 (FE complex) | 6 (ATS-033..041) | 27 | Pipeline dashboard + Form builder UI |
| Sprint 4 | M06, M07 | 8 (ATS-042..049) | 38 | Resume review + Test mgmt |
| Sprint 5 | M08, M09 | 7 (ATS-050..056) | 31 | Interview + Transfer |
| Sprint 6 | M10, M11, M12 | 7 (ATS-057..063) | 29 | Talent pool + Candidates + Requisitions |
| Sprint 7 | M13, M14, M15 | 10 (ATS-064..073) | 36 | Notifications + Setup + i18n/Theme |

## User Story Index (Batch 1)

| Story ID | Title | Module | Points | Priority |
|---|---|---|---|---|
| US-M01-01 | Admin Login | M01 | 3 | Must |
| US-M01-02 | Forced Password Change on First Login | M01 | 3 | Must |
| US-M01-03 | Session Timeout 30min | M01 | 2 | Must |
| US-M01-04 | List Users | M01 | 2 | Must |
| US-M01-05 | Create New User | M01 | 3 | Must |
| US-M01-06 | Edit and Delete User | M01 | 3 | Must |
| US-M01-07 | Admin Resets User Password | M01 | 2 | Must |
| US-M02-01 | HR/Supervisor Login | M02 | 2 | Must |
| US-M02-02 | HR Self-Service OTP Password Reset | M02 | 5 | Must |
| US-M03-01 | Applicant Accesses Job via URL | M03 | 3 | Must |
| US-M03-02 | Submit Application Form | M03 | 5 | Must |
| US-M03-03 | Block Blacklisted Candidate | M03 | 3 | Must |
| US-M03-04 | Auto-trigger OCR on Resume | M03 | 3 | Must |
| US-M03-05 | Online Test Re-openable + Auto-save | M03 | 5 | Must |
| US-M04-01 | Create Job Posting with Hard Skill Weights | M04 | 5 | Must |
| US-M04-02 | Job Scheduling Auto-Publish / Auto-Close | M04 | 5 | Must |
| US-M04-03 | Edit, Filter, Manage Job Postings | M04 | 3 | Must |
| US-M04-04 | Pipeline Dashboard 5-Stage Milestone View | M04 | 5 | Must |
| US-M04-05 | Dashboard Stats (Total/Remaining/Dropped/Hired) | M04 | 3 | Must |
| US-M04-06 | Email Trigger from Pipeline Dashboard | M04 | 5 | Must |
| US-M05-01 | Build Form with Standard Fields | M05 | 5 | Must |
| US-M05-02 | Add Custom Question Types | M05 | 5 | Must |
| US-M05-03 | Publish Form as JSON to Database | M05 | 3 | Must |
| US-M05-04 | Edit and Delete Forms with Reference Check | M05 | 3 | Must |

## User Story Index (Batch 2)

| Story ID | Title | Module | Points | Priority |
|---|---|---|---|---|
| US-M06-01 | 3-Panel Resume Review Layout | M06 | 3 | Must |
| US-M06-02 | Overall Score Calculation with Weight Popup | M06 | 3 | Must |
| US-M06-03 | HR Marks Up to 3 Suitable Positions | M06 | 2 | Must |
| US-M06-04 | Resume Review Action Buttons | M06 | 5 | Must |
| US-M06-05 | Auto-Compute Similarity for Starred Candidates on New Job Publish | M06 | 3 | Must |
| US-M07-01 | HR Initiates Action Test with Round Badge | M07 | 3 | Must |
| US-M07-02 | HR Selects Applicants and Assigns Test | M07 | 3 | Must |
| US-M07-03 | Online Test Delivery with Anti-Cheat | M07 | 5 | Must |
| US-M07-04 | MCQ Auto-Grading and Manual Grading Inbox | M07 | 5 | Must |
| US-M07-05 | HR Opens Phase 2 Selection and Leaderboard | M07 | 5 | Must |
| US-M07-06 | HR and Supervisor CRUD Test Library | M07 | 3 | Must |
| US-M08-01 | HR Initiates Action Interview with Round Badge | M08 | 3 | Must |
| US-M08-02 | HR Defines Scoring Criteria per Round | M08 | 3 | Must |
| US-M08-03 | HR Assigns Interviewers — Supervisor and Temp Interviewers | M08 | 5 | Must |
| US-M08-04 | HR Views and Manages Invite Links List | M08 | 2 | Must |
| US-M08-05 | Phase 1 Interview Grading — Real-Time Per-Interviewer Scores | M08 | 5 | Must |
| US-M08-06 | Phase 2 Interview Selection — Leaderboard and Final Decisions | M08 | 5 | Must |
| US-M08-07 | HR Sends Interview Invitation Email (TPL05) | M08 | 2 | Must |
| US-M09-01 | HR Initiates Transfer During or After Interview | M09 | 5 | Must |
| US-M09-02 | Transfer History Carries Over and Unlimited Transfers Supported | M09 | 3 | Must |
| US-M09-03 | System Warns HR if Target Position is Closed | M09 | 2 | Must |
| US-M09-04 | Block Sending Test from Wrong Position After Transfer | M09 | 2 | Must |
| US-M10-01 | Auto-Compute Similarity for Starred Candidates Against All New Jobs | M10 | 3 | Must |
| US-M10-02 | Starred Candidates Removed from Pool on Blacklist | M10 | 2 | Must |
| US-M11-01 | Candidate List with Search and Filter | M11 | 3 | Must |
| US-M11-02 | Candidate Detail — Full Application History and Documents | M11 | 5 | Must |
| US-M11-03 | OCR Auto-Trigger on Resume Upload | M11 | 3 | Must |
| US-M11-04 | Blacklist Candidate with Mandatory Reason | M11 | 3 | Must |
| US-M11-05 | Unblacklist Candidate with Mandatory Reason and History Preserved | M11 | 2 | Must |
| US-M12-01 | Supervisor Creates and Submits Requisition | M12 | 3 | Must |
| US-M12-02 | HR Reviews and Edits Requisition | M12 | 3 | Must |
| US-M12-03 | HR Converts Requisition to Job Posting (1-Click Pre-Fill) | M12 | 3 | Must |
| US-M13-01 | In-App Notification Bell with Unread Badge | M13 | 3 | Must |
| US-M13-02 | Per-User Notification Channel Toggle (In-App vs Email) | M13 | 3 | Must |
| US-M13-03 | System Triggers All 10 Defined Notification Events | M13 | 5 | Must |
| US-M14-01 | HR Manages System and Custom Email Templates | M14 | 5 | Must |
| US-M14-02 | HR Creates and Manages Job Description Templates | M14 | 2 | Must |
| US-M14-03 | HR Manages Master Data | M14 | 5 | Must |
| US-M14-04 | JD and Hard Skill Embedding Sent to AI Service After Job Create/Edit | M14 | 3 | Must |
| US-M15-01 | Thai / English Language Toggle | M15 | 3 | Must |
| US-M15-02 | All UI Text via Translation Keys (next-intl) | M15 | 3 | Must |
| US-M15-03 | Dark / Light Theme Toggle with CSS Variables | M15 | 3 | Must |
| US-M15-04 | WCAG 2.1 AA Contrast Compliance for Both Themes | M15 | 2 | Must |

## Architectural Decisions

| Decision | Detail | Date |
|---|---|---|
| Auth module | M01+M02 share one auth module — implement once, roles differentiated by flag | 2026-03-16 |
| OCR trigger | M03 OCR is async via FastAPI AI service queue — must not block form submission | 2026-03-16 |
| Auto-publish/close | M04 requires background scheduler (Go worker / cron) | 2026-03-16 |
| Form JSON schema | Must support 14 field types: Text/Textarea/Dropdown/Checkbox/Radio/FileUpload/DatePicker/MCQ/ShortAnswer/Essay/TrueFalse/RatingScale/Image/Video/Link | 2026-03-16 |
| pipeline_actions | Central orchestration table — action-based pipeline, not linear; round_number + phase for multi-round | 2026-03-16 |
| JSONB embeddings | Vectors stored as JSONB (pgvector migration path available later) | 2026-03-16 |
| Temp interviewers | invite_links + is_temp user flag; 256-bit token, time-limited | 2026-03-16 |
| Batch2 conflict resolution | job_embeddings=JSONB, blacklist_logs=candidate_id+action, email_templates=code+type | 2026-03-16 |
| Design token flags | Flagged items (weights, spacing, radius, shadows) approved for dev — refine at UAT | 2026-03-16 |

## Notion DB Property Reference

Team Directory DB: `Agent` (title), `Code` (rich_text), `Role` (select), `Phase` (multi_select), `Status` (select), `Deliverables` (rich_text)
Sprint Roadmap DB: `Sprint` (title), `Focus` (rich_text), `Modules` (rich_text), `Tasks` (number), `Points` (number), `Capacity` (rich_text), `Status` (select)
Sprint 1 Kanban DB: `Task` (title), `Status` (select: To Do/In Progress/In Review/Done), `Agent` (select), `Tag` (select), `Points` (number), `ID` (rich_text)
API Registry DB: `Endpoint` (title), `Method` (select), `Module` (select), `Auth` (select), `Status` (select), `Description` (rich_text)
Activity Log DB: `Event` (title), `Phase` (select), `Agent` (select), `Date` (date), `Outcome` (select), `Details` (rich_text)

## Gate Failure Patterns

_None recorded yet._

## Sprint Velocity Benchmarks

_None recorded yet — Phase 3 PM will establish._

## Bug Patterns by Module

_None recorded yet._
