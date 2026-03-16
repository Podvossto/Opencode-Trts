---
name: notion-workspace-architect
description: |
  Use this skill whenever the user asks to design or structure a Notion Workspace for a software development team — including mentions of "Notion workspace", "SDLC", "sprint planning", "project structure in Notion", or requests for a Notion template for software project management.
  This skill produces a Workspace Structure, Page Hierarchy, Workflow overview, and Layout guidelines suited for software teams following SDLC principles.
---

# Notion Workspace Architect

## Purpose

Design a Notion Workspace that helps software teams plan, design, develop, and track projects systematically. Prioritize clarity, simplicity, and practical usability over completeness.

---

## Step 1 — Gather Context

If the user has not provided enough information, ask:
- Team size (solo / small 2–5 / medium 6–15)
- Project type (Web App, Mobile, API, Internal Tool, etc.)
- Methodology (Scrum / Kanban / Ad-hoc)
- Notion plan (Free / Plus / Business)

If unspecified, assume: medium team, Scrum, Notion Plus.

---

## Step 2 — Produce Output

Always output all 4 sections below.

---

## Section 1: Workspace Structure

Present as a tree hierarchy. Adjust depth and pages to fit the user's project size.

```
[Project Name] Workspace
├── Home Dashboard
├── Project Hub
│   ├── Project Overview
│   ├── Scope & Requirements
│   └── Roadmap
├── Sprint & Task Management
│   ├── Backlog
│   ├── Active Sprint Board
│   └── Done / Archive
├── Design
│   ├── UI/UX References
│   └── Wireframes & Mockups
├── Architecture & Technical
│   ├── Architecture Overview
│   ├── Folder Structure
│   ├── API List
│   └── Database Overview
├── Documentation
│   ├── Meeting Notes
│   ├── Decision Log
│   └── Onboarding Guide
└── Team & Settings
    ├── Team Members
    └── Working Agreements
```

---

## Section 2: Page Structure

Describe each main page with its purpose, content, and recommended Notion blocks.

**Home Dashboard**
- Purpose: Entry point — gives the team an instant project overview
- Content: Project status, sprint progress summary (linked DB), quick links, announcements
- Blocks: Callout, Columns, Linked Database View

**Project Overview**
- Purpose: Central reference for the entire team
- Content: Project name, objective, target date, tech stack (table), team roles, key milestones
- Blocks: Toggle, Table, Callout

**Scope & Requirements**
- Purpose: Define boundaries and requirements clearly
- Content: In Scope / Out of Scope, user stories or feature list (database), acceptance criteria
- Blocks: Toggle List, Database, Checkbox

**Active Sprint Board**
- Purpose: Real-time task tracking during a sprint
- Content: Kanban columns (To Do / In Progress / In Review / Done), assignee, priority, story points, sprint number
- Blocks: Database (Board View), Filter, Group By

**Backlog**
- Purpose: Holding area for all unplanned tasks and stories
- Content: Filtered view of the same Sprint Board database, priority ranking, estimation column
- Blocks: Database (Table View)

**Architecture Overview**
- Purpose: High-level system design reference
- Content: System diagram (embed or image), component list and relationships, tech stack rationale, deployment environments
- Blocks: Image, Embed, Toggle, Table

**Folder Structure**
- Purpose: Document the codebase layout
- Content: Folder tree (code block), brief description per folder
- Blocks: Code Block, Callout

**API List**
- Purpose: Central registry of all API endpoints
- Content: Database with Method, Endpoint, Description, Auth Required, Status (Draft / Active / Deprecated), Module
- Blocks: Database (Table View)

**Database Overview**
- Purpose: High-level data model reference
- Content: Table or collection names with short descriptions, relationship diagram (embed). No need to document every field.
- Blocks: Table, Embed, Toggle

---

## Section 3: Workflow

```
Project Kickoff
    → Project Overview → Scope → Architecture Design
    → Backlog Creation → Sprint Planning
    → Active Sprint  (Daily standup references Sprint Board)
    → Sprint Review → Mark Done → Retrospective
    → Repeat
```

- **Daily**: Each team member checks the Sprint Board for their tasks
- **Sprint Planning**: Pull items from Backlog into the active sprint
- **Sprint Review**: Update Done column and Roadmap
- **Documentation**: Log all meetings in Meeting Notes; log key decisions in Decision Log

---

## Section 4: Layout Guidelines

**Page Layout**
- Use 2-column layout on the Home Dashboard (quick links | status summary)
- Use Callout blocks for important notices, warnings, and status
- Use Toggle to hide details not needed at a glance
- Use Dividers to separate logical sections within a page

**Database Views**
| Page | Primary View | Secondary View |
|------|-------------|----------------|
| Sprint Board | Board (Kanban) | Table |
| Backlog | Table | — |
| API List | Table (filtered by Module) | — |
| Feature List | Table (grouped by Priority) | — |

**What to Avoid**
- Pages nested deeper than 3 levels
- Databases with more than 10 properties unless clearly needed
- Duplicating data across pages — use Linked Databases instead
- Pages with no clear owner

---

## Reference Files

- `references/notion-blocks.md` — Block types and database property recommendations
- `references/sdlc-phases.md` — SDLC phases mapped to Notion pages
- `references/templates.md` — Starter template text for key pages