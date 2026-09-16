# Calculator API

A REST API for a calculator, written in Go with an MVC-style layout and a
strict separation of concerns between the arithmetic model, the HTTP
controllers, and the JSON presentation layer.

## Requirements

- Go 1.22 or newer
- No external dependencies (standard library only)

## Setup

This is the backend half of a monorepo; see the [root README](../README.md)
for how it relates to `frontend/`.

```bash
git clone <this-repository>
cd Calculator/backend
cp .env.example .env   # optional, see Configuration below
```

## Running the project

```bash
go run ./cmd/api
```

Run this from inside `backend/` (or point `go run` at `./backend/cmd/api`
from the repository root).

By default the server listens on port `8080`. Configuration is read from
environment variables (see [Configuration](#configuration)); on Windows
PowerShell you can set them inline, e.g.:

```powershell
$env:SERVER_PORT = "8080"; go run ./cmd/api
```

## Configuration

All configuration comes from environment variables, with sensible
defaults so the server also runs with none of them set. See
[.env.example](.env.example) for the full list:

| Variable | Default | Description |
|---|---|---|
| `SERVER_PORT` | `8080` | TCP port the HTTP server listens on |
| `CORS_ALLOWED_ORIGINS` | `*` | Comma-separated list of origins allowed by CORS (`*` allows any origin) |
| `CALC_PRECISION` | `2` | Number of decimal places every result is rounded to |

`cmd/api/main.go` loads `.env` (via `config.LoadDotEnv`, a small
stdlib-only parser — no third-party dotenv library) into the process
environment before reading configuration, so editing `.env` and
restarting the server is enough; a variable already set in the real
environment always takes precedence over the file. If running the
frontend locally, add its dev server origin (`http://localhost:5173`
by default) to `CORS_ALLOWED_ORIGINS`.

## Running the tests

```bash
go test ./... -cover
```

Latest run:

```
ok  	calculator/internal/calculator	coverage: 95.2% of statements
ok  	calculator/internal/config	coverage: 95.8% of statements
ok  	calculator/internal/handler	coverage: 86.6% of statements
ok  	calculator/internal/response	coverage: 100.0% of statements
ok  	calculator/internal/router	coverage: 100.0% of statements
```

(`cmd/api` has no tests of its own since `main.go` only wires dependencies
together and starts the server — there is no branching logic to cover.)

## API reference

Every endpoint is a `POST` under `/api/v1/`, accepts a JSON body, and
returns a JSON body with the same envelope shape.

**Success envelope:**

```json
{ "success": true, "operation": "<name>", "result": <number> }
```

**Error envelope:**

```json
{ "success": false, "error": { "code": "<CODE>", "message": "<text>", "field": "<optional>" } }
```

| Error code | HTTP status | Meaning |
|---|---|---|
| `INVALID_JSON` | 400 | Malformed or unparseable request body |
| `MISSING_FIELD` | 400 | A required field is missing |
| `INVALID_NUMBER` | 400 | A field is present but is not a number |
| `DIVISION_BY_ZERO` | 422 | Division (or an equivalent power case) by zero |
| `NEGATIVE_SQRT_INPUT` | 422 | Square root of a negative number |
| `RESULT_NOT_FINITE` | 422 | Result overflowed to `Inf`/`NaN` |
| `NOT_FOUND` | 404 | Unknown route |
| `METHOD_NOT_ALLOWED` | 405 | Known route, wrong HTTP method |
| `INTERNAL_ERROR` | 500 | Unexpected failure (should not occur for expected input) |

### POST /api/v1/add

```bash
curl -X POST http://localhost:8080/api/v1/add -d '{"a": 2, "b": 3.5}'
```

```json
{ "success": true, "operation": "add", "result": 5.5 }
```

### POST /api/v1/subtract

```bash
curl -X POST http://localhost:8080/api/v1/subtract -d '{"a": 5, "b": 3}'
```

```json
{ "success": true, "operation": "subtract", "result": 2 }
```

### POST /api/v1/multiply

```bash
curl -X POST http://localhost:8080/api/v1/multiply -d '{"a": 4, "b": 5}'
```

```json
{ "success": true, "operation": "multiply", "result": 20 }
```

### POST /api/v1/divide

```bash
curl -X POST http://localhost:8080/api/v1/divide -d '{"a": 10, "b": 0}'
```

```json
{ "success": false, "error": { "code": "DIVISION_BY_ZERO", "message": "division by zero" } }
```

### POST /api/v1/power

```bash
curl -X POST http://localhost:8080/api/v1/power -d '{"base": 2, "exponent": -1}'
```

```json
{ "success": true, "operation": "power", "result": 0.5 }
```

### POST /api/v1/sqrt

```bash
curl -X POST http://localhost:8080/api/v1/sqrt -d '{"value": -9}'
```

```json
{ "success": false, "error": { "code": "NEGATIVE_SQRT_INPUT", "message": "cannot compute the square root of a negative number" } }
```

### POST /api/v1/percentage

Computes `percentage`% of `value`.

```bash
curl -X POST http://localhost:8080/api/v1/percentage -d '{"value": 200, "percentage": 15}'
```

```json
{ "success": true, "operation": "percentage", "result": 30 }
```

### Validation error example (any endpoint)

```bash
curl -X POST http://localhost:8080/api/v1/add -d '{"a": 2}'
```

```json
{ "success": false, "error": { "code": "MISSING_FIELD", "message": "field \"b\" is required", "field": "b" } }
```

## Design decisions

- **One endpoint per operation, not a generic `/calculate`.** Each route
  has its own request shape and its own domain errors, so validation and
  error handling stay specific and simple instead of branching on an
  operator string.

- **Strict package boundaries.** `internal/calculator` never imports
  `net/http`; it is a pure arithmetic library that can be tested and
  reused with no server running. `internal/handler` never builds a JSON
  body directly — it only decides *which* error code and message apply,
  and delegates writing the response to `internal/response`. This keeps
  the JSON envelope consistent by construction: there is only one place
  in the codebase that calls `json.Marshal` on a response body.

- **Domain errors as sentinel values.** `calculator` returns
  `ErrDivisionByZero`, `ErrNegativeSqrtInput` and `ErrResultNotFinite`
  instead of generic errors or panics. Handlers map them to HTTP status
  codes with `errors.Is`, which keeps the mapping explicit and makes it
  impossible for an expected edge case to fall through to a 500.

- **Rounding strategy.** Every result is rounded with "round half away
  from zero" (`math.Round` after scaling by `10^precision`), applied once
  to the final result of each operation. The precision is a constructor
  parameter of `Calculator`, not a package-level constant, so it comes
  from configuration without coupling the model to environment loading.

- **Strict JSON decoding.** Request bodies are decoded into structs with
  `*float64` fields and `DisallowUnknownFields`. A `nil` pointer after
  decoding means the field was missing (`MISSING_FIELD`); a
  `json.UnmarshalTypeError` means the field had the wrong type
  (`INVALID_NUMBER`, naming the offending field); anything else
  (including unexpected extra fields) is treated as `INVALID_JSON`. The
  trade-off is that a client sending an unrecognized field gets a
  slightly generic "not valid JSON" message rather than a dedicated
  "unknown field" code — acceptable given how rarely that case matters
  in practice.

- **`0^0 = 1` convention.** `power` follows Go's own `math.Pow`
  convention for `0^0`, which matches the convention used by most
  calculators and programming languages, instead of treating it as an
  error.

- **`percentage` semantics.** `percentage` is defined as "`percentage`%
  of `value`" (`result = value * (percentage / 100)`), mirroring the `%`
  button on a standard calculator, rather than "what percent is A of B."

- **CORS lives in `router/`, not a separate `middleware/` package.**
  Since CORS is inseparable from how requests are routed and is the only
  cross-cutting HTTP concern this API needs, it did not justify its own
  package under the five defined responsibilities (model, handlers,
  presentation, routing, config).

- **404/405 share the same JSON envelope as every other error.** Method
  mismatches are caught by a small `withMethod` wrapper in `router/`
  before a handler ever runs, and unmatched routes fall through to a
  dedicated `notFound` handler — both go through `response.WriteError`,
  so a client never sees Go's default plain-text error pages.

The prompts used to develop this project are documented in
[PROMPTS.txt](../PROMPTS.txt) at the repository root.
