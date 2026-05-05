# skild - Documento MVP

## 1) Objetivo

`skild` es un CLI en Go para instalar skills de OpenCode desde un unico repositorio remoto configurable por variables de entorno.

El foco del MVP es simplicidad, aprendizaje y ejecucion confiable para uso interno.

---

## 2) Alcance del MVP

### Incluye
- Soporte para **OpenCode unicamente**.
- Fuente de skills desde **un solo repositorio git**.
- Configuracion por variables de entorno.
- Comandos basicos:
  - `skild list`
  - `skild install <skill-name>`
  - `skild install --all`
  - `skild update`
  - `skild version`
  - `skild help`
- Instalacion por defecto en modo `copy`.
- Descubrimiento de skills por presencia de `SKILL.md`.

### No incluye (fuera del MVP)
- Multi-agent (Claude, Codex, Cursor, etc.).
- Multiples providers/fuentes simultaneas.
- Lockfile complejo o versionado avanzado.
- Telemetria y auditorias de seguridad.
- UI interactiva avanzada (fuzzy prompts).

---

## 3) Supuestos de estructura del repo de skills

Se asume un repositorio con una ruta raiz (por defecto `skills/`) y subcarpetas por skill:

```text
skills/
  skill-a/
    SKILL.md
  skill-b/
    SKILL.md
```

Cada skill valida debe tener `SKILL.md`.

---

## 4) Configuracion por variables de entorno

- `SKILD_REPO_URL` (**requerida**): URL del repo git de skills.
- `SKILD_REPO_REF` (opcional, default: `main`): branch/tag/commit.
- `SKILD_CACHE_DIR` (opcional, default: `~/.cache/skild`): directorio de cache local.
- `SKILD_ROOT_PATH` (opcional, default: `skills`): subdirectorio dentro del repo donde viven las skills.
- `SKILD_OPENCODE_DIR` (opcional, default: `./.opencode/skills`): destino de instalacion.
- `SKILD_INSTALL_MODE` (opcional, default: `copy`): `copy` o `symlink` (MVP recomienda `copy`).

---

## 5) Comportamiento esperado por comando

## `skild list`
- Sincroniza/accede al repo cacheado.
- Descubre skills validas.
- Muestra lista de nombres disponibles.

## `skild install <skill-name>`
- Valida existencia de la skill.
- Instala en `SKILD_OPENCODE_DIR/<skill-name-sanitized>`.
- Si ya existe, reemplaza instalacion anterior de forma segura.

## `skild install --all`
- Instala todas las skills descubiertas.

## `skild update`
- Actualiza repo cacheado (`fetch` + checkout de `SKILD_REPO_REF`).
- Reinstala skills instaladas previamente (o todas, segun simplificacion elegida para MVP).

## `skild version`
- Muestra version del binario.

## `skild help`
- Muestra uso y ejemplos.

---

## 6) Flujo tecnico MVP (alto nivel)

1. Cargar configuracion desde env + defaults.
2. Validar `SKILD_REPO_URL`.
3. Preparar cache:
   - Si no existe clone local: clonar.
   - Si existe: fetch + checkout ref.
4. Descubrir skills:
   - Buscar subdirectorios en `SKILD_ROOT_PATH` con `SKILL.md`.
5. Instalar:
   - Sanitizar nombre para evitar traversal o rutas invalidas.
   - Copiar carpeta skill al destino OpenCode.
6. Reportar resultado por consola.

---

## 7) Estructura sugerida del proyecto Go

```text
cmd/skild/main.go
internal/config/
internal/repo/
internal/discovery/
internal/install/
internal/cli/
internal/errors/
test/
```

Responsabilidades:
- `config`: env vars, defaults, validacion.
- `repo`: clone/fetch/checkout.
- `discovery`: detectar skills validas.
- `install`: copy/symlink, reemplazo seguro, sanitizacion.
- `cli`: parseo de comandos y handlers.

---

## 8) Reglas de seguridad minimas

- Nunca permitir path traversal (`..`) en nombres de skill o rutas derivadas.
- Sanitizar nombres de carpeta destino.
- Validar que operaciones de escritura queden dentro de `SKILD_OPENCODE_DIR`.
- Ignorar carpetas sin `SKILL.md`.
- Mensajes de error claros y accionables.

---

## 9) Criterios de aceptacion (Definition of Done)

1. Con `SKILD_REPO_URL` definida, `skild list` funciona.
2. `skild install <skill>` instala correctamente en `.opencode/skills/`.
3. `skild install --all` instala todas las skills detectadas.
4. `skild update` trae cambios del repo y actualiza instalacion.
5. Funciona en macOS/Linux (Windows opcional para siguiente fase).

---

## 10) Decisiones MVP

- **Modo default:** `copy` (menos friccion que symlink en multiplataforma).
- **Fuente unica:** 1 repo configurable.
- **Compatibilidad inicial:** solo OpenCode.
- **Deteccion skill:** basada en carpeta + `SKILL.md`.

---

## 11) Roadmap post-MVP

1. `skild remove` y `skild installed`.
2. Lockfile local (`skills-lock.json`) para reproducibilidad.
3. Soporte de repos privados (`GITHUB_TOKEN`).
4. Soporte multi-agent.
5. Validacion de esquema/frontmatter de `SKILL.md`.

---

## 12) Ejemplo de uso

```bash
export SKILD_REPO_URL="https://github.com/tu-org/company-skills.git"
export SKILD_REPO_REF="main"

skild list
skild install prompt-engineering
skild install --all
skild update
```
