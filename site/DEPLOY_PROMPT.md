# Prompt de despliegue recomendado

```text
Deploy this existing project as a ChatGPT Site.

Treat the repository's `site/` directory as the Site project root. Preserve the current product behavior and visual design. Confirm that the generated artifacts are compatible with the current Sites runtime.

Provision durable storage:
- D1 relational database bound as `DB`.
- R2 object storage bound as `FILES`.

The app must not use localStorage or IndexedDB for curriculum data, proposals, votes, comments, progress, or uploaded files. Keep all authorization in server-side code. Preserve Sign in with ChatGPT using the platform paths and identity headers.

Review and apply `migrations/0001_initial.sql`, or confirm that the idempotent runtime bootstrap creates the same schema. Set `MAINTAINER_EMAILS` in Site settings rather than committing it.

Run `npm run verify`. Show me the source diff, storage bindings, migration plan, and access model. Save a version first. Do not deploy the version publicly until I explicitly approve the reviewed version.
```
