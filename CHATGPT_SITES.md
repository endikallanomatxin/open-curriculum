# Open Curriculum en ChatGPT Sites

Este repositorio conserva el prototipo original en Go/PostgreSQL y añade en `site/` un nuevo corte vertical pensado para ChatGPT Sites.

## Qué se ha trasladado del concepto original

El manifiesto y el prototipo describen cuatro ideas que ya forman un producto coherente:

1. Un currículo global representado como grafo dirigido.
2. Unidades de aprendizaje con dependencias explícitas y recursos reutilizables.
3. Una vista estable para estudiantes que solo muestra la última versión aceptada.
4. Un flujo de desarrollo abierto basado en propuestas, discusión, votación y aceptación.

La primera implementación de Sites cubre ese circuito completo. Funciones como certificación institucional, ofertas de empleo, ejercicios supervisados, reputación y ponderación estadística de votos se dejan para etapas posteriores.

## Arquitectura

```text
Navegador
  ├─ Frontend estático, sin estado de producto local
  ├─ Sign in with ChatGPT
  └─ API same-origin
        ├─ Worker de Sites
        ├─ D1: grafo, propuestas, votos, comentarios y progreso
        └─ R2: recursos adjuntos
```

La identidad se recibe exclusivamente mediante las cabeceras verificadas por Sites. El frontend nunca decide quién puede resolver una propuesta; la API compara el email autenticado con `MAINTAINER_EMAILS`.

## Modelo de publicación

- Las unidades de `units` son la versión aprobada y visible.
- Una propuesta contiene el estado deseado de una unidad y sus prerrequisitos.
- Los adjuntos permanecen privados mientras la propuesta esté abierta o rechazada.
- Al aceptar una propuesta se actualizan unidad, dependencias, adjuntos y resolución mediante un batch de D1.
- Antes de escribir dependencias se verifica en servidor que no aparezca un ciclo.

## Entrada al proyecto

Consulta [`site/README.md`](site/README.md) para validar, configurar y desplegar el Site.
