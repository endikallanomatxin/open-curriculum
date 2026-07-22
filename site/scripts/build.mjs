import { cp, mkdir, rm } from "node:fs/promises";

await rm("dist", { recursive: true, force: true });
await mkdir("dist/assets", { recursive: true });
await mkdir("dist/server", { recursive: true });
await mkdir("dist/.openai", { recursive: true });

await cp("index.html", "dist/index.html");
await cp("src/app.js", "dist/assets/app.js");
await cp("src/styles.css", "dist/assets/styles.css");
await cp("worker/index.js", "dist/server/index.js");
await cp("worker/lib.js", "dist/server/lib.js");
await cp(".openai/hosting.json", "dist/.openai/hosting.json");
