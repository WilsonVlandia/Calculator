// useCalculator is the controller: it owns the calculator's state and
// decides what happens on each button press, including when to call
// the API client and how to fold the result (or error) back into
// state. This is React's own vocabulary for the role a backend
// "handler" plays — a hook is how React extracts stateful logic away
// from rendering, the same way the backend used Go's own "handler"
// instead of borrowing "controller" from another framework.
import { useCallback, useRef, useState } from 'react';
import { apiClient } from '../model/apiClient';
import type { ApiError, ApiResult } from '../model/types';

export type BinaryOperatorKey = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'percentage';

// OPERATOR_SYMBOLS is the single source of truth for how each operator
// is displayed, shared by the Keypad's button labels and the pending
// operation row so the two can never drift apart.
export const OPERATOR_SYMBOLS: Record<BinaryOperatorKey, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: 'xʸ',
  percentage: '%',
};

interface CalculatorState {
  display: string;
  pendingValue: number | null;
  pendingOperator: BinaryOperatorKey | null;
  waitingForOperand: boolean;
  errorMessage: string | null;
  isLoading: boolean;
}

const INITIAL_STATE: CalculatorState = {
  display: '0',
  pendingValue: null,
  pendingOperator: null,
  waitingForOperand: false,
  errorMessage: null,
  isLoading: false,
};

// MAX_DIGITS keeps the display (and the numbers sent to the API) from
// growing without bound while the user is typing.
const MAX_DIGITS = 15;

// ERROR_MESSAGES gives every backend error code its own distinct,
// human-readable message instead of collapsing them into one generic
// "Error". NETWORK_ERROR falls back to the message apiClient already
// generated for it.
const ERROR_MESSAGES: Partial<Record<ApiError['code'], string>> = {
  DIVISION_BY_ZERO: 'Cannot divide by zero',
  NEGATIVE_SQRT_INPUT: 'Cannot take the square root of a negative number',
  RESULT_NOT_FINITE: 'Result is too large to display',
  INVALID_NUMBER: 'Enter a valid number',
  MISSING_FIELD: 'A required value is missing',
  INVALID_JSON: 'The request could not be understood by the server',
  NOT_FOUND: 'Requested operation was not found',
  METHOD_NOT_ALLOWED: 'Requested operation is not allowed',
  INTERNAL_ERROR: 'Something went wrong. Please try again',
};

function errorMessageFor(error: ApiError): string {
  return ERROR_MESSAGES[error.code] ?? error.message;
}

function callForOperator(operator: BinaryOperatorKey, a: number, b: number): Promise<ApiResult> {
  switch (operator) {
    case 'add':
      return apiClient.add(a, b);
    case 'subtract':
      return apiClient.subtract(a, b);
    case 'multiply':
      return apiClient.multiply(a, b);
    case 'divide':
      return apiClient.divide(a, b);
    case 'power':
      return apiClient.power(a, b);
    case 'percentage':
      return apiClient.percentage(a, b);
  }
}

export function useCalculator() {
  // The async handlers below (chooseOperator, equals, sqrt) need to
  // read "the state as of right now" before awaiting an API call, and
  // again read/write it after the call resolves. Relying on React's
  // setState-updater form for that would depend on exactly when React
  // chooses to invoke the updater relative to the surrounding code,
  // which is not guaranteed. Keeping the actual value in a ref and
  // using forceRender purely to trigger a re-render sidesteps that
  // ambiguity: reads and writes are always synchronous and immediate.
  const stateRef = useRef<CalculatorState>(INITIAL_STATE);
  const [, forceRender] = useState(0);

  const setState = useCallback((updater: CalculatorState | ((prev: CalculatorState) => CalculatorState)) => {
    const next = typeof updater === 'function' ? (updater as (prev: CalculatorState) => CalculatorState)(stateRef.current) : updater;
    stateRef.current = next;
    forceRender((n) => n + 1);
  }, []);

  const inputDigit = useCallback(
    (digit: string) => {
      setState((prev) => {
        if (prev.isLoading) return prev;
        if (prev.errorMessage) return { ...INITIAL_STATE, display: digit };
        if (prev.waitingForOperand) return { ...prev, display: digit, waitingForOperand: false };

        const digitsOnly = prev.display.replace(/[-.]/g, '');
        if (digitsOnly.length >= MAX_DIGITS) return prev;

        const nextDisplay = prev.display === '0' ? digit : prev.display + digit;
        return { ...prev, display: nextDisplay };
      });
    },
    [setState],
  );

  // inputDecimal refuses a second decimal point, so the API never
  // receives an unparseable number built by the user.
  const inputDecimal = useCallback(() => {
    setState((prev) => {
      if (prev.isLoading) return prev;
      if (prev.errorMessage) return { ...INITIAL_STATE, display: '0.' };
      if (prev.waitingForOperand) return { ...prev, display: '0.', waitingForOperand: false };
      if (prev.display.includes('.')) return prev;
      return { ...prev, display: prev.display + '.' };
    });
  }, [setState]);

  const clear = useCallback(() => {
    setState(INITIAL_STATE);
  }, [setState]);

  // backspace removes the last character of the current entry. It only
  // touches `display` — pendingValue/pendingOperator are left alone,
  // so it never needs to re-derive anything about a pending operation.
  const backspace = useCallback(() => {
    setState((prev) => {
      if (prev.isLoading) return prev;
      if (prev.errorMessage) return INITIAL_STATE;
      if (prev.waitingForOperand) return prev;

      const next = prev.display.slice(0, -1);
      return { ...prev, display: next === '' || next === '-' ? '0' : next };
    });
  }, [setState]);

  // chooseOperator handles both starting a new pending operation and,
  // when one is already in flight (the user typed a second number but
  // hasn't pressed "="), evaluating it first so operators can be
  // chained (e.g. 2 + 3 x -> evaluates 2 + 3, then waits to multiply).
  const chooseOperator = useCallback(
    async (operator: BinaryOperatorKey) => {
      const base = stateRef.current.errorMessage ? INITIAL_STATE : stateRef.current;
      const inputValue = Number.parseFloat(base.display);

      if (base.pendingOperator !== null && !base.waitingForOperand) {
        const a = base.pendingValue as number;
        const pendingOperator = base.pendingOperator;
        const b = inputValue;
        setState({ ...base, isLoading: true, errorMessage: null });

        const result = await callForOperator(pendingOperator, a, b);

        setState((prev) => {
          if (!result.ok) {
            return { ...INITIAL_STATE, errorMessage: errorMessageFor(result.error) };
          }
          return {
            ...prev,
            display: String(result.result),
            pendingValue: result.result,
            pendingOperator: operator,
            waitingForOperand: true,
            isLoading: false,
          };
        });
        return;
      }

      setState({
        ...base,
        pendingValue: base.pendingOperator === null ? inputValue : base.pendingValue,
        pendingOperator: operator,
        waitingForOperand: true,
        errorMessage: null,
      });
    },
    [setState],
  );

  const equals = useCallback(async () => {
    const base = stateRef.current.errorMessage ? INITIAL_STATE : stateRef.current;
    if (base.pendingOperator === null || base.pendingValue === null) return;

    const a = base.pendingValue;
    const operator = base.pendingOperator;
    const b = Number.parseFloat(base.display);
    setState({ ...base, isLoading: true });

    const result = await callForOperator(operator, a, b);

    setState(() => {
      if (!result.ok) {
        return { ...INITIAL_STATE, errorMessage: errorMessageFor(result.error) };
      }
      return { ...INITIAL_STATE, display: String(result.result), waitingForOperand: true };
    });
  }, [setState]);

  // sqrt is unary: it applies to the currently displayed value right
  // away, without waiting for "=".
  const sqrt = useCallback(async () => {
    const base = stateRef.current.errorMessage ? INITIAL_STATE : stateRef.current;
    const value = Number.parseFloat(base.display);
    setState({ ...base, isLoading: true, errorMessage: null });

    const result = await apiClient.sqrt(value);

    setState((prev) => {
      if (!result.ok) {
        return { ...INITIAL_STATE, errorMessage: errorMessageFor(result.error) };
      }
      return { ...prev, display: String(result.result), waitingForOperand: true, isLoading: false };
    });
  }, [setState]);

  const state = stateRef.current;

  return {
    display: state.display,
    errorMessage: state.errorMessage,
    isLoading: state.isLoading,
    pendingOperator: state.pendingOperator,
    pendingValue: state.pendingValue,
    inputDigit,
    inputDecimal,
    clear,
    backspace,
    chooseOperator,
    equals,
    sqrt,
  };
}

export type UseCalculatorResult = ReturnType<typeof useCalculator>;
