// Types mirroring the backend's request/response contract. This file
// has no logic of its own — it only describes the shape of the API.

export interface BinaryRequest {
  a: number;
  b: number;
}

export interface PowerRequest {
  base: number;
  exponent: number;
}

export interface SqrtRequest {
  value: number;
}

export interface PercentageRequest {
  value: number;
  percentage: number;
}

export interface SuccessResponse {
  success: true;
  operation: string;
  result: number;
}

// BackendErrorCode lists every error code the backend can return, per
// its README.
export type BackendErrorCode =
  | 'INVALID_JSON'
  | 'MISSING_FIELD'
  | 'INVALID_NUMBER'
  | 'DIVISION_BY_ZERO'
  | 'NEGATIVE_SQRT_INPUT'
  | 'RESULT_NOT_FINITE'
  | 'NOT_FOUND'
  | 'METHOD_NOT_ALLOWED'
  | 'INTERNAL_ERROR';

export interface BackendErrorDetail {
  code: BackendErrorCode;
  message: string;
  field?: string;
}

export interface ErrorResponse {
  success: false;
  error: BackendErrorDetail;
}

// NetworkErrorDetail is a client-side addition for failures that never
// reached the backend at all (the server is unreachable, or the
// response could not be parsed).
export interface NetworkErrorDetail {
  code: 'NETWORK_ERROR';
  message: string;
}

export type ApiError = BackendErrorDetail | NetworkErrorDetail;

export type ApiResult =
  | { ok: true; operation: string; result: number }
  | { ok: false; error: ApiError };
