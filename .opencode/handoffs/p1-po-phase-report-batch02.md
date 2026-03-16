===PHASE REPORT===
Phase   : P1
Agent   : product-owner-agile
Task ID : ATS-BACKLOG-002 / Batch 2 (M06–M15)
Status  : COMPLETE

[OUTPUT]

---

## M06 — Resume Review (AI Similarity + LLM Summary)

---

### **[US-M06-01] 3-Panel Resume Review Layout**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want a 3-panel Resume Review layout so that I can see the applicant list, resume/documents, and AI analysis simultaneously without switching screens.

**Gherkin AC:**
```
Feature: Resume Review 3-Panel Layout
  Scenario: HR opens Resume Review for a job
    Given HR navigates to Resume Review for an active job
    Then the screen displays 3 panels:
      | Panel        | Content                                              |
      | Left         | Waiting list with applicant badges, Dept filter, Include Starred checkbox |
      | Centre       | Resume viewer + uploaded documents                   |
      | Right-top    | Applicant personal info + LLM Summary                |
      | Right-bottom | Position panel with Overall Score ranked list        |
    And selecting an applicant in the left panel updates centre and right panels

  Scenario: HR filters waiting list by department
    Given the waiting list contains applicants from multiple departments
    When HR selects a department in the Dept filter
    Then only applicants from that department appear in the waiting list

  Scenario: HR toggles Include Starred
    Given there are starred candidates matching the job
    When HR checks the Include Starred checkbox
    Then starred candidates appear in the waiting list with a star badge
    When HR unchecks Include Starred
    Then starred candidates are hidden from the waiting list
```

Out-of-Scope (M06): Applicant self-service view, multi-job simultaneous review, bulk interest/reject.

---

### **[US-M06-02] Overall Score Calculation with Weight Popup**
**3 pts | Must Have**
Tags: `frontend`, `backend`, `ai-service`

As an HR user I want the Overall Score to be calculated as a weighted average of JD similarity and Hard Skill similarity using values from the Weight Popup so that scoring reflects the job's priorities.

**Gherkin AC:**
```
Feature: Overall Score Weighted Calculation
  Scenario: HR views pre-computed scores on load
    Given HR opens Resume Review for a job with published weights
    Then each applicant in the right-bottom Position panel shows an Overall Score
    And scores are ranked in descending order

  Scenario: Score formula uses Weight Popup values
    Given a job has JD weight = 40% and Hard Skill weight = 60%
    And an applicant has JD similarity = 0.80 and Hard Skill similarity = 0.70
    Then the applicant's Overall Score = (0.80 * 0.40) + (0.70 * 0.60) = 0.74

  Scenario: System pre-computes score for Next Applicant on load
    Given HR is viewing applicant A
    Then the system pre-computes the similarity score for the next applicant in the list
    So that switching to the next applicant shows the score immediately without delay
```

Out-of-Scope (M06): Manual score override, score history, score visibility to applicant.

---

### **[US-M06-03] HR Marks Up to 3 Suitable Positions**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to mark up to 3 suitable positions for an applicant during Resume Review so that I can record position fit for reference without creating new applications.

**Gherkin AC:**
```
Feature: Position Marking in Resume Review
  Scenario: HR marks a suitable position
    Given HR is reviewing an applicant and clicks Interest
    Then HR is prompted to select up to 3 suitable positions from active job list
    And the selections are saved as position_marks on the applicant record
    And no new Application record is created

  Scenario: HR tries to mark more than 3 positions
    Given HR has already marked 3 positions for an applicant
    When HR attempts to select a 4th position
    Then the system prevents selection and shows a warning "Maximum 3 positions allowed"

  Scenario: Position marks are reference-only
    Given HR has marked positions for an applicant
    Then the marks appear in the applicant's detail view for HR reference
    And they do not trigger any pipeline action or notification
```

Out-of-Scope (M06): Auto-transfer based on marks, applicant visibility of marks.

---

### **[US-M06-04] Resume Review Action Buttons (Interest / Reject / Star / End Review / Next)**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want clear action buttons in Resume Review so that I can decisively move each applicant to the correct next state.

**Gherkin AC:**
```
Feature: Resume Review Action Buttons
  Scenario: HR clicks Interest
    Given HR is reviewing an applicant with status Pending
    When HR clicks Interest
    Then the applicant's status changes to In Progress
    And the applicant disappears from the Pending waiting list
    And HR is prompted to mark up to 3 suitable positions (US-M06-03)

  Scenario: HR clicks Reject
    Given HR is reviewing an applicant
    When HR clicks Reject
    Then the applicant's status changes to Dropped
    And the system automatically sends rejection email SYS02 to the applicant
    And the applicant is removed from the waiting list

  Scenario: HR clicks Star
    Given HR is reviewing an applicant
    When HR clicks Star
    Then the applicant's status changes to Starred
    And the applicant is added to the Talent Pool
    And the applicant is auto-rejected from the current job (Dropped)
    And NO rejection email is sent to the applicant
    And the applicant disappears from the current job's waiting list

  Scenario: HR clicks End Review
    Given HR is in Resume Review
    When HR clicks End Review
    Then HR is returned to the job dashboard without further changes

  Scenario: HR clicks Next
    Given HR is reviewing applicant A
    When HR clicks Next
    Then the system moves to the next applicant in the waiting list
    And the pre-computed score (US-M06-02) is displayed immediately

  Scenario: HR clicks Retry OCR when OCR failed
    Given an applicant's resume failed OCR processing
    Then the LLM Summary panel shows a Retry OCR button instead of the summary
    When HR clicks Retry OCR
    Then the system re-queues the OCR job for that resume
    And the panel shows a loading indicator until OCR completes
```

Out-of-Scope (M06): Bulk actions, undo/reverse of Reject or Star.

---

### **[US-M06-05] Auto-Compute Similarity for Starred Candidates on New Job Publish**
**3 pts | Must Have**
Tags: `backend`, `ai-service`

As the system I want to automatically compute similarity scores for all Starred candidates whenever a new job is published so that HR can immediately see scored talent pool applicants in Resume Review.

**Gherkin AC:**
```
Feature: Auto-Compute Similarity on Job Publish
  Scenario: New job is published with starred candidates in pool
    Given there are 5 starred candidates in the Talent Pool
    When HR publishes a new job posting
    Then the system asynchronously computes JD and Hard Skill similarity for all 5 starred candidates against the new job
    And when computation is complete, each starred candidate has an Overall Score for that job

  Scenario: Starred candidates appear in new job's Resume Review
    Given similarity has been computed for starred candidates against job J
    When HR opens Resume Review for job J and enables Include Starred
    Then starred candidates appear in the waiting list with their star badge and Overall Score

  Scenario: Computation does not block job publish
    Given HR clicks Publish on a new job
    Then the job becomes Active immediately
    And similarity computation runs as a background task without delaying the publish action
```

Out-of-Scope (M06): Real-time similarity re-compute on weight change, compute for non-starred candidates.

---

## M07 — Action Pipeline — Test (Multi-round)

---

### **[US-M07-01] HR Initiates Action Test with Round Badge**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to initiate an Action Test for a group of applicants and see a round badge per applicant so that I can track how many test rounds each applicant has completed.

**Gherkin AC:**
```
Feature: Initiate Action Test
  Scenario: HR initiates first test round
    Given a job has applicants with status In Progress
    When HR clicks Action Test on the pipeline dashboard
    Then the system shows all applicants with status In Progress who are not Dropped or Hired
    And each applicant shows a badge "Round 1"

  Scenario: HR initiates a subsequent test round
    Given applicants have already completed Round 1
    When HR clicks Action Test again
    Then the system shows eligible applicants with badge "Round 2"
    And the badge count increments automatically per applicant per round

  Scenario: Round badge is per-applicant
    Given applicant A is in Round 2 and applicant B is in Round 1
    When HR views the Action Test applicant list
    Then applicant A shows badge "Round 2" and applicant B shows badge "Round 1"
```

Out-of-Scope (M07): Transfer from Test stage, test round history export.

---

### **[US-M07-02] HR Selects Applicants and Assigns Test**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to select which applicants join each test round and assign the appropriate test so that only chosen candidates receive the test.

**Gherkin AC:**
```
Feature: Applicant Selection and Test Assignment
  Scenario: First round uses default test linked to job
    Given a job has a default test configured (FRH06)
    When HR opens Action Test Round 1
    Then the default test is pre-selected for the round
    And HR can select a subset of eligible applicants to include

  Scenario: Subsequent rounds allow HR to choose from library
    Given HR is setting up Round 2 or later
    When HR opens Action Test configuration
    Then HR can select any test from the General Test library
    And HR confirms the applicant selection before sending

  Scenario: HR confirms selection before sending test
    Given HR has selected applicants and a test
    When HR clicks Confirm
    Then the system sends the test link to selected applicants via email
    And non-selected eligible applicants remain in the list for future rounds
```

Out-of-Scope (M07): Supervisor assigning General Test (Supervisor assigns Technical Test only).

---

### **[US-M07-03] Online Test Delivery with Anti-Cheat**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As the system I want to deliver online tests with anti-cheat controls and timer so that test integrity is maintained.

**Gherkin AC:**
```
Feature: Online Test Delivery and Anti-Cheat
  Scenario: Applicant receives and opens test link
    Given the system has sent a test link to applicant A
    When applicant A clicks the link
    Then the test loads with a timer displayed (per-test or per-question as configured)
    And the test is accessible without login

  Scenario: Copy-paste is blocked
    Given applicant A is taking an online test
    When applicant A attempts to copy text from or paste text into the test
    Then the action is blocked by the system

  Scenario: Tab-switch detection
    Given applicant A is taking an online test
    When applicant A switches to another browser tab or window
    Then the system records a tab-switch event
    And optionally warns the applicant that tab-switching was detected

  Scenario: Questions displayed in random order
    Given a test has random order enabled
    When two different applicants open the same test
    Then the question order is different for each applicant

  Scenario: Test link is re-openable until Submit or expiry
    Given applicant A has started but not submitted the test
    When applicant A closes and re-opens the test link
    Then the test resumes with previously auto-saved answers loaded
    When the test has been submitted or expired
    Then opening the link shows a closed/submitted message and no further editing is allowed
```

Out-of-Scope (M07): Proctoring video, IP restriction, browser lockdown.

---

### **[US-M07-04] MCQ Auto-Grading and Manual Grading Inbox**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As HR and Supervisor I want MCQ questions to be auto-graded and essay/short-answer questions to be graded via a dedicated inbox so that test results are processed efficiently.

**Gherkin AC:**
```
Feature: Test Grading
  Scenario: MCQ questions are auto-graded on submit
    Given an applicant submits a test containing MCQ questions
    Then the system immediately calculates the MCQ score
    And the score is visible to HR in real-time (Phase 1 view)

  Scenario: HR grades essay/short-answer for General Test
    Given a General Test has essay or short-answer questions
    And an applicant has submitted the test
    Then HR sees the ungraded questions in the Grading Inbox with a badge count
    When HR assigns scores per question and submits
    Then the total score is updated

  Scenario: Supervisor grades essay/short-answer for Technical Test
    Given a Technical Test was created by Supervisor S
    And an applicant has submitted the test
    Then Supervisor S sees the ungraded questions in their Grading Inbox
    When Supervisor S grades and submits
    Then the system sends notification FRNOT05 to HR

  Scenario: HR sees scores real-time during Phase 1
    Given a test round is in Phase 1 (Grading)
    When any grader submits scores
    Then HR sees the updated score for that applicant immediately without page refresh
    And HR can see which applicants still have ungraded manual questions

  Scenario: Grading inbox badge count
    Given there are 3 applicants with ungraded questions
    Then the Grading Inbox shows a badge count of 3
    When all questions for an applicant are graded
    Then the badge count decreases by 1
```

Out-of-Scope (M07): Partial score saving by grader (submit is atomic), grade dispute workflow.

---

### **[US-M07-05] HR Opens Phase 2 Selection and Leaderboard**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to manually open Phase 2 Selection after grading and view a live leaderboard so that I can make informed pass/drop decisions.

**Gherkin AC:**
```
Feature: Phase 2 Test Selection
  Scenario: HR opens Phase 2
    Given a test round is in Phase 1 and HR is ready to proceed
    When HR clicks "Start Selection"
    Then the view transitions to Phase 2 Selection
    And a leaderboard appears sorted by total score descending

  Scenario: Leaderboard updates live
    Given Phase 2 is active and a grader submits a late score
    Then the leaderboard score updates immediately without a page refresh
    And HR can see which applicants still have incomplete grading

  Scenario: HR makes pass/drop decisions in Phase 2
    Given HR is viewing the Phase 2 leaderboard
    When HR selects an applicant and chooses Test Again, Interview, Hire, or Drop
    Then the applicant's next action is recorded
    And the decision cannot be changed once confirmed (no undo)

  Scenario: Drop triggers rejection email
    Given HR drops an applicant in Phase 2
    Then the system automatically sends rejection email SYS02 to the applicant
```

Out-of-Scope (M07): Transfer from Test Phase 2 (Transfer is only available after Interview).

---

### **[US-M07-06] HR and Supervisor CRUD Test Library**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As HR and Supervisor I want to create, edit, and delete tests in the test library so that the right tests are always available for pipeline actions.

**Gherkin AC:**
```
Feature: Test Library Management
  Scenario: HR creates a General Test
    Given HR is in the Test Library
    When HR creates a new test with title, questions, and time settings
    Then the test is saved to the General Test library and available for Action Test selection

  Scenario: Supervisor creates a Technical Test
    Given Supervisor S is in the Test Library
    When Supervisor S creates a Technical Test with MCQ and essay questions
    Then the test is saved to Supervisor S's Technical Test library

  Scenario: HR/Supervisor edits an existing test
    Given a test exists in the library
    When HR or the owning Supervisor edits the test
    Then changes are saved and take effect for future test rounds (not active rounds)

  Scenario: HR/Supervisor deletes a test with reference check
    Given a test is referenced as the default test for an active job
    When HR attempts to delete the test
    Then the system blocks deletion and shows "Test is in use by active job(s)"
    When a test has no active references
    Then deletion proceeds after confirmation

  Scenario: HR links default test to job during job creation
    Given HR is creating a job posting
    When HR selects a default General Test and/or Technical Test from the library
    Then the selected tests are linked as default for Round 1 of Action Test for that job
```

Out-of-Scope (M07): Test versioning, test cloning, shared test library between HR and Supervisor.

---

## M08 — Action Pipeline — Interview (Multi-round, Scoring Criteria, Invite Links)

---

### **[US-M08-01] HR Initiates Action Interview with Round Badge**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to initiate an Action Interview for a group of applicants and see a round badge per applicant so that I can track how many interview rounds each applicant has completed.

**Gherkin AC:**
```
Feature: Initiate Action Interview
  Scenario: HR initiates first interview round
    Given a job has applicants with status In Progress
    When HR clicks Action Interview on the pipeline dashboard
    Then the system shows all eligible applicants who are not Dropped or Hired
    And each applicant shows a badge "Round 1"

  Scenario: HR initiates a subsequent interview round
    Given applicants have completed Round 1 interview
    When HR clicks Action Interview again
    Then eligible applicants appear with badge "Round 2"
    And the badge count increments automatically per applicant per round

  Scenario: Round badge is per-applicant
    Given applicant A is in Round 2 and applicant B is in Round 1
    When HR views the Action Interview applicant list
    Then applicant A shows badge "Round 2" and applicant B shows badge "Round 1"
```

Out-of-Scope (M08): Transfer initiation (Transfer is a separate action from Interview Phase 2), test round history export.

---

### **[US-M08-02] HR Defines Scoring Criteria per Round**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to define scoring criteria (criteria name + max score) for each interview round so that all interviewers evaluate applicants on the same structured basis.

**Gherkin AC:**
```
Feature: Interview Scoring Criteria Definition
  Scenario: HR creates scoring criteria for a round
    Given HR is setting up an interview round
    When HR adds criteria items (e.g. "Communication Skills", max score 10)
    Then each criterion is saved with a name and max_score for that round
    And interviewers will see these criteria when grading

  Scenario: Supervisor creates criteria only for assigned rounds
    Given Supervisor S is assigned to interview round R
    When Supervisor S adds a scoring criterion to round R
    Then the criterion is saved for that round
    When Supervisor S attempts to add criteria to a round they are not assigned to
    Then the system blocks the action with "Access denied"

  Scenario: Criteria are locked once grading begins
    Given at least one interviewer has submitted scores for a round
    When HR attempts to edit or delete a criterion for that round
    Then the system blocks the change and shows "Grading in progress — criteria cannot be changed"
```

Out-of-Scope (M08): Criteria templates, global criteria library, criteria reuse across rounds.

---

### **[US-M08-03] HR Assigns Interviewers — Supervisor and Temp Interviewers**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to assign both regular Supervisors and temporary interviewers (via invite link) to an interview round so that the right people evaluate each applicant.

**Gherkin AC:**
```
Feature: Interviewer Assignment
  Scenario: HR assigns a Supervisor as interviewer
    Given HR is configuring an interview round
    When HR selects a Supervisor from the active user list
    Then the Supervisor is added as an interviewer for that round
    And the Supervisor sees the assigned applicants and criteria in their grading view

  Scenario: HR generates a temp interviewer invite link
    Given HR needs to invite an external or temporary interviewer
    When HR clicks "Generate Invite Link" for a round
    Then the system generates a 256-bit token invite URL
    And HR can copy and share the link with the temp interviewer

  Scenario: Temp account auto-created on first link click
    Given a temp invite link has been generated
    When the temp interviewer clicks the link for the first time
    Then the system auto-creates a temporary account scoped to the assigned applicants and criteria for that round
    And the temp interviewer can grade without a pre-existing account

  Scenario: Temp account is auto-deactivated after round
    Given a temp interviewer was invited for round R
    When round R is closed or the invite link expires
    Then the system automatically deactivates the temp account
    And the temp interviewer can no longer access the grading view

  Scenario: HR revokes a temp invite link
    Given a temp invite link is active
    When HR clicks Revoke on the invite link
    Then the link is immediately invalidated
    And the temp account (if created) is deactivated
    And any scores already submitted by that temp interviewer are retained
```

Out-of-Scope (M08): Temp interviewers creating scoring criteria, temp interviewers viewing other rounds.

---

### **[US-M08-04] HR Views and Manages Invite Links List**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to see all active and expired invite links for an interview round so that I can track and revoke access as needed.

**Gherkin AC:**
```
Feature: Invite Links Management
  Scenario: HR views invite links for a round
    Given HR is in the interview round management view
    When HR opens the Invite Links panel
    Then HR sees a list of all generated links with:
      | Column      | Content                              |
      | Link token  | Masked URL (copy button)             |
      | Status      | Active / Revoked / Expired           |
      | Created at  | Timestamp                            |
      | Temp account| Email or "not yet activated"         |

  Scenario: HR revokes an active link from the list
    Given the invite links list shows an Active link
    When HR clicks Revoke
    Then the link status changes to Revoked immediately
    And the associated temp account is deactivated

  Scenario: Expired links are shown but not actionable
    Given a link has passed its expiry time
    Then it appears in the list with status "Expired"
    And the Revoke button is disabled for expired links
```

Out-of-Scope (M08): Bulk revoke, resending invite links via email automatically.

---

### **[US-M08-05] Phase 1 Interview Grading — Real-Time Per-Interviewer Scores**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to see per-interviewer scores in real-time during Phase 1 grading and know who has not yet submitted so that I can monitor interview progress before opening Phase 2.

**Gherkin AC:**
```
Feature: Phase 1 Interview Grading
  Scenario: Interviewer grades an applicant per criterion
    Given an interviewer is in the grading view for round R
    When the interviewer enters a score per criterion and adds an optional note
    Then the scores are saved per criterion per applicant per interviewer

  Scenario: Interviewer submits final scores
    Given an interviewer has entered scores for all assigned applicants
    When the interviewer clicks Submit
    Then the submission is recorded as complete for that interviewer
    And the scores become visible to HR

  Scenario: HR sees per-interviewer scores in real-time
    Given an interviewer has submitted scores
    Then HR sees that interviewer's scores for each applicant immediately without page refresh
    And scores are shown per interviewer (not aggregated)

  Scenario: HR sees who has not submitted
    Given round R has 3 assigned interviewers
    And only 2 have submitted
    Then HR sees a visual indicator showing interviewer 3 as "Pending"
    And HR can see which specific applicants are missing scores from interviewer 3
```

Out-of-Scope (M08): HR editing an interviewer's submitted scores, score dispute workflow.

---

### **[US-M08-06] Phase 2 Interview Selection — Leaderboard and Final Decisions**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to open Phase 2 Selection after all grading is complete and view an avg-of-averages leaderboard with per-interviewer comparison so that I can make final decisions.

**Gherkin AC:**
```
Feature: Phase 2 Interview Selection
  Scenario: HR opens Phase 2 Selection
    Given all interviewers have submitted scores for round R
    When HR clicks "Start Selection"
    Then the view transitions to Phase 2 Selection
    And a leaderboard appears ranked by avg-of-averages score (average of each interviewer's average) descending

  Scenario: Leaderboard shows per-interviewer comparison
    Given HR is in Phase 2 Selection
    When HR expands an applicant row
    Then HR sees each interviewer's individual average score alongside the combined average

  Scenario: Leaderboard updates live for late submissions
    Given Phase 2 is open and a late interviewer submits scores
    Then the leaderboard recalculates and reorders immediately

  Scenario: HR makes a final decision (no undo)
    Given HR is viewing Phase 2 leaderboard
    When HR selects an applicant and chooses Test, Interview, Transfer, Hire, or Drop
    Then the decision is recorded immediately
    And the decision cannot be changed once confirmed

  Scenario: Drop triggers rejection email
    Given HR drops an applicant in Phase 2
    Then the system automatically sends rejection email SYS02 to the applicant
```

Out-of-Scope (M08): Auto-advance to Phase 2, partial decision saves, bulk decisions.

---

### **[US-M08-07] HR Sends Interview Invitation Email (TPL05)**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to send interview invitation emails to selected applicants using template TPL05 so that applicants receive timely, correctly formatted interview notices.

**Gherkin AC:**
```
Feature: Interview Invitation Email
  Scenario: HR sends invitation to selected applicants
    Given HR has selected applicants for an interview round
    When HR clicks "Send Invitation"
    Then the system sends email using TPL05 to each selected applicant
    And a delivery timestamp is recorded per applicant

  Scenario: System blocks send if TPL05 does not exist
    Given TPL05 has not been configured in Setup
    When HR attempts to send interview invitations
    Then the system blocks the action and shows "Email template TPL05 is required — configure it in Setup"

  Scenario: Sent invitations are logged
    Given invitation emails have been sent
    Then HR can see a "Sent" status per applicant in the interview round view
    With the timestamp of when the invitation was sent
```

Out-of-Scope (M08): Applicant reply tracking within the ATS, bulk re-send, scheduling via calendar integration.

---

## M09 — Cross-Position Transfer

---

### **[US-M09-01] HR Initiates Transfer During or After Interview**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to transfer an applicant to a different position during or after Interview (Phase 1 or Phase 2) so that I can redirect qualified candidates without losing their history.

**Gherkin AC:**
```
Feature: Cross-Position Transfer
  Scenario: HR transfers an applicant from Interview Phase 1
    Given applicant A is in Interview Phase 1 for position P1
    When HR selects applicant A and clicks Transfer
    Then HR is prompted to choose a target position from active job list
    When HR confirms the target position P2
    Then the application on P1 is closed and the headcount slot on P1 is returned
    And a new Application is created on P2 carrying over applicant history
    And applicant A is now visible in P2's pipeline

  Scenario: HR transfers an applicant from Interview Phase 2
    Given applicant A is in Interview Phase 2 for position P1
    When HR selects Transfer as the next action for applicant A
    Then the same transfer flow executes: P1 closed, slot returned, P2 Application created with history

  Scenario: Transfer is blocked if source is not Interview stage
    Given applicant A is in Test Phase 2 (not Interview)
    When HR attempts to transfer applicant A
    Then the system blocks the action and shows "Transfer is only available during or after Interview"

  Scenario: Applicant is unaware of transfer until Offer
    Given applicant A has been transferred from P1 to P2
    Then no notification is sent to applicant A
    And applicant A is not informed of the transfer until an Offer is issued
```

Out-of-Scope (M09): Transfer to a closed or archived position without a warning, applicant-initiated transfer.

---

### **[US-M09-02] Transfer History Carries Over and Unlimited Transfers Supported**
**3 pts | Must Have**
Tags: `backend`

As the system I want to preserve full application history across transfers and allow unlimited sequential transfers so that HR always has a complete view of the applicant's journey.

**Gherkin AC:**
```
Feature: Transfer History and Chaining
  Scenario: History is carried over on transfer
    Given applicant A has Test scores and Interview scores on position P1
    When HR transfers applicant A to P2
    Then the new Application on P2 includes a reference to all prior history from P1
    And HR can see the full history (P1 test scores, P1 interview scores) in applicant A's detail view on P2

  Scenario: Unlimited sequential transfers
    Given applicant A has been transferred from P1 to P2
    When HR transfers applicant A from P2 to P3
    Then a third Application is created on P3 with history from both P1 and P2 carried over
    And this can repeat without a system-imposed transfer limit

  Scenario: Transfer chain is visible as timeline
    Given applicant A has been transferred P1→P2→P3
    When HR views applicant A's Candidate Detail
    Then the timeline shows each transfer event in order with source position, target position, and timestamp
```

Out-of-Scope (M09): Reverting a transfer, merging duplicate applicant records from transfers.

---

### **[US-M09-03] System Warns HR if Target Position is Closed**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want the system to warn me before confirming a transfer to a closed position so that I do not inadvertently transfer an applicant to a position that is no longer accepting candidates.

**Gherkin AC:**
```
Feature: Transfer Target Position Validation
  Scenario: HR selects a closed position as transfer target
    Given applicant A is eligible for transfer
    When HR selects position P2 as the transfer target
    And position P2 has status Closed
    Then the system shows a warning: "Position P2 is currently Closed. Transferring will reopen the slot. Do you want to continue?"
    And HR must explicitly confirm before the transfer proceeds

  Scenario: HR cancels after warning
    Given the warning is shown for a closed target position
    When HR clicks Cancel
    Then no transfer is created and applicant A remains in the original position unchanged

  Scenario: Transfer to active position shows no warning
    Given HR selects an Active position as the transfer target
    Then no warning is shown and the transfer confirmation dialog proceeds normally
```

Out-of-Scope (M09): Auto-reopening a closed position without HR confirmation, blocking transfer to closed positions entirely.

---

### **[US-M09-04] Block Sending Test from Wrong Position After Transfer**
**2 pts | Must Have**
Tags: `backend`

As the system I want to block any test delivery to an applicant from a position they have been transferred away from so that transferred applicants do not receive tests from their old position.

**Gherkin AC:**
```
Feature: Block Test Send After Transfer
  Scenario: HR attempts to send a test from the closed source position
    Given applicant A was transferred from P1 to P2
    And the Application on P1 is now Closed
    When HR (or any system action) tries to send a test to applicant A from P1
    Then the system blocks the action and shows "Applicant A has been transferred — tests cannot be sent from this position"

  Scenario: Test can be sent from the new position
    Given applicant A is now in P2's pipeline
    When HR sends a test to applicant A from P2
    Then the test is delivered normally

  Scenario: Closed application is excluded from all action lists
    Given applicant A's P1 Application is Closed
    When HR opens Action Test or Action Interview on P1
    Then applicant A does not appear in any eligible applicant list for P1
```

Out-of-Scope (M09): Retroactive cancellation of already-sent tests from the source position.

---

## M10 — Talent Pool / Star Feature

---

### **[US-M10-01] Auto-Compute Similarity for Starred Candidates Against All New Jobs**
**3 pts | Must Have**
Tags: `backend`, `ai-service`

As the system I want to automatically compute similarity scores for all Starred candidates against every new job when it is published so that HR can immediately see scored talent pool candidates in any new job's Resume Review.

**Gherkin AC:**
```
Feature: Auto-Compute Similarity on Job Publish (Talent Pool)
  Scenario: System triggers computation for all starred candidates on new job publish
    Given there are N starred candidates in the Talent Pool
    When a new job J is published
    Then the system enqueues similarity computation jobs for all N starred candidates against job J
    And computation runs asynchronously without blocking the publish action

  Scenario: Scored starred candidates appear in the new job's Resume Review
    Given similarity computation has completed for starred candidates against job J
    When HR opens Resume Review for job J and enables the Include Starred checkbox
    Then all starred candidates appear in the waiting list with:
      | Field        | Value                         |
      | Star badge   | Visible                       |
      | Overall Score| Computed JD+HardSkill weighted|
      | Ranked       | By Overall Score descending   |

  Scenario: Starred candidate with no similarity result yet shows pending state
    Given similarity computation is still running for candidate C against job J
    When HR opens Resume Review for job J with Include Starred enabled
    Then candidate C appears with a "Score pending" indicator instead of a numeric score
    And the score updates automatically once computation completes
```

Out-of-Scope (M10): Star button behaviour (belongs to M06), re-compute on weight change, compute for non-starred candidates.

---

### **[US-M10-02] Starred Candidates Removed from Pool on Blacklist**
**2 pts | Must Have**
Tags: `backend`

As the system I want to automatically remove a candidate from the Talent Pool when they are blacklisted so that blacklisted individuals never appear in the starred candidate list for future jobs.

**Gherkin AC:**
```
Feature: Talent Pool Blacklist Removal
  Scenario: Blacklisted candidate is removed from Talent Pool
    Given candidate C is in the Talent Pool with Starred status
    When HR blacklists candidate C
    Then candidate C's Starred status is removed
    And candidate C no longer appears in the Talent Pool
    And candidate C no longer appears in Resume Review waiting lists for any job (Include Starred enabled)

  Scenario: Previously computed similarity scores are retained but not surfaced
    Given candidate C had similarity scores computed for jobs J1 and J2 before blacklisting
    When HR blacklists candidate C
    Then the similarity score records are retained in the database for audit purposes
    But candidate C does not appear in any job's Include Starred list

  Scenario: Unblacklisting does not auto-restore Starred status
    Given candidate C was blacklisted and removed from Talent Pool
    When HR unblacklists candidate C
    Then candidate C's status is restored to normal (not Starred)
    And candidate C is NOT automatically re-added to the Talent Pool
    And HR must explicitly re-star the candidate to add them back to the pool
```

Out-of-Scope (M10): Manual bulk removal from Talent Pool, Talent Pool capacity limits, applicant visibility of Starred status.

---

## M11 — Candidate Management (List, Detail, OCR, Blacklist)

---

### **[US-M11-01] Candidate List with Search and Filter**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to view a searchable, filterable list of all candidates so that I can quickly find and access any candidate's record.

**Gherkin AC:**
```
Feature: Candidate List
  Scenario: HR views the candidate list
    Given HR navigates to Candidate Management
    Then HR sees a paginated list of all candidates with columns:
      | Column        | Content                         |
      | Name          | Full name                       |
      | Email         | Email address                   |
      | App Count     | Number of applications          |
      | Latest Status | Most recent application status  |

  Scenario: HR searches by name or email
    Given the candidate list is displayed
    When HR types a name or email fragment in the search box
    Then the list filters to show only candidates whose name or email contains the search text
    And the filter applies without requiring a page reload

  Scenario: HR filters by status
    Given the candidate list is displayed
    When HR selects a status from the filter dropdown (e.g. In Progress, Dropped, Hired)
    Then only candidates whose latest application matches that status are shown

  Scenario: HR filters by position
    Given the candidate list is displayed
    When HR selects a job position from the position filter
    Then only candidates who have applied for that position are shown
```

Out-of-Scope (M11): Bulk actions on candidate list, export to CSV, candidate merge.

---

### **[US-M11-02] Candidate Detail — Full Application History and Documents**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to view a candidate's complete detail page including personal info, all application history, test/interview scores, documents, and a transfer timeline so that I have full context for every hiring decision.

**Gherkin AC:**
```
Feature: Candidate Detail View
  Scenario: HR opens a candidate detail page
    Given HR clicks on a candidate in the candidate list
    Then HR sees the candidate detail page with:
      | Section              | Content                                              |
      | Personal Info        | Name, email, phone, education, experience            |
      | Application History  | All applications with position, status, dates        |
      | Test Scores          | Scores per round per position                        |
      | Interview Scores     | Per-interviewer scores per round per position        |
      | Documents            | Resume and uploaded files with download buttons      |
      | Timeline             | Chronological events including transfer history      |

  Scenario: HR edits candidate name
    Given HR is on the candidate detail page
    When HR clicks Edit on the name field and enters a corrected name
    Then the candidate's name is updated in the system
    And the change is logged in the timeline

  Scenario: HR downloads a candidate document
    Given the candidate has uploaded a resume and other documents
    When HR clicks the download button on a document
    Then the file downloads to HR's device

  Scenario: Timeline shows transfer history
    Given candidate C has been transferred from P1 to P2 to P3
    When HR views the timeline on the candidate detail page
    Then each transfer event appears with source position, target position, and timestamp in chronological order
```

Out-of-Scope (M11): HR editing fields other than name, applicant self-editing their profile.

---

### **[US-M11-03] OCR Auto-Trigger on Resume Upload**
**3 pts | Must Have**
Tags: `backend`, `ai-service`

As the system I want to automatically trigger OCR processing on a candidate's resume when it is uploaded so that the resume text and vectors are available for AI similarity scoring without HR intervention.

**Gherkin AC:**
```
Feature: OCR Auto-Trigger
  Scenario: OCR triggered automatically after resume upload
    Given a candidate submits an application with a PDF/DOCX/image resume
    When the application is saved
    Then the system enqueues an OCR job for the resume to the AI service
    And the OCR job runs asynchronously without blocking the submission confirmation

  Scenario: OCR success stores text and vector
    Given an OCR job completes successfully
    Then the extracted text is stored on the candidate record
    And a vector embedding is generated and stored for similarity computation

  Scenario: OCR failure shows Retry OCR in candidate detail
    Given an OCR job has failed for candidate C
    When HR opens candidate C's detail page
    Then the Documents section shows a "Retry OCR" button next to the failed resume
    When HR clicks Retry OCR
    Then the system re-enqueues the OCR job
    And the button changes to a loading indicator until the job completes or fails again
```

Out-of-Scope (M11): Manual text extraction by HR, OCR for non-resume documents.

---

### **[US-M11-04] Blacklist Candidate with Mandatory Reason**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to blacklist a candidate with a mandatory reason so that future applications from that candidate are automatically blocked and the reason is on record.

**Gherkin AC:**
```
Feature: Candidate Blacklist
  Scenario: HR blacklists a candidate
    Given HR is on a candidate's detail page
    When HR clicks Blacklist
    Then a dialog appears requiring a mandatory reason field
    When HR enters a reason and confirms
    Then the candidate's status is set to Blacklisted
    And the reason, timestamp, and HR user are recorded in the Blacklist History Log
    And the candidate is removed from the Talent Pool if starred (per M10)

  Scenario: Blacklisted candidate's future applications are blocked
    Given candidate C is blacklisted
    When candidate C attempts to submit a new application via the applicant portal
    Then the system blocks the submission and shows an appropriate message (US-M03-03)

  Scenario: Blacklisted candidate is hidden from waiting lists
    Given candidate C is blacklisted
    When HR opens any job's Resume Review waiting list
    Then candidate C does not appear in the list even if they have a pending application

  Scenario: Blacklist reason is mandatory
    Given HR clicks Blacklist on a candidate
    When HR attempts to confirm without entering a reason
    Then the system shows a validation error "Reason is required" and does not save
```

Out-of-Scope (M11): Auto-blacklist based on score thresholds, blacklisting by non-HR roles.

---

### **[US-M11-05] Unblacklist Candidate with Mandatory Reason and History Preserved**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to unblacklist a candidate with a mandatory reason so that the candidate can apply again and the full blacklist/unblacklist history is preserved.

**Gherkin AC:**
```
Feature: Candidate Unblacklist
  Scenario: HR unblacklists a candidate
    Given candidate C is blacklisted
    When HR clicks Unblacklist on candidate C's detail page
    Then a dialog appears requiring a mandatory reason field
    When HR enters a reason and confirms
    Then candidate C's status is restored to normal (not Starred)
    And the unblacklist event is added to the Blacklist History Log with reason, timestamp, and HR name

  Scenario: Unblacklisted candidate can submit new applications
    Given candidate C has been unblacklisted
    When candidate C submits a new application
    Then the application is accepted normally

  Scenario: Full blacklist history is visible in candidate detail
    Given candidate C has been blacklisted once and unblacklisted once
    When HR views the timeline on candidate C's detail page
    Then both the blacklist event and the unblacklist event appear in chronological order
    Each showing: event type, reason, timestamp, and HR name
```

Out-of-Scope (M11): Unblacklist restoring Starred status (must be re-starred manually), bulk unblacklist.

---

## M12 — Supervisor Requisition

---

### **[US-M12-01] Supervisor Creates and Submits Requisition**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As a Supervisor I want to create a job requisition with full position details and submit it to HR so that I can formally request a new hire.

**Gherkin AC:**
```
Feature: Supervisor Creates Requisition
  Scenario: Supervisor creates a new requisition
    Given Supervisor S is on the Requisition page
    When Supervisor S fills in:
      | Field             | Notes                              |
      | Title             | Required                           |
      | Department        | Required, from master list         |
      | Headcount         | Required, number                   |
      | Job Description   | Rich text                          |
      | Hard Skills       | Multi-select from master list      |
      | Education         | Required                           |
      | Experience Level  | Required, from master list         |
      | Job Type          | Full-time / Part-time / Contract   |
    And clicks Save Draft
    Then the requisition is saved with status Draft

  Scenario: Supervisor submits requisition to HR
    Given Supervisor S has a Draft requisition
    When Supervisor S clicks Submit
    Then the requisition status changes to Pending
    And the system sends notification FRNOT01 to HR
    And Supervisor S can no longer edit the requisition while it is Pending

  Scenario: Supervisor tracks requisition status
    Given Supervisor S has submitted requisitions
    When Supervisor S views the Requisition list
    Then each requisition shows its current status: Draft / Pending / In-Progress / Published
```

Out-of-Scope (M12): Supervisor editing a Pending or In-Progress requisition, Supervisor converting requisition to job posting.

---

### **[US-M12-02] HR Reviews and Edits Requisition**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to view and edit a submitted requisition before converting it to a job posting so that I can adjust details to match HR's standards.

**Gherkin AC:**
```
Feature: HR Reviews Requisition
  Scenario: HR receives notification and views new requisition
    Given Supervisor S has submitted a requisition
    And HR has received notification FRNOT01
    When HR opens the Requisitions list
    Then the new requisition appears with status Pending
    And HR can open it to view all fields submitted by the Supervisor

  Scenario: HR edits a requisition before converting
    Given HR is viewing a Pending requisition
    When HR edits any field (e.g. title, headcount, JD)
    Then the changes are saved on the requisition record
    And the requisition status remains Pending until HR converts it

  Scenario: HR sees all requisitions regardless of Supervisor
    Given multiple Supervisors have submitted requisitions
    When HR opens the Requisitions list
    Then HR sees all submitted requisitions from all Supervisors
    And can filter by status or department
```

Out-of-Scope (M12): HR rejecting a requisition (HR converts or ignores), Supervisor visibility of HR edits.

---

### **[US-M12-03] HR Converts Requisition to Job Posting (1-Click Pre-Fill)**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to convert an approved requisition into a job posting with one click and have the job form pre-filled from the requisition data so that I do not need to re-enter information.

**Gherkin AC:**
```
Feature: Convert Requisition to Job Posting
  Scenario: HR converts a requisition to a job posting
    Given HR is viewing a Pending requisition
    When HR clicks "Convert to Job Posting"
    Then HR is taken to the Create Job Posting form
    And the following fields are pre-filled from the requisition:
      | Job Form Field    | Source Requisition Field |
      | Title             | Title                    |
      | Department        | Department               |
      | Headcount         | Headcount                |
      | Job Description   | Job Description          |
      | Hard Skills       | Hard Skills              |
      | Education         | Education                |
      | Experience Level  | Experience Level         |
      | Job Type          | Job Type                 |
    And HR can edit any pre-filled field before publishing

  Scenario: Requisition status updates on conversion
    Given HR has clicked "Convert to Job Posting" for requisition R
    When HR saves or publishes the resulting job posting
    Then requisition R's status changes to In-Progress (or Published when the job goes live)

  Scenario: HR notifies Supervisor when job is published
    Given a job posting converted from Supervisor S's requisition is published
    Then the system sends notification FRNOT02 to Supervisor S
```

Out-of-Scope (M12): One requisition converting to multiple job postings, automatic publish without HR review.

---

## M13 — Notification System

---

### **[US-M13-01] In-App Notification Bell with Unread Badge**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user or Supervisor I want a notification bell in the main navigation with an unread badge count so that I am immediately aware of pending actions without leaving the current screen.

**Gherkin AC:**
```
Feature: In-App Notification Bell
  Scenario: Unread badge appears on new notification
    Given a new notification event is triggered for HR user U
    When user U views any page in apps/main
    Then the notification bell in the main nav shows a badge with the unread count
    And the badge increments for each additional unread notification

  Scenario: HR opens the notification inbox
    Given HR has unread notifications
    When HR clicks the notification bell
    Then a notification inbox panel opens showing all notifications in reverse chronological order
    And each notification shows: event description, related entity name, and timestamp

  Scenario: Notification is marked read on click
    Given the inbox is open with unread notifications
    When HR clicks on a notification
    Then that notification is marked as read
    And the unread badge count decreases by 1

  Scenario: Badge disappears when all notifications are read
    Given HR reads all notifications in the inbox
    Then the unread badge is removed from the bell icon
```

Out-of-Scope (M13): Bulk mark-all-as-read button, notification archiving, notification snooze.

---

### **[US-M13-02] Per-User Notification Channel Toggle (In-App vs Email)**
**3 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user or Supervisor I want to independently toggle In-App and Email notifications per event type so that I receive only the notification channels I need for each event.

**Gherkin AC:**
```
Feature: Notification Channel Preferences
  Scenario: User views notification preference settings
    Given user U navigates to Notification Settings
    Then U sees a table of all 10 notification event types with two toggle columns:
      | Column  | Description                  |
      | In-App  | Toggle In-App on/off per event|
      | Email   | Toggle Email on/off per event |

  Scenario: User toggles Email off for an event
    Given notification event FRNOT03 (new applicant) has Email enabled for user U
    When user U turns off the Email toggle for FRNOT03
    Then future FRNOT03 events send In-App only (if enabled) and no email to user U
    And the In-App toggle for FRNOT03 is unaffected

  Scenario: User toggles In-App off for an event
    Given user U turns off the In-App toggle for an event
    Then future occurrences of that event produce no in-app notification for user U
    But email for that event (if enabled) is still sent

  Scenario: Default channel settings applied on first login
    Given a new user account is created
    When the user logs in for the first time
    Then notification preferences are set to system defaults:
      | Event           | Default Channels               |
      | FRNOT03         | In-App only                    |
      | FRNOT09         | In-App only                    |
      | FRNOT10         | In-App only                    |
      | All other events| In-App + Email                 |
```

Out-of-Scope (M13): Admin-level override of user notification preferences, SMS channel, push notifications.

---

### **[US-M13-03] System Triggers All 10 Defined Notification Events**
**5 pts | Must Have**
Tags: `backend`

As the system I want to automatically trigger the correct notification to the correct recipient(s) for each of the 10 defined events so that HR and Supervisors are always informed of relevant actions without manual intervention.

**Gherkin AC:**
```
Feature: Notification Event Triggers
  Scenario: FRNOT01 — New Requisition submitted
    Given Supervisor S submits a new requisition
    Then the system sends FRNOT01 to all HR users via their configured channels

  Scenario: FRNOT02 — Job published (converted from requisition)
    Given HR publishes a job converted from Supervisor S's requisition
    Then the system sends FRNOT02 to Supervisor S via their configured channels

  Scenario: FRNOT03 — New Applicant submitted
    Given a candidate submits an application for job J
    Then the system sends FRNOT03 to the HR owner of job J

  Scenario: FRNOT04 — All tests submitted for a round
    Given all applicants in a test round have submitted their tests
    Then the system sends FRNOT04 to the responsible HR user

  Scenario: FRNOT05 — Technical Test graded by Supervisor
    Given Supervisor S completes grading of a Technical Test
    Then the system sends FRNOT05 to the HR user managing that job

  Scenario: FRNOT06 — All interview scores submitted
    Given all interviewers have submitted scores for an interview round
    Then the system sends FRNOT06 to the responsible HR user

  Scenario: FRNOT07 — Assigned as interviewer
    Given HR assigns Supervisor S as an interviewer for a round
    Then the system sends FRNOT07 to Supervisor S

  Scenario: FRNOT08 — Assigned to grade Technical Test
    Given HR or the system assigns Supervisor S to grade a Technical Test
    Then the system sends FRNOT08 to Supervisor S

  Scenario: FRNOT09 — Offer response received
    Given a candidate responds to an offer
    Then the system sends FRNOT09 to the responsible HR user

  Scenario: FRNOT10 — Job auto-closed by scheduler
    Given a job's close date is reached and the scheduler auto-closes it
    Then the system sends FRNOT10 to the HR owner of that job

  Scenario: Notification respects user channel preference
    Given user U has disabled Email for FRNOT05
    When a FRNOT05 event is triggered for user U
    Then the system sends only an In-App notification to user U
    And no email is sent
```

Out-of-Scope (M13): Notification retry on email delivery failure beyond system responsibility, notification templates customisation.

---

## M14 — Setup / Pre-configuration

---

### **[US-M14-01] HR Manages System and Custom Email Templates**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to view hardcoded system email templates and create/edit/delete custom HR-managed templates so that all outgoing emails use the correct, approved content.

**Gherkin AC:**
```
Feature: Email Template Management
  Scenario: HR views system (SYS) templates
    Given HR navigates to Setup > Email Templates
    Then HR sees the 3 system templates: SYS01, SYS02, SYS03
    And the templates are displayed as read-only (HR cannot edit or delete them)

  Scenario: HR creates a custom (TPL) template
    Given HR is on the Email Templates page
    When HR clicks Create Template and fills in name, subject, and body with merge tags
    Then the new template is saved and available for use in pipeline actions

  Scenario: HR edits a custom template
    Given a custom template TPL04 exists
    When HR edits the subject or body
    Then changes are saved and take effect for future emails only (not already-sent emails)

  Scenario: HR deletes a custom template with reference check
    Given TPL05 is referenced as the Interview Invitation template for active jobs
    When HR attempts to delete TPL05
    Then the system blocks deletion and shows "Template TPL05 is in use — cannot be deleted"
    When a template has no active references
    Then deletion proceeds after confirmation

  Scenario: System blocks pipeline action if required TPL is missing
    Given HR attempts to send an Interview Invitation
    And TPL05 does not exist in the system
    Then the system blocks the action and shows "Email template TPL05 is required — configure it in Setup before proceeding"
```

Out-of-Scope (M14): HR editing SYS templates, email template versioning, preview with real applicant data.

---

### **[US-M14-02] HR Creates and Manages Job Description Templates**
**2 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to create and manage standard Job Description templates so that a default JD is available when creating new job postings.

**Gherkin AC:**
```
Feature: Job Description Template Management
  Scenario: HR creates a JD template
    Given HR navigates to Setup > Job Description Templates
    When HR creates a new template with a title and rich-text body
    Then the template is saved and available as a default when creating a new job posting

  Scenario: JD template pre-fills the Job Posting form
    Given at least one JD template exists
    When HR creates a new job posting and selects a JD template
    Then the Job Description field in the posting form is pre-filled with the template content
    And HR can edit the pre-filled content before saving

  Scenario: HR edits or deletes a JD template
    Given a JD template exists
    When HR edits the title or body
    Then the change is saved and does not affect already-created job postings
    When HR deletes the template
    Then it is removed from the selection list for new job postings
```

Out-of-Scope (M14): Versioning of JD templates, applicant-visible JD templates.

---

### **[US-M14-03] HR Manages Master Data (Department, Employment Type, Experience Level, Hard Skill)**
**5 pts | Must Have**
Tags: `frontend`, `backend`

As an HR user I want to create, edit, and delete master data items (Department, Employment Type, Experience Level, Hard Skill) with reference checks so that dropdown options across the system are always accurate.

**Gherkin AC:**
```
Feature: Master Data Management
  Scenario: HR creates a new master data item
    Given HR is in Setup > Manage Data
    When HR selects a category (e.g. Department) and clicks Add
    Then HR enters the item name and saves
    And the new item appears in the relevant dropdowns system-wide immediately

  Scenario: HR edits an existing item
    Given a Department named "Engineering" exists
    When HR renames it to "Software Engineering"
    Then all existing records that reference "Engineering" reflect the updated name

  Scenario: HR deletes an item with reference check
    Given Department "HR" is referenced by 3 active job postings
    When HR attempts to delete "HR"
    Then the system blocks deletion and shows "Item is referenced by active records — cannot be deleted"
    When an item has no active references
    Then deletion proceeds after confirmation

  Scenario: HR filters and refreshes the master data list
    Given HR is on the Manage Data page
    When HR selects a category filter (e.g. Hard Skill)
    Then only Hard Skill items are shown
    When HR clicks Refresh
    Then the list reloads with the latest data
```

Out-of-Scope (M14): Bulk import of master data, reordering of dropdown items, applicant-facing display of master data labels.

---

### **[US-M14-04] JD and Hard Skill Embedding Sent to AI Service After Job Create/Edit**
**3 pts | Must Have**
Tags: `backend`, `ai-service`

As the system I want to send the job's JD text and Hard Skill list to the AI service for embedding and vector storage whenever a job is created or edited so that similarity computation uses up-to-date vectors.

**Gherkin AC:**
```
Feature: JD and Hard Skill Embedding
  Scenario: Embedding triggered on job create
    Given HR creates a new job posting with a JD and Hard Skills
    When the job is saved
    Then the system sends the JD text and Hard Skill list to the AI service asynchronously
    And the AI service stores the resulting vectors and weight values for that job

  Scenario: Embedding triggered on job edit
    Given HR edits the JD or Hard Skills of an existing job
    When the edit is saved
    Then the system re-sends the updated JD and Hard Skill list to the AI service
    And the stored vectors for that job are updated

  Scenario: Embedding does not block job save
    Given HR saves a job posting
    Then the save confirmation is immediate
    And the embedding request runs as a background task without delaying the user

  Scenario: Weight values stored alongside vectors
    Given HR has configured JD weight and Hard Skill weight via the Weight Popup (US-M06-02)
    When the embedding job completes
    Then the weight values are stored with the job's AI profile for use in Overall Score calculation
```

Out-of-Scope (M14): Real-time re-embedding on weight slider change, embedding for non-job content.

---

## M15 — Localization and Theme

---

### **[US-M15-01] Thai / English Language Toggle**
**3 pts | Must Have**
Tags: `frontend`

As any user (HR, Supervisor, Admin, Applicant) I want to toggle between Thai and English language so that I can use the system in my preferred language.

**Gherkin AC:**
```
Feature: Language Toggle
  Scenario: User switches language from Thai to English
    Given the app is displayed in Thai (default)
    When the user clicks the language toggle in the Header
    Then all UI text switches to English immediately without a full page reload
    And the selected language is persisted in the user's profile (logged-in) or localStorage (applicant portal)

  Scenario: User switches language from English to Thai
    Given the app is displayed in English
    When the user clicks the language toggle
    Then all UI text switches to Thai immediately
    And the selection is persisted

  Scenario: Default language is Thai on first load
    Given a new user opens the app for the first time
    Then the UI is displayed in Thai
    And no language preference is set yet in their profile

  Scenario: User-generated content is NOT translated
    Given the app is in English mode
    When HR views a Job Description or Scoring Criteria written in Thai
    Then those fields display as-is in Thai (no auto-translation applied)
```

Out-of-Scope (M15): Machine translation of user-generated content, additional languages beyond Thai and English.

---

### **[US-M15-02] All UI Text via Translation Keys (next-intl)**
**3 pts | Must Have**
Tags: `frontend`

As a developer I want all UI text to be driven by translation keys in `messages/th.json` and `messages/en.json` with nested namespaces so that no hardcoded strings appear in components.

**Gherkin AC:**
```
Feature: Translation Key Coverage
  Scenario: All static UI strings use translation keys
    Given the codebase uses next-intl
    When any UI component renders a user-visible string
    Then the string is fetched via a translation key (e.g. t('common.save'))
    And no hardcoded English or Thai string appears directly in JSX

  Scenario: Date and number formats follow locale
    Given the app locale is set to Thai
    When a date is displayed (e.g. application submission date)
    Then the date is formatted per Thai locale conventions (e.g. DD/MM/YYYY, Buddhist Era if applicable)
    When the app locale is set to English
    Then the date is formatted per English locale conventions (e.g. MMM D, YYYY)

  Scenario: Missing translation key falls back gracefully
    Given a translation key exists in th.json but is missing from en.json
    When the app is in English mode and renders that key
    Then the Thai value is shown as fallback rather than a raw key string
```

Out-of-Scope (M15): Runtime translation key editing, translation management UI inside the app.

---

### **[US-M15-03] Dark / Light Theme Toggle with CSS Variables**
**3 pts | Must Have**
Tags: `frontend`

As any user I want to toggle between Dark and Light mode so that I can use the interface comfortably in any lighting environment.

**Gherkin AC:**
```
Feature: Dark / Light Theme Toggle
  Scenario: User switches to Dark mode
    Given the app is in Light mode (default or OS preference)
    When the user clicks the theme toggle in the Header
    Then all UI components switch to the dark color scheme instantly
    And the preference is persisted in the user's profile (logged-in) or localStorage (applicant portal)

  Scenario: User switches to Light mode
    Given the app is in Dark mode
    When the user clicks the theme toggle
    Then all UI components switch to the light color scheme instantly

  Scenario: Default theme follows OS preference on first load
    Given a new user opens the app with no saved preference
    Then the theme matches the user's OS dark/light mode preference

  Scenario: All components respect the active theme
    Given Dark mode is active
    Then all of the following components use dark-mode CSS variable values:
      | Component  |
      | Sidebar    |
      | Navbar     |
      | Table      |
      | Card       |
      | Modal      |
      | Form Input |
      | Button     |
      | Badge      |
      | Chart      |
      | Dropdown   |
      | Tooltip    |
    And no component retains hardcoded light-mode colors in dark mode
```

Out-of-Scope (M15): Per-component theme override, custom user-defined color themes.

---

### **[US-M15-04] WCAG 2.1 AA Contrast Compliance for Both Themes**
**2 pts | Must Have**
Tags: `frontend`

As any user I want all UI text and interactive elements to meet WCAG 2.1 AA contrast ratio (4.5:1 minimum) in both Light and Dark mode so that the system is accessible.

**Gherkin AC:**
```
Feature: Accessibility Contrast Compliance
  Scenario: Text contrast meets 4.5:1 in Light mode
    Given the app is in Light mode
    When any text element is rendered
    Then the contrast ratio between the text color and its background is at least 4.5:1
    As verified by an automated accessibility audit tool

  Scenario: Text contrast meets 4.5:1 in Dark mode
    Given the app is in Dark mode
    When any text element is rendered
    Then the contrast ratio between the text color and its background is at least 4.5:1

  Scenario: Interactive element contrast (buttons, links) meets 4.5:1
    Given either Light or Dark mode is active
    When a Button, Link, or form control is rendered
    Then the contrast ratio for its label/icon against its background meets WCAG 2.1 AA (4.5:1)

  Scenario: TigerSoft brand colors are used within contrast constraints
    Given Vivid Red #F4001A is used for CTA buttons and Oxford Blue #0B1F3A is used for text
    Then these colors are applied in contexts where the 4.5:1 ratio is maintained
    And alternative accessible pairings are used where brand colors alone would fail the ratio
```

Out-of-Scope (M15): WCAG AAA compliance, screen reader full audit, keyboard navigation audit (separate QA scope).

---

### Summary Table

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

### Total
Stories: 43 | Points: 143

===PHASE REPORT END===
Status: COMPLETE
Phase: P1 Batch 2 (M06–M15)
