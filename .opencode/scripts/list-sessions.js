const { Database } = require('bun:sqlite');
const db = new Database('C:/Users/Pitchaya.h/.local/share/opencode/opencode.db');

// Show all sessions to understand what's there
const all = db.query("SELECT id, slug, title, time_created, time_updated FROM session ORDER BY time_updated DESC LIMIT 30").all();
console.log('All recent sessions:');
all.forEach(s => console.log(JSON.stringify(s)));
