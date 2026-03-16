// P3 Notion Setup: Create Sprint 1 + all 21 tasks (corrected property names)
const NOTION_KEY = 'ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA';
const SPRINT_DB = '3250531c-537a-8131-b7a3-efe66af39c4f';
const TASKS_DB = '3250531c-537a-8109-aa3d-cd3961fe0257';
const NOTION_VERSION = '2022-06-28';

async function notionPost(url, body) {
  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${NOTION_KEY}`,
      'Content-Type': 'application/json',
      'Notion-Version': NOTION_VERSION,
    },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.text();
    throw new Error(`Notion ${res.status}: ${err}`);
  }
  return res.json();
}

async function createSprint() {
  console.log('Creating Sprint 1...');
  const result = await notionPost('https://api.notion.com/v1/pages', {
    parent: { database_id: SPRINT_DB },
    properties: {
      'Sprint Name': { title: [{ text: { content: 'Sprint 1 — Auth Foundation + Portal Intake' } }] },
      Status: { select: { name: 'Active' } },
      Feature: { rich_text: [{ text: { content: 'M01 Auth + M02 HR Auth + M03 Portal (partial)' } }] },
      Velocity: { number: 46 },
    },
  });
  console.log(`  Sprint 1 created: ${result.id}`);
  return result.id;
}

const TASKS = [
  { id: 'ATS-001', title: 'Docker Compose + DB migration setup', tag: 'infra', priority: 'normal', sp: 2, assignee: 'dev', desc: 'Configure docker-compose.yml services, run initial + additive migrations, verify all 35 tables created' },
  { id: 'ATS-002', title: 'Redis setup + JWT blacklist config', tag: 'infra', priority: 'normal', sp: 2, assignee: 'dev', desc: 'Configure Redis container, implement JWT blacklist store in pkg/redis, wire into auth middleware' },
  { id: 'ATS-003', title: 'Login endpoint — admin/hr/supervisor', tag: 'backend', priority: 'Must Have', sp: 3, assignee: 'dev', desc: 'POST /api/auth/login — validate credentials, issue JWT with role claim, return token' },
  { id: 'ATS-004', title: 'Logout + JWT blacklist', tag: 'backend', priority: 'Must Have', sp: 2, assignee: 'dev', desc: 'POST /api/auth/logout — add token to Redis blacklist, middleware rejects blacklisted tokens' },
  { id: 'ATS-005', title: 'Force password change on first login', tag: 'backend', priority: 'Must Have', sp: 2, assignee: 'dev', desc: 'POST /api/auth/change-password — check must_change_password flag, enforce change before any other action' },
  { id: 'ATS-006', title: 'Session timeout 30min', tag: 'backend', priority: 'normal', sp: 1, assignee: 'dev', desc: 'JWT expiry set to 30min, middleware validates exp claim, auto-reject expired tokens' },
  { id: 'ATS-007', title: 'OTP request + verify + password reset', tag: 'backend', priority: 'Must Have', sp: 3, assignee: 'dev', desc: 'POST /api/auth/otp/request — generate OTP, send via email. POST /api/auth/otp/verify — validate OTP, reset password' },
  { id: 'ATS-008', title: 'List users with pagination + filter', tag: 'backend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'GET /api/users — paginated list, filter by role/department/active status' },
  { id: 'ATS-009', title: 'Create user (admin only)', tag: 'backend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'POST /api/users — admin creates user with must_change_password=true, validates unique email' },
  { id: 'ATS-010', title: 'Edit + Delete (deactivate) user', tag: 'backend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'PUT /api/users/:id + DELETE /api/users/:id — edit profile, soft-delete via is_active=false' },
  { id: 'ATS-011', title: 'Admin resets user password', tag: 'backend', priority: 'normal', sp: 1, assignee: 'dev', desc: 'POST /api/users/:id/reset-password — admin resets, sets must_change_password=true' },
  { id: 'ATS-012', title: 'Public job listing + detail endpoints', tag: 'backend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'GET /api/portal/jobs (list open jobs) + GET /api/portal/jobs/:id (detail with form). No auth required' },
  { id: 'ATS-013', title: 'Login page + JWT storage', tag: 'frontend', priority: 'Must Have', sp: 3, assignee: 'dev', desc: 'Login form (email/password), call /api/auth/login, store JWT, redirect to dashboard. TigerSoft branding' },
  { id: 'ATS-014', title: 'Force password change modal', tag: 'frontend', priority: 'Must Have', sp: 2, assignee: 'dev', desc: 'Modal on first login when must_change_password=true, blocks navigation until changed' },
  { id: 'ATS-015', title: 'OTP password reset flow', tag: 'frontend', priority: 'Must Have', sp: 3, assignee: 'dev', desc: 'Request OTP form > enter OTP > new password form. Email input, OTP verification, password update' },
  { id: 'ATS-016', title: 'Session timeout auto-logout', tag: 'frontend', priority: 'normal', sp: 1, assignee: 'dev', desc: 'Detect JWT expiry, show timeout warning, auto-redirect to login' },
  { id: 'ATS-017', title: 'User list table with search/filter', tag: 'frontend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'Data table with pagination, search by name/email, filter by role/department. TigerSoft branding' },
  { id: 'ATS-018', title: 'Create/Edit/Delete user forms', tag: 'frontend', priority: 'normal', sp: 3, assignee: 'dev', desc: 'User form modal (create/edit), delete confirmation dialog, form validation' },
  { id: 'ATS-019', title: 'Public job listing page (Portal)', tag: 'frontend', priority: 'normal', sp: 2, assignee: 'dev', desc: 'apps/portal — list open jobs with search, job detail page with apply button. TigerSoft branding' },
  { id: 'ATS-020', title: 'Go tests — auth + users endpoints', tag: 'test', priority: 'Must Have', sp: 3, assignee: 'qa', desc: 'Unit + integration tests for all auth and user management endpoints. Cover login, logout, OTP, CRUD' },
  { id: 'ATS-021', title: 'Playwright E2E — login, password change, user CRUD, portal', tag: 'test', priority: 'Must Have', sp: 3, assignee: 'qa', desc: 'E2E tests: login flow, forced password change, user list/create/edit/delete, portal job listing' },
];

async function createTask(task) {
  const properties = {
    'Task ID': { title: [{ text: { content: task.id } }] },
    Status: { select: { name: 'READY FOR DEVELOPMENT' } },
    Tag: { multi_select: [{ name: task.tag }] },
    Priority: { select: { name: task.priority } },
    Points: { number: task.sp },
    Assignee: { select: { name: task.assignee } },
    Story: { rich_text: [{ text: { content: `${task.title}\n${task.desc}` } }] },
    Sprint: { rich_text: [{ text: { content: 'Sprint 1' } }] },
  };

  const result = await notionPost('https://api.notion.com/v1/pages', {
    parent: { database_id: TASKS_DB },
    properties,
  });
  return result.id;
}

async function main() {
  try {
    const sprintId = await createSprint();
    console.log(`Sprint 1 Page ID: ${sprintId}\n`);

    let created = 0;
    for (const task of TASKS) {
      try {
        const pageId = await createTask(task);
        created++;
        console.log(`  [${created}/21] ${task.id} -> ${pageId}`);
        await new Promise(r => setTimeout(r, 400));
      } catch (err) {
        console.error(`  FAILED ${task.id}: ${err.message}`);
      }
    }

    console.log(`\n=== DONE: Sprint 1 + ${created}/21 tasks created ===`);
    console.log(`SPRINT_1_PAGE_ID=${sprintId}`);
  } catch (err) {
    console.error('Fatal:', err.message);
    process.exit(1);
  }
}

main();
