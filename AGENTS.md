# Agent Notes

Compact, high-signal guidance for working in this repo.

## Stack & Layout

- Monorepo: `backend/` (Go 1.26.3, Gin, GORM, PostgreSQL) and `frontend/` (SvelteKit 2, Svelte 5 runes, Vite 8, Tailwind CSS 4).
- Root `.env` is committed and loaded by the backend via `godotenv`.
- Backend entrypoint: `backend/main.go`. API routes are grouped under `/api`.
- Frontend is client-only: `frontend/src/routes/+layout.ts` sets `ssr = false`.

## Running Locally

Use Docker Compose for the full stack:

```bash
docker compose up --build
```

- PostgreSQL: `localhost:5432`
- Backend dev server (with Air hot-reload): `http://localhost:8080`
- Frontend dev server: `http://localhost:5173`

Backend Air config enables polling (`poll = true`) because the source is bind-mounted into the container.

## Backend

### Commands

```bash
cd backend

# Build binary
go build -o backend .

# Format (PR checklist requirement)
go fmt ./...

# Run tests (some require Postgres — see Testing)
go test ./...

# Run a single package
go test ./internal/inventario
```

### Important Runtime Behavior

- `backend/internal/database/database.go` runs GORM `AutoMigrate` on startup.
- `backend/internal/config/setup.go` seeds initial data if tables are empty:
  - Roles: `Admin`, `Empleado`
  - Default admin user: `admin` / `Admin123.`
  - Default cash register: `Caja Principal`
  - Payment methods: `Efectivo`, `Tarjeta`, `Fiado`
- `config.LoadConfig()` fails fatally if any `DB_*` variable or `JWT_SECRET` is missing.

### Package Conventions

Each domain lives under `backend/internal/<domain>/` with roughly:

- `models.go` / `modelos.go` — GORM models
- `repository.go` / `repositorio.go` — data access
- `controller.go` / `controladores.go` — HTTP handlers
- `routes.go` — route wiring

Existing domains: `auth`, `cajas`, `clientes`, `empleados`, `inventario`, `promociones`, `ventas`.

## Frontend

### Commands

```bash
cd frontend
npm install

# Dev server
npm run dev

# Type-check Svelte/TS
npm run check

# Lint = prettier check + eslint
npm run lint

# Auto-format
npm run format
```

### Important Conventions

- **Runes mode is forced** in `vite.config.ts` for all non-`node_modules` files. Use `$state`, `$derived`, `$props`, etc.
- **Tailwind 4 is CSS-first**: imported in `src/routes/layout.css` with `@import 'tailwindcss'` and `@theme` custom properties.
- Prettier config is in `frontend/.prettierrc`: tabs, single quotes, no trailing commas, `printWidth: 100`, Svelte + Tailwind plugins.
- API base URL comes from `import.meta.env.VITE_API_URL` (fallback `http://localhost:8080/api`). Calls use `credentials: 'include'`.
- `$lib/api.ts` is the central API client. Error messages come from fields like `detalle`, `error`, `mensaje`.

## Testing

- Backend tests are in `backend/internal/inventario/`:
  - `controller_test.go` — unit tests with `gin.TestMode`, no DB needed.
  - `repository_test.go` — **requires a local Postgres** with hardcoded DSN:
    ```
    host=localhost user=postgres password=admin123 dbname=pos_db port=5432 sslmode=disable
    ```
    It auto-migrates `Categoria` and `Producto` for the test.
- No frontend test framework is configured.

## Deployment

- `.github/workflows/deploy-cd.yaml` deploys on every push to `main`.
- It SSHs into the production server, pulls `main`, builds `gps-backend` in `~/GPS_2026_1/backend`, kills the previous process, and restarts it with `nohup`.
- Only the backend Go binary is deployed via CD; the frontend is served through the Docker/compose dev flow locally.

## PR Checklist (from `.github/pull_request_template.md`)

Before requesting review:

- Code compiles locally.
- Run `go fmt ./...`.
- Add validations to prevent corrupt DB data.
- Add `_test.go` unit tests when applicable.
