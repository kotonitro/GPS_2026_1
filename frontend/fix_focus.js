const fs = require('fs');
const path = require('path');

function walk(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  list.forEach(file => {
    file = path.join(dir, file);
    const stat = fs.statSync(file);
    if (stat && stat.isDirectory()) {
      results = results.concat(walk(file));
    } else if (file.endsWith('+page.svelte')) {
      results.push(file);
    }
  });
  return results;
}

const files = walk('c:/Users/naits/GPS_2026_1/frontend/src/routes/dashboard');
let changed = 0;

files.forEach(f => {
  let content = fs.readFileSync(f, 'utf8');
  // First, temporarily remove any focus:ring-1 focus:ring-primario that follows focus:outline-none
  // so we don't duplicate it.
  content = content.replace(/focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario/g, 'focus:border-primario focus:outline-none');
  
  // Now replace all focus:border-primario focus:outline-none with the full string
  const newContent = content.replace(/focus:border-primario focus:outline-none/g, 'focus:border-primario focus:outline-none focus:ring-1 focus:ring-primario');
  
  if (newContent !== fs.readFileSync(f, 'utf8')) {
    fs.writeFileSync(f, newContent);
    changed++;
    console.log('Fixed', f);
  }
});
console.log(`Updated ${changed} files.`);
