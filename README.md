# Calculator

A calculator built as two independent projects in one repository: a Go
REST API and a React frontend that consumes it.

- [`backend/`](backend/README.md) — the REST API (Go, standard library
  only). Implements add, subtract, multiply, divide, power, sqrt and
  percentage as separate endpoints, with strict input validation and a
  consistent JSON envelope for both results and errors.
- [`frontend/`](frontend/README.md) — an Apple Calculator-styled web UI
  (React, TypeScript, Vite) that consumes the API above.

See each project's own README for setup, how to run it, its API/UI
details, and its design decisions.

## Quick start

```bash
# terminal 1: backend
cd backend
cp .env.example .env   # add http://localhost:5173 to CORS_ALLOWED_ORIGINS for local frontend dev
go run ./cmd/api

# terminal 2: frontend
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Deployment / Docker

Prerequisites: [Git](https://git-scm.com/downloads) and
[Docker Desktop](https://www.docker.com/products/docker-desktop/)
(includes Docker Compose v2). Nothing else — this does not require
installing Go or Node locally.

1. Clone the repository and move into its root:

   ```bash
   git clone <this-repository>
   cd Calculator
   ```

2. Build and start both services with one command:

   ```bash
   docker compose up --build
   ```

3. Open http://localhost:3000 in a browser and use the calculator.

Docker Compose builds the backend and frontend images and starts both
containers, each on its default port:

| Service | URL |
|---|---|
| Frontend | http://localhost:3000 |
| Backend | http://localhost:8080 |

### If port 8080 or 3000 is already in use

Create a `.env` file in the repository root, next to
`docker-compose.yml` (this is separate from `backend/.env` and
`frontend/.env`, which only matter for running each service locally
without Docker — see each project's own README), setting whichever
port(s) you need to change:

```
BACKEND_PORT=9090
FRONTEND_PORT=4000
```

Then run `docker compose up --build` again. The `--build` is required,
not optional: the frontend image bakes the backend's URL directly into
its compiled JavaScript when the image is built (see below), so a
changed `BACKEND_PORT` only takes effect on a fresh build, not on a
plain `docker compose up`.

### Stopping

```bash
docker compose down
```

**Why `VITE_API_BASE_URL` is a build arg, not a runtime setting:** Vite
inlines every `import.meta.env.VITE_*` value directly into the compiled
JavaScript at build time — there is no Node process left at runtime to
read an environment variable from, since the final frontend image is
just nginx serving static files. `docker-compose.yml` passes it as a
build `arg` pointing at `http://localhost:${BACKEND_PORT:-8080}` — the
backend's port **as your browser sees it** (mapped to the host), not
`http://backend:8080` (the internal Compose service name/network,
which only resolves between containers, not from a browser running on
your machine). The same logic applies in reverse to the backend's
`CORS_ALLOWED_ORIGINS`, which is set to
`http://localhost:${FRONTEND_PORT:-3000}` — the origin the browser will
actually load the frontend from.

If you change `BACKEND_PORT` or `FRONTEND_PORT`, both of these stay
consistent automatically, since they're derived from the same
variables rather than hardcoded.

## Prompts used

The prompts used to develop this project are documented in
[PROMPTS.txt](PROMPTS.txt).
