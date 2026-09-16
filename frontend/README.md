# Calculator Frontend

An Apple Calculator-styled web UI (React, TypeScript, Vite) that
consumes the [backend REST API](../backend/README.md).

## Requirements

- Node.js 18 or newer
- The backend running locally (see [backend/README.md](../backend/README.md))

## Setup

```bash
cd frontend
cp .env.example .env
npm install
```

## Running the dev server

```bash
npm run dev
```

Vite's dev server defaults to `http://localhost:5173`. **The backend's
`CORS_ALLOWED_ORIGINS` must include this origin** (or `*`) for requests
from the dev server to succeed — see the backend README's Configuration
section. Without it, every API call fails with a CORS error in the
browser console instead of reaching the calculator.

## Running the tests

```bash
npm test          # single run
npm run test:watch  # watch mode
```

## Building for production

```bash
npm run build     # type-checks with tsc, then builds with Vite into dist/
npm run preview   # serve the production build locally
```

## Environment variables

| Variable | Default (`.env.example`) | Description |
|---|---|---|
| `VITE_API_BASE_URL` | `http://localhost:8080` | Base URL of the backend REST API |

Vite only exposes variables prefixed with `VITE_` to client code, and
only reads `.env` (never committed) — `.env.example` is the template.

## Design decisions

### Model / Controller / View mapping

This project mirrors the backend's MVC-inspired separation, but each
layer is named with the vocabulary its own ecosystem actually uses —
the same principle the backend followed by calling its layer `handler`
instead of borrowing `controller` from another framework:

| Role | Backend | Frontend | Why this name |
|---|---|---|---|
| Model | `internal/calculator` | [`src/model/`](src/model/) | Represents the data contract. The backend computes; the frontend's model only *types* that contract (`types.ts`) and calls whoever computes (`apiClient.ts`) — no calculation logic lives here. |
| Controller | `internal/handler` | [`src/hooks/useCalculator.ts`](src/hooks/useCalculator.ts) | A **hook** is React's own mechanism for extracting stateful logic away from rendering — the direct equivalent of what a controller does in an MVC framework, expressed in the vocabulary the React ecosystem actually uses. |
| View | `internal/response` + routing | [`src/components/`](src/components/) | Purely presentational **components**: they render props and forward events, with no calculator logic of their own. |

`App.tsx` is the one place that connects the hook to the view — it
calls `useCalculator()` and passes the result down as props to
`<Calculator>`. `Calculator`, `Display`, `Keypad` and `Button` never
import the hook or the API client themselves, which keeps them trivial
to test in isolation with plain props.

### Why a ref-backed state container instead of plain `useState` updaters

`useCalculator` needs to read "the state as of right now" synchronously
inside `chooseOperator`/`equals`/`sqrt` (to decide what to send to the
API), then write the result back after an `await`. Doing this by
stuffing a side effect inside a `setState(prev => ...)` updater and
reading it back immediately after calling `setState` is a common
pattern, but it's fragile: React does not guarantee the updater runs
synchronously at call time. The hook instead keeps state in a `useRef`
(always synchronously read/write, never batched) and uses a dummy
`useState` purely to force a re-render after each change. This trades a
little indirection for state transitions that are deterministic
regardless of how React happens to batch the surrounding calls.

### Calculator state machine

The keypad's four "function" keys (`AC`, `√`, `%`, `xʸ`) plus the
standard `+ − × ÷` map onto a small state machine tracking the current
display string, a pending operand/operator pair, and whether the next
digit should start a fresh number (`waitingForOperand`):

- `add`, `subtract`, `multiply`, `divide`, `power` and `percentage` are
  all treated as binary operators with the same flow: number → operator
  → number → `=`. `percentage` follows the backend's own semantics
  (`percentage`% of `value`), matching the `%` key on a standard
  calculator rather than "what percent is A of B."
- `sqrt` is unary and applies immediately to the currently displayed
  number, without waiting for `=`.
- Pressing an operator while a second number has already been typed
  (but `=` hasn't been pressed yet) evaluates the pending operation
  first, so operators can be chained (`2 + 3 ×` computes `5` before
  setting up the multiplication).
- A decimal point is rejected if the current entry already has one, and
  digit entry is capped at 15 characters — both are frontend validation
  on top of whatever the backend re-validates.
- `backspace()` erases the last character of the current entry. It only
  exists as a controller method (not a state-machine transition change)
  and is reachable both by clicking the ⌫ icon in the display's corner
  and by pressing the physical Backspace key.

### Keyboard support

`src/hooks/useKeyboardInput.ts` attaches a single `window` `keydown`
listener (wired up once, from `App.tsx`) that calls the exact same
`useCalculator` methods the Keypad's buttons call — nothing is
reimplemented for the keyboard path. Mapping: digits `0`-`9`, `.` for
the decimal point, `+ - * / %` for their operators (`x` also works for
multiply), `Enter` or `=` for equals, `Backspace` to erase, and
`Escape` to clear. `xʸ` (power) and `√` (sqrt) are deliberately left
click-only: neither has a single, unambiguous physical key. Only the
keys this hook actually recognizes call `preventDefault()` (so, for
example, `/` doesn't trigger Firefox's quick-find), and the listener is
removed on unmount.

### Number formatting

`src/model/displayFormat.ts` exports a single pure function,
`formatDisplayNumber(rawValue: string): string`, used for both the main
display value and the pending-operation line (e.g. `"1,000 +"`) so the
two can never format differently. It takes the exact string from
calculator state — not a parsed number — so it can add thousands
separators to the integer part without disturbing a decimal point the
user is still typing (`"1234."` stays `"1,234."`) or trailing zeros
already present (`"1000.50"` stays `"1,000.50"`). Once the combined
integer+decimal digit count exceeds 15 — the largest digit count for
which every value is still an exactly representable double
(`Number.MAX_SAFE_INTEGER` is 16 digits, and not every 16-digit integer
is exact) — it switches to scientific notation with 6 significant
digits instead (`toExponential(5)`). `Display` additionally picks a
smaller font tier as the formatted string gets longer, since a
15-digit number with separators (19 characters) is longer than the
worst-case scientific notation string.

### Error handling

Every one of the backend's error codes (`INVALID_JSON`, `MISSING_FIELD`,
`INVALID_NUMBER`, `DIVISION_BY_ZERO`, `NEGATIVE_SQRT_INPUT`,
`RESULT_NOT_FINITE`, `NOT_FOUND`, `METHOD_NOT_ALLOWED`,
`INTERNAL_ERROR`) is mapped to its own distinct, human-readable message
in `useCalculator.ts`, plus a frontend-only `NETWORK_ERROR` for when the
backend can't be reached at all (e.g. it isn't running, or CORS
blocked the request). None of these collapse into a generic "Error" —
the `Display` component shows the specific message in place of the
numeric value. Pressing any digit after an error clears it and starts a
fresh entry.

### Styling

Plain CSS Modules, one stylesheet per component, no UI kit or utility
framework — the Apple Calculator look (strict black/gray/white palette,
circular buttons, no accent colors) is specific enough that a component
library would fight against it more than help. Sizes use `clamp()` so
the display and keypad scale down on narrow viewports without a
separate mobile layout.
