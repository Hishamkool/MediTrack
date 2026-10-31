# MediTrack

Medicine master, batch/expiry tracking and dealer purchase orders for a hospital pharmacy.

- `backend/`: Go API (chi, pgx, golang-migrate, JWT)
- `frontend/`: React + Vite + TypeScript + Tailwind

## Local setup

Prerequisites: Go 1.22+, Node 20+, and either Docker (for a local Postgres) or a Neon database.

### 1. Database

Option A, local Docker Postgres on port 5433:

```bash
docker compose up -d
```

Option B, Neon: create a project, copy its connection string (it ends in `?sslmode=require`), and use it as `DATABASE_URL` below.

### 2. Backend

```bash
cd backend
cp .env.example .env      # then edit: DATABASE_URL, JWT_SECRET, ADMIN_*
make run                  # or: go run ./cmd/api
```

On start the API:
1. applies any pending migrations (turn this off with `AUTO_MIGRATE=false`)
2. creates the first admin from `ADMIN_EMAIL` / `ADMIN_PASSWORD` if no active admin exists
3. listens on `http://localhost:8080` (check `http://localhost:8080/healthz`)

Other commands (run them from `backend/`):

| Command | What it does |
|---|---|
| `make test` | run the Go tests |
| `make migrate-up` | apply pending migrations |
| `make migrate-down` | roll back the last migration |
| `make migrate-version` | show the current migration version |

No `make`? Run the command shown in the `Makefile` directly, e.g. `go run ./cmd/migrate up`.

### 3. Frontend

```bash
cd frontend
cp .env.example .env.local   # VITE_API_URL=http://localhost:8080
npm install
npm run dev                  # http://localhost:5173
```

## Roles

| | admin | staff |
|---|---|---|
| Manage users | ✅ | ❌ |
| Medicines: add/edit | ✅ | ✅ |
| Dealers: view | ✅ | ✅ |
| Dealers: add/edit | ✅ | ❌ |
| Batches: add/edit | ✅ | ✅ |
| Dashboard, view orders | ✅ | ✅ |
| Create/send/commit orders | ✅ | ❌ |
| Delete records | ✅ | ❌ |

These rules are enforced by backend middleware. The frontend only hides what a user can't use.

## Deployment (summary, full steps come in Phase 5)

- **API (Render)**: Web Service, root directory `backend`, build command `go build -o bin/api ./cmd/api`, start command `./bin/api`, health check path `/healthz`. Set `DATABASE_URL`, `JWT_SECRET`, `CORS_ALLOWED_ORIGINS` (your frontend URL) and `ADMIN_*` as environment variables.
- **Frontend (Vercel/Netlify)**: root directory `frontend`, build command `npm run build`, output directory `dist`, env var `VITE_API_URL` = your API URL. `vercel.json` / `public/_redirects` make page refreshes on routes like `/users` work.
