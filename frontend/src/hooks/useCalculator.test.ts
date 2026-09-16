import { act, renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { apiClient } from '../model/apiClient';
import type { ApiResult } from '../model/types';
import { useCalculator } from './useCalculator';

vi.mock('../model/apiClient', () => ({
  apiClient: {
    add: vi.fn(),
    subtract: vi.fn(),
    multiply: vi.fn(),
    divide: vi.fn(),
    power: vi.fn(),
    sqrt: vi.fn(),
    percentage: vi.fn(),
  },
}));

const mockedApiClient = vi.mocked(apiClient);

beforeEach(() => {
  vi.clearAllMocks();
});

describe('useCalculator', () => {
  it('builds a number by pressing digits', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('4');
      result.current.inputDigit('2');
    });

    expect(result.current.display).toBe('42');
  });

  it('does not allow a second decimal point', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.inputDecimal();
      result.current.inputDigit('5');
      result.current.inputDecimal();
    });

    expect(result.current.display).toBe('1.5');
  });

  it('performs a full add operation via the API client', async () => {
    mockedApiClient.add.mockResolvedValue({ ok: true, operation: 'add', result: 5 });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.chooseOperator('add');
      result.current.inputDigit('3');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(mockedApiClient.add).toHaveBeenCalledWith(2, 3);
    expect(result.current.display).toBe('5');
  });

  it('maps a domain error to its dedicated message', async () => {
    mockedApiClient.divide.mockResolvedValue({
      ok: false,
      error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' },
    });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.chooseOperator('divide');
      result.current.inputDigit('0');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(result.current.errorMessage).toBe('Cannot divide by zero');
  });

  it('maps a NEGATIVE_SQRT_INPUT error to its own dedicated message', async () => {
    mockedApiClient.sqrt.mockResolvedValue({
      ok: false,
      error: { code: 'NEGATIVE_SQRT_INPUT', message: 'cannot compute the square root of a negative number' },
    });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('9');
    });

    await act(async () => {
      await result.current.sqrt();
    });

    expect(result.current.errorMessage).toBe('Cannot take the square root of a negative number');
  });

  it('clears the error and starts fresh on the next digit press', async () => {
    mockedApiClient.divide.mockResolvedValue({
      ok: false,
      error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' },
    });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.chooseOperator('divide');
      result.current.inputDigit('0');
    });
    await act(async () => {
      await result.current.equals();
    });
    expect(result.current.errorMessage).not.toBeNull();

    act(() => {
      result.current.inputDigit('9');
    });

    expect(result.current.errorMessage).toBeNull();
    expect(result.current.display).toBe('9');
  });

  it('computes sqrt immediately without needing equals', async () => {
    mockedApiClient.sqrt.mockResolvedValue({ ok: true, operation: 'sqrt', result: 4 });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('1');
      result.current.inputDigit('6');
    });

    await act(async () => {
      await result.current.sqrt();
    });

    expect(mockedApiClient.sqrt).toHaveBeenCalledWith(16);
    expect(result.current.display).toBe('4');
  });

  it('computes percentage as percentage% of value', async () => {
    mockedApiClient.percentage.mockResolvedValue({ ok: true, operation: 'percentage', result: 30 });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.inputDigit('0');
      result.current.inputDigit('0');
      result.current.chooseOperator('percentage');
      result.current.inputDigit('1');
      result.current.inputDigit('5');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(mockedApiClient.percentage).toHaveBeenCalledWith(200, 15);
    expect(result.current.display).toBe('30');
  });

  it('chains operators by evaluating the pending operation first', async () => {
    mockedApiClient.add.mockResolvedValue({ ok: true, operation: 'add', result: 5 });
    mockedApiClient.multiply.mockResolvedValue({ ok: true, operation: 'multiply', result: 10 });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.chooseOperator('add');
      result.current.inputDigit('3');
    });

    await act(async () => {
      await result.current.chooseOperator('multiply');
    });

    expect(mockedApiClient.add).toHaveBeenCalledWith(2, 3);
    expect(result.current.display).toBe('5');

    act(() => {
      result.current.inputDigit('2');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(mockedApiClient.multiply).toHaveBeenCalledWith(5, 2);
    expect(result.current.display).toBe('10');
  });

  it('does nothing when equals is pressed with no pending operator', async () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('7');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(mockedApiClient.add).not.toHaveBeenCalled();
    expect(result.current.display).toBe('7');
  });

  it('clears all state when AC is pressed', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('9');
      result.current.clear();
    });

    expect(result.current.display).toBe('0');
    expect(result.current.errorMessage).toBeNull();
  });

  it('sets isLoading while an operation is in flight', async () => {
    let resolvePromise: (value: ApiResult) => void = () => {};
    mockedApiClient.add.mockImplementation(
      () =>
        new Promise<ApiResult>((resolve) => {
          resolvePromise = resolve;
        }),
    );
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.chooseOperator('add');
      result.current.inputDigit('3');
    });

    let equalsPromise!: Promise<void>;
    act(() => {
      equalsPromise = result.current.equals();
    });

    expect(result.current.isLoading).toBe(true);

    await act(async () => {
      resolvePromise({ ok: true, operation: 'add', result: 5 });
      await equalsPromise;
    });

    expect(result.current.isLoading).toBe(false);
  });

  it('exposes the pending operator and value while an operation is set up', () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.inputDigit('0');
      result.current.inputDigit('0');
      result.current.chooseOperator('add');
    });

    expect(result.current.pendingOperator).toBe('add');
    expect(result.current.pendingValue).toBe(200);
  });

  it('clears the pending operator and value once equals resolves', async () => {
    mockedApiClient.add.mockResolvedValue({ ok: true, operation: 'add', result: 5 });
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit('2');
      result.current.chooseOperator('add');
      result.current.inputDigit('3');
    });

    await act(async () => {
      await result.current.equals();
    });

    expect(result.current.pendingOperator).toBeNull();
    expect(result.current.pendingValue).toBeNull();
  });

  describe('backspace', () => {
    it('removes the last character of the current entry', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('1');
        result.current.inputDigit('2');
        result.current.inputDigit('3');
        result.current.backspace();
      });

      expect(result.current.display).toBe('12');
    });

    it('resets to "0" once the last digit is erased', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('7');
        result.current.backspace();
      });

      expect(result.current.display).toBe('0');
    });

    it('clears an error and starts fresh instead of editing it', async () => {
      mockedApiClient.divide.mockResolvedValue({
        ok: false,
        error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' },
      });
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('1');
        result.current.chooseOperator('divide');
        result.current.inputDigit('0');
      });
      await act(async () => {
        await result.current.equals();
      });
      expect(result.current.errorMessage).not.toBeNull();

      act(() => {
        result.current.backspace();
      });

      expect(result.current.errorMessage).toBeNull();
      expect(result.current.display).toBe('0');
    });

    it('does nothing while waiting for the second operand', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('5');
        result.current.chooseOperator('add');
        result.current.backspace();
      });

      expect(result.current.display).toBe('5');
    });

    it('does not touch pendingValue or pendingOperator', () => {
      const { result } = renderHook(() => useCalculator());

      act(() => {
        result.current.inputDigit('2');
        result.current.chooseOperator('add');
        result.current.inputDigit('3');
        result.current.backspace();
      });

      expect(result.current.display).toBe('0');
      expect(result.current.pendingOperator).toBe('add');
      expect(result.current.pendingValue).toBe(2);
    });
  });
});
