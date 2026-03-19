// Generate constants.go content from JSON files
import { readFileSync } from "fs";

const t3 = JSON.parse(readFileSync("constants_t3.json", "utf8"));
const t4 = JSON.parse(readFileSync("constants_t4.json", "utf8"));

function formatRC(data, varName) {
  const lines = [];
  lines.push(`var ${varName} = [${data.roundConstants.length}][${data.t}]string{`);
  for (const row of data.roundConstants) {
    lines.push(`\t{${row.map(v => `"${v}"`).join(", ")}},`);
  }
  lines.push(`}`);
  return lines.join("\n");
}

function formatMDS(data, varName) {
  const lines = [];
  lines.push(`var ${varName} = [${data.t}][${data.t}]string{`);
  for (const row of data.mds) {
    lines.push(`\t{${row.map(v => `"${v}"`).join(", ")}},`);
  }
  lines.push(`}`);
  return lines.join("\n");
}

console.log(`package poseidon

// Auto-generated Poseidon constants for P-256 base field.
// Generated from @noble/curves Grain LFSR.

`);
console.log(formatRC(t3, "rcHexT3"));
console.log();
console.log(formatMDS(t3, "mdsHexT3"));
console.log();
console.log(formatRC(t4, "rcHexT4"));
console.log();
console.log(formatMDS(t4, "mdsHexT4"));
