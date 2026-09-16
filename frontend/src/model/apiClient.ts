// Thin wrapper around fetch: one function per backend operation. It
// never throws — every call resolves to a typed ApiResult, so callers
// (the useCalculator hook) never need a try/catch of their own.
import type {
  ApiResult,
  BinaryRequest,
  ErrorResponse,
  PercentageRequest,
  PowerRequest,
  SqrtRequest,
  SuccessResponse,
} from './types';

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

async function post<TBody>(path: string, body: TBody): Promise<ApiResult> {
  let response: Response;
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
  } catch {
    return {
      ok: false,
      error: {
        code: 'NETWORK_ERROR',
        message: 'Unable to reach the server. Check your connection and try again.',
      },
    };
  }

  let payload: SuccessResponse | ErrorResponse;
  try {
    payload = await response.json();
  } catch {
    return {
      ok: false,
      error: { code: 'NETWORK_ERROR', message: 'The server returned an unreadable response.' },
    };
  }

  if (!payload.success) {
    return { ok: false, error: payload.error };
  }
  return { ok: true, operation: payload.operation, result: payload.result };
}

export const apiClient = {
  add: (a: number, b: number) => post<BinaryRequest>('/api/v1/add', { a, b }),
  subtract: (a: number, b: number) => post<BinaryRequest>('/api/v1/subtract', { a, b }),
  multiply: (a: number, b: number) => post<BinaryRequest>('/api/v1/multiply', { a, b }),
  divide: (a: number, b: number) => post<BinaryRequest>('/api/v1/divide', { a, b }),
  power: (base: number, exponent: number) => post<PowerRequest>('/api/v1/power', { base, exponent }),
  sqrt: (value: number) => post<SqrtRequest>('/api/v1/sqrt', { value }),
  percentage: (value: number, percentage: number) =>
    post<PercentageRequest>('/api/v1/percentage', { value, percentage }),
};
