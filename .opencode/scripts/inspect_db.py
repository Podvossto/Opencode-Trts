"""Introspect the Notion DB to get real property names and types."""
import json
import urllib.request

NOTION_TOKEN = "ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA"
DATABASE_ID = "3250531c-537a-8109-aa3d-cd3961fe0257"

req = urllib.request.Request(
    f"https://api.notion.com/v1/databases/{DATABASE_ID}",
    headers={
        "Authorization": f"Bearer {NOTION_TOKEN}",
        "Notion-Version": "2022-06-28",
    }
)
with urllib.request.urlopen(req) as resp:
    db = json.loads(resp.read().decode("utf-8"))

print("DB title:", db.get("title", [{}])[0].get("plain_text", ""))
print("\nProperties:")
for name, prop in db["properties"].items():
    print(f"  '{name}': {prop['type']}")
