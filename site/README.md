# Open Curriculum — ChatGPT Site

Prototipo funcional del concepto descrito en el repositorio original. Está preparado para desplegarse como un **ChatGPT Site** con almacenamiento persistente gestionado por la plataforma:

- **D1 (`DB`)**: currículo, dependencias, propuestas, votos, comentarios y progreso.
- **R2 (`FILES`)**: documentos, imágenes, audio y vídeo asociados a propuestas.
- **Sign in with ChatGPT**: identidad de participantes y autorización de mantenimiento.
- **Sin `localStorage` ni IndexedDB**: el navegador no es la fuente de verdad.

## Funcionalidad incluida

- Currículo navegable como grafo dirigido de unidades y prerrequisitos.
- Selección de una unidad objetivo y resaltado de la ruta necesaria.
- Filtros por texto y área temática.
- Progreso personal por unidad (`in_progress` / `completed`).
- Propuestas para crear o editar unidades.
- Cambios de prerrequisitos con detección de ciclos en servidor.
- Votos, abstenciones y comentarios por propuesta.
- Archivos opcionales guardados en R2; solo se publican al aceptar la propuesta.
- Aceptación o rechazo por personas mantenedoras configuradas.
- Aplicación atómica del cambio aceptado y registro básico de auditoría.
- Datos iniciales para poder evaluar el producto desde el primer despliegue.

## Estructura compatible con Sites

El comando de build genera:

```text
dist/
├── .openai/hosting.json
├── index.html
├── assets/
│   ├── app.js
│   └── styles.css
└── server/
    ├── index.js
    └── lib.js
```

`dist/server/index.js` exporta el Worker que sirve la API y delega los estáticos en `env.ASSETS`.

## Comprobación local

No hay dependencias de runtime ni de build.

```bash
npm run verify
```

Esto comprueba la sintaxis, ejecuta las pruebas y construye `dist/`.

## Despliegue con ChatGPT Sites

1. Abre esta carpeta `site/` como proyecto local en ChatGPT Work o Codex de escritorio.
2. Pide: `Deploy this existing project with Sites. Provision D1 as DB and R2 as FILES. Save a version before deploying.`
3. Revisa el diff y la migración antes de desplegar.
4. En **Site settings**, configura:

```text
MAINTAINER_EMAILS=correo1@example.com,correo2@example.com
```

5. Guarda una versión y, después de revisarla, despliega esa versión.

`.openai/hosting.json` comienza sin `project_id`; Sites lo añadirá al vincular el proyecto alojado.

## Inicialización de datos

El Worker ejecuta `CREATE TABLE IF NOT EXISTS` y una semilla idempotente la primera vez que se consulta la API. La definición equivalente se conserva en `migrations/0001_initial.sql` para revisión, evolución y migraciones futuras.

## Decisiones deliberadas

- La vista pública solo muestra unidades ya aceptadas.
- Las personas anónimas pueden explorar; contribuir y guardar progreso requiere inicio de sesión.
- La decisión final es manual por mantenimiento. El criterio estadístico de cierre de votaciones descrito en el manifiesto se deja para una fase posterior, cuando existan identidad robusta, reputación y suficiente volumen de votos.
- Los recursos se adjuntan a una propuesta y no se hacen públicos hasta que la propuesta se acepta.
- La certificación, los grupos docentes, los ejercicios corregidos y la ponderación de votos quedan fuera de este primer corte vertical.

## Variables y bindings

| Nombre | Tipo | Obligatorio | Uso |
|---|---|---:|---|
| `DB` | D1 binding | Sí | Estado estructurado persistente |
| `FILES` | R2 binding | Sí para adjuntos | Contenido binario persistente |
| `ASSETS` | Sites asset binding | Sí | Frontend compilado |
| `MAINTAINER_EMAILS` | Variable de entorno | Sí para resolver propuestas | Lista CSV de emails autorizados |

No guardes secretos en `.openai/hosting.json` ni en el repositorio.
