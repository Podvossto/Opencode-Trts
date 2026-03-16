"""
Seed Batch 2 user stories (M06–M15, 43 stories) into Notion Backlog DB.
Uses urllib only — no requests, no curl.
Run: python .opencode/scripts/seed_batch2.py
"""

import json
import urllib.request
import urllib.error
import time

NOTION_TOKEN = "ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA"
DATABASE_ID = "3250531c-537a-8109-aa3d-cd3961fe0257"

HEADERS = {
    "Authorization": f"Bearer {NOTION_TOKEN}",
    "Content-Type": "application/json",
    "Notion-Version": "2022-06-28",
}

STORIES = [
    # M06
    {"id": "US-M06-01", "title": "3-Panel Resume Review Layout",                                   "module": "M06", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M06-02", "title": "Overall Score Calculation with Weight Popup",                    "module": "M06", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend", "ai-service"]},
    {"id": "US-M06-03", "title": "HR Marks Up to 3 Suitable Positions",                           "module": "M06", "points": 2,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M06-04", "title": "Resume Review Action Buttons",                                   "module": "M06", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M06-05", "title": "Auto-Compute Similarity for Starred Candidates on New Job Publish", "module": "M06", "points": 3, "priority": "Must Have", "tags": ["backend", "ai-service"]},
    # M07
    {"id": "US-M07-01", "title": "HR Initiates Action Test with Round Badge",                     "module": "M07", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M07-02", "title": "HR Selects Applicants and Assigns Test",                        "module": "M07", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M07-03", "title": "Online Test Delivery with Anti-Cheat",                          "module": "M07", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M07-04", "title": "MCQ Auto-Grading and Manual Grading Inbox",                     "module": "M07", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M07-05", "title": "HR Opens Phase 2 Selection and Leaderboard",                    "module": "M07", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M07-06", "title": "HR and Supervisor CRUD Test Library",                           "module": "M07", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    # M08
    {"id": "US-M08-01", "title": "HR Initiates Action Interview with Round Badge",                "module": "M08", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-02", "title": "HR Defines Scoring Criteria per Round",                         "module": "M08", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-03", "title": "HR Assigns Interviewers — Supervisor and Temp Interviewers",    "module": "M08", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-04", "title": "HR Views and Manages Invite Links List",                        "module": "M08", "points": 2,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-05", "title": "Phase 1 Interview Grading — Real-Time Per-Interviewer Scores",  "module": "M08", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-06", "title": "Phase 2 Interview Selection — Leaderboard and Final Decisions", "module": "M08", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M08-07", "title": "HR Sends Interview Invitation Email (TPL05)",                   "module": "M08", "points": 2,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    # M09
    {"id": "US-M09-01", "title": "HR Initiates Transfer During or After Interview",               "module": "M09", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M09-02", "title": "Transfer History Carries Over and Unlimited Transfers Supported","module": "M09", "points": 3,  "priority": "Must Have",  "tags": ["backend"]},
    {"id": "US-M09-03", "title": "System Warns HR if Target Position is Closed",                  "module": "M09", "points": 2,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M09-04", "title": "Block Sending Test from Wrong Position After Transfer",          "module": "M09", "points": 2,  "priority": "Must Have",  "tags": ["backend"]},
    # M10
    {"id": "US-M10-01", "title": "Auto-Compute Similarity for Starred Candidates Against All New Jobs", "module": "M10", "points": 3, "priority": "Must Have", "tags": ["backend", "ai-service"]},
    {"id": "US-M10-02", "title": "Starred Candidates Removed from Pool on Blacklist",             "module": "M10", "points": 2,  "priority": "Must Have",  "tags": ["backend"]},
    # M11
    {"id": "US-M11-01", "title": "Candidate List with Search and Filter",                         "module": "M11", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M11-02", "title": "Candidate Detail — Full Application History and Documents",     "module": "M11", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M11-03", "title": "OCR Auto-Trigger on Resume Upload",                             "module": "M11", "points": 3,  "priority": "Must Have",  "tags": ["backend", "ai-service"]},
    {"id": "US-M11-04", "title": "Blacklist Candidate with Mandatory Reason",                     "module": "M11", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M11-05", "title": "Unblacklist Candidate with Mandatory Reason and History Preserved", "module": "M11", "points": 2, "priority": "Must Have", "tags": ["frontend", "backend"]},
    # M12
    {"id": "US-M12-01", "title": "Supervisor Creates and Submits Requisition",                    "module": "M12", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M12-02", "title": "HR Reviews and Edits Requisition",                              "module": "M12", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M12-03", "title": "HR Converts Requisition to Job Posting (1-Click Pre-Fill)",     "module": "M12", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    # M13
    {"id": "US-M13-01", "title": "In-App Notification Bell with Unread Badge",                    "module": "M13", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M13-02", "title": "Per-User Notification Channel Toggle (In-App vs Email)",        "module": "M13", "points": 3,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M13-03", "title": "System Triggers All 10 Defined Notification Events",            "module": "M13", "points": 5,  "priority": "Must Have",  "tags": ["backend"]},
    # M14
    {"id": "US-M14-01", "title": "HR Manages System and Custom Email Templates",                  "module": "M14", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M14-02", "title": "HR Creates and Manages Job Description Templates",              "module": "M14", "points": 2,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M14-03", "title": "HR Manages Master Data",                                        "module": "M14", "points": 5,  "priority": "Must Have",  "tags": ["frontend", "backend"]},
    {"id": "US-M14-04", "title": "JD and Hard Skill Embedding Sent to AI Service After Job Create/Edit", "module": "M14", "points": 3, "priority": "Must Have", "tags": ["backend", "ai-service"]},
    # M15
    {"id": "US-M15-01", "title": "Thai / English Language Toggle",                                "module": "M15", "points": 3,  "priority": "Must Have",  "tags": ["frontend"]},
    {"id": "US-M15-02", "title": "All UI Text via Translation Keys (next-intl)",                  "module": "M15", "points": 3,  "priority": "Must Have",  "tags": ["frontend"]},
    {"id": "US-M15-03", "title": "Dark / Light Theme Toggle with CSS Variables",                  "module": "M15", "points": 3,  "priority": "Must Have",  "tags": ["frontend"]},
    {"id": "US-M15-04", "title": "WCAG 2.1 AA Contrast Compliance for Both Themes",               "module": "M15", "points": 2,  "priority": "Must Have",  "tags": ["frontend"]},
]


def create_page(story):
    url = "https://api.notion.com/v1/pages"
    payload = {
        "parent": {"database_id": DATABASE_ID},
        "properties": {
            "Task ID": {
                "title": [{"text": {"content": f"[{story['id']}] {story['title']}"}}]
            },
            "Status": {
                "select": {"name": "BACKLOG"}
            },
            "Tag": {
                "multi_select": [{"name": t} for t in story["tags"]]
            },
            "Priority": {
                "select": {"name": "normal"}
            },
            "Points": {
                "number": story["points"]
            },
            "Story": {
                "rich_text": [{"text": {"content": story["id"]}}]
            },
            "Gherkin AC": {
                "rich_text": [{"text": {"content": f"Module: {story['module']} | MoSCoW: {story['priority']}"}}]
            },
        }
    }
    data = json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(url, data=data, headers=HEADERS, method="POST")
    try:
        with urllib.request.urlopen(req) as resp:
            result = json.loads(resp.read().decode("utf-8"))
            return result.get("id", "unknown")
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8")
        print(f"  ERROR {e.code}: {body}")
        return None


def main():
    print(f"Seeding {len(STORIES)} Batch 2 stories into Notion DB {DATABASE_ID}...\n")
    success = 0
    failed = 0
    for i, story in enumerate(STORIES, 1):
        print(f"  [{i:02d}/{len(STORIES)}] {story['id']} — {story['title'][:55]}...", end=" ")
        page_id = create_page(story)
        if page_id:
            print(f"OK ({page_id[:8]}...)")
            success += 1
        else:
            print("FAILED")
            failed += 1
        time.sleep(0.35)  # stay well under Notion rate limit (3 req/s)

    print(f"\nDone. Success: {success} | Failed: {failed}")
    if failed == 0:
        print("All 43 Batch 2 stories inserted successfully.")


if __name__ == "__main__":
    main()
