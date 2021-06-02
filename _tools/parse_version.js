const fs = require("fs/promises");

async function parse() {
  const content = await fs.readFile("version.go", "utf8");
  const regex = new RegExp(`^const +version *= *"(?<version>.+)?"`, "m");
  const version = regex.exec(content)?.groups?.version;
  if (version === undefined) {
    process.on("exit", function() {
      process.exit(1);
    });
  } else {
    console.log("::set-output name=VERSION::" + version);
  }
}

parse();
