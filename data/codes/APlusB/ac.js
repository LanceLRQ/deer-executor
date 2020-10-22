const fs = require('fs');
const data = fs.readFileSync('/dev/stdin');
const result = data.toString('ascii').trim().split(' ').map(x => parseInt(x)).reduce((a, b) => a + b, 0);
console.log(result);
process.exit();