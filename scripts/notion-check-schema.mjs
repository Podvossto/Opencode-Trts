// Quick script to check DB property names
const NOTION_KEY = 'ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA';

async function getDB(dbId, label) {
  const res = await fetch(`https://api.notion.com/v1/databases/${dbId}`, {
    headers: {
      'Authorization': `Bearer ${NOTION_KEY}`,
      'Notion-Version': '2022-06-28',
    },
  });
  const data = await res.json();
  console.log(`\n=== ${label} ===`);
  console.log('Title:', data.title?.[0]?.plain_text);
  for (const [key, val] of Object.entries(data.properties || {})) {
    console.log(`  "${key}" -> type: ${val.type}${val.type === 'select' ? ', options: [' + (val.select?.options?.map(o => o.name).join(', ') || 'none') + ']' : ''}${val.type === 'multi_select' ? ', options: [' + (val.multi_select?.options?.map(o => o.name).join(', ') || 'none') + ']' : ''}`);
  }
}

async function main() {
  await getDB('3250531c-537a-8131-b7a3-efe66af39c4f', 'Sprint DB');
  await getDB('3250531c-537a-8109-aa3d-cd3961fe0257', 'Tasks DB');
}

main();
