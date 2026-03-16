const { Database } = require('bun:sqlite');
const db = new Database('C:/Users/Pitchaya.h/.local/share/opencode/opencode.db');
// Check schema first
const tables = db.query("SELECT name FROM sqlite_master WHERE type='table'").all();
console.log('Tables:', tables.map(t => t.name).join(', '));
const cols = db.query("PRAGMA table_info(session)").all();
console.log('Session columns:', cols.map(c => c.name).join(', '));
