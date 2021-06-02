const fs = require("fs/promises");

(async function main() {
  const content = await fs.readFile("version.go", "utf8");
  const regex = new RegExp(`^const +version *= *"(?<version>.+)?"`, "m");
  const version = regex.exec(content)?.groups?.version;
  if (version === undefined) {
    throw new Error("version is undefined.");
  }
  console.log("::set-output name=VERSION::" + version);
})().catch(error => {
  console.error(error);
  process.exit(1);
});
