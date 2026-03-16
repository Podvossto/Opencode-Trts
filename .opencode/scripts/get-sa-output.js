const { Database } = require('bun:sqlite');
const db = new Database('C:/Users/Pitchaya.h/.local/share/opencode/opencode.db');

// Find SA sessions
const sessions = db.query("SELECT id, slug, title, time_created FROM session WHERE slug LIKE '%sa%' OR slug LIKE '%ats%' ORDER BY time_created DESC LIMIT 20").all();
console.log('SA/ATS sessions:');
sessions.forEach(s => console.log(s.id, s.slug, s.title, s.time_created));

// Get message cols
const msgCols = db.query("PRAGMA table_info(message)").all();
console.log('\nMessage columns:', msgCols.map(c => c.name).join(', '));

// Get part cols
const partCols = db.query("PRAGMA table_info(part)").all();
console.log('\nPart columns:', partCols.map(c => c.name).join(', '));
