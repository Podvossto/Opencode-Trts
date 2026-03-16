===PHASE REPORT===
Phase   : P1
Agent   : product-owner-agile (inline)
Task ID : ATS-BACKLOG-001 / Batch 1 (M01–M05)
Status  : COMPLETE

[OUTPUT]

## BATCH 1 SUMMARY — M01–M05

Total Stories : 24
Total Points  : 86
MoSCoW        : Must=24 Should=0 Could=0 Won't=0
Stories ≥13pts: NONE — all ≤8pts ✅

---

## M01 — Admin Auth & User Management (7 stories | 18 pts)

US-M01-01 | Admin Login                          | 3 pts | Must
US-M01-02 | Forced Password Change on First Login | 3 pts | Must
US-M01-03 | Session Timeout 30min                | 2 pts | Must
US-M01-04 | List Users                           | 2 pts | Must
US-M01-05 | Create New User                      | 3 pts | Must
US-M01-06 | Edit and Delete User                 | 3 pts | Must
US-M01-07 | Admin Resets User Password           | 2 pts | Must

OUT OF SCOPE: SSO/OAuth, granular permissions, MFA beyond OTP, admin audit log

---

## M02 — HR Auth & Password Reset (2 stories | 7 pts)

US-M02-01 | HR/Supervisor Login                 | 2 pts | Must
US-M02-02 | HR Self-Service OTP Password Reset   | 5 pts | Must

OUT OF SCOPE: Social login, account lockout policy, Admin-triggered HR OTP path (M01-07)

---

## M03 — Applicant Portal — Job Application Submission (5 stories | 19 pts)

US-M03-01 | Applicant Accesses Job via URL       | 3 pts | Must
US-M03-02 | Submit Application Form              | 5 pts | Must
US-M03-03 | Block Blacklisted Candidate          | 3 pts | Must
US-M03-04 | Auto-trigger OCR on Resume           | 3 pts | Must
US-M03-05 | Online Test Re-openable + Auto-save  | 5 pts | Must

OUT OF SCOPE: Applicant account/login, form draft-save, status tracking page, file format validation rules

---

## M04 — Job Posting & Pipeline Dashboard (6 stories | 26 pts)

US-M04-01 | Create Job Posting with Hard Skill Weights | 5 pts | Must
US-M04-02 | Job Scheduling Auto-Publish / Auto-Close   | 5 pts | Must
US-M04-03 | Edit, Filter, Manage Job Postings          | 3 pts | Must
US-M04-04 | Pipeline Dashboard 5-Stage Milestone View  | 5 pts | Must
US-M04-05 | Dashboard Stats (Total/Remaining/Dropped/Hired) | 3 pts | Must
US-M04-06 | Email Trigger from Pipeline Dashboard      | 5 pts | Must

OUT OF SCOPE: External job board posting, job approval workflow, advanced analytics, duplicate/clone job

---

## M05 — Application Form Builder (4 stories | 16 pts)

US-M05-01 | Build Form with Standard Fields           | 5 pts | Must
US-M05-02 | Add Custom Question Types                 | 5 pts | Must
US-M05-03 | Publish Form as JSON to Database          | 3 pts | Must
US-M05-04 | Edit and Delete Forms with Reference Check| 3 pts | Must

OUT OF SCOPE: Form versioning, form templates/library, conditional field logic, form analytics

---

[HANDOFF NOTES]
For P2 (SA):
- All 24 stories are Must; no splitting required
- M01+M02 share auth mechanism — architect as single auth module
- M03 OCR (US-M03-04) requires async queue — FastAPI AI service endpoint needed
- M04 auto-publish/close (US-M04-02) requires a background scheduler (cron job or Go worker)
- M05 Form JSON schema must support: Text, Textarea, Dropdown, Checkbox, Radio, File Upload, Date Picker, MCQ, Short Answer, Essay, True/False, Rating Scale, Image/Video/Link media fields
- Recommended sprint split (for PM reference): Sprint 1 = M01+M02+M03 (28 pts auth+portal foundation); Sprint 2 = M04+M05 (42 pts job+form)

===END REPORT===
