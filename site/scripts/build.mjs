import { cp, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";

async function concatenate(directory, prefix, output) {
  const names = (await readdir(directory))
    .filter((name) => name.startsWith(prefix) && name.endsWith(".js"))
    .sort();
  if (!names.length) throw new Error(`No source fragments found for ${directory}/${prefix}`);
  const chunks = await Promise.all(names.map((name) => readFile(`${directory}/${name}`, "utf8")));
  await writeFile(output, chunks.join(""));
}

await rm("dist", { recursive: true, force: true });
await mkdir("dist/assets", { recursive: true });
await mkdir("dist/server", { recursive: true });
await mkdir("dist/.openai", { recursive: true });

await cp("index.html", "dist/index.html");
await cp("src/styles.css", "dist/assets/styles.css");
await concatenate("src", "app.part", "dist/assets/app.js");
await concatenate("worker", "index.part", "dist/server/index.js");
await cp("worker/lib.js", "dist/server/lib.js");
await cp(".openai/hosting.json", "dist/.openai/hosting.json");
