"""
Set all pages in the User Stories & Backlog DB to Status: IN DESIGN.
Uses urllib only — no requests, no curl.
Run: python .opencode/scripts/set_in_design.py
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


def query_all_pages():
    """Fetch all page IDs from the database."""
    pages = []
    has_more = True
    start_cursor = None

    while has_more:
        payload = {"page_size": 100}
        if start_cursor:
            payload["start_cursor"] = start_cursor

        data = json.dumps(payload).encode("utf-8")
        req = urllib.request.Request(
            f"https://api.notion.com/v1/databases/{DATABASE_ID}/query",
            data=data,
            headers=HEADERS,
            method="POST",
        )
        try:
            with urllib.request.urlopen(req) as resp:
                result = json.loads(resp.read())
                pages.extend(result.get("results", []))
                has_more = result.get("has_more", False)
                start_cursor = result.get("next_cursor")
        except urllib.error.HTTPError as e:
            print(f"Query error: {e.code} {e.read().decode()}")
            break

    return pages


def update_status(page_id, task_id_title):
    payload = {
        "properties": {
            "Status": {
                "select": {"name": "IN DESIGN"}
            }
        }
    }
    data = json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(
        f"https://api.notion.com/v1/pages/{page_id}",
        data=data,
        headers=HEADERS,
        method="PATCH",
    )
    try:
        with urllib.request.urlopen(req) as resp:
            resp.read()
            return True
    except urllib.error.HTTPError as e:
        body = e.read().decode()
        print(f"  ERROR [{task_id_title}]: {e.code} — {body[:200]}")
        return False


def main():
    print("Fetching all pages from Backlog DB...")
    pages = query_all_pages()
    print(f"Found {len(pages)} pages.")

    ok = 0
    fail = 0
    for page in pages:
        page_id = page["id"]
        # Extract title from Task ID property
        title_parts = page.get("properties", {}).get("Task ID", {}).get("title", [])
        task_label = title_parts[0]["plain_text"] if title_parts else page_id
        success = update_status(page_id, task_label)
        if success:
            ok += 1
            print(f"  OK {task_label}")
        else:
            fail += 1
        time.sleep(0.15)  # Respect Notion rate limit

    print(f"\nDone: {ok} OK, {fail} failed out of {len(pages)} pages.")


if __name__ == "__main__":
    main()
