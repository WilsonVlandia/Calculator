import { afterEach, describe, expect, it, vi } from 'vitest';
import { apiClient } from './apiClient';

function mockFetchResolving(body: unknown) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      json: () => Promise.resolve(body),
    } as Response),
  );
}

describe('apiClient', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('returns a typed success result', async () => {
    mockFetchResolving({ success: true, operation: 'add', result: 5 });

    const result = await apiClient.add(2, 3);

    expect(result).toEqual({ ok: true, operation: 'add', result: 5 });
  });

  it('returns a typed error result for a domain error', async () => {
    mockFetchResolving({ success: false, error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' } });

    const result = await apiClient.divide(10, 0);

    expect(result).toEqual({
      ok: false,
      error: { code: 'DIVISION_BY_ZERO', message: 'division by zero' },
    });
  });

  it('sends the right payload shape and path for power', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true, operation: 'power', result: 8 }),
    } as Response);
    vi.stubGlobal('fetch', fetchMock);

    await apiClient.power(2, 3);

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/power'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ base: 2, exponent: 3 }),
      }),
    );
  });

  it('sends the right payload shape for percentage', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ success: true, operation: 'percentage', result: 30 }),
    } as Response);
    vi.stubGlobal('fetch', fetchMock);

    await apiClient.percentage(200, 15);

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/percentage'),
      expect.objectContaining({
        body: JSON.stringify({ value: 200, percentage: 15 }),
      }),
    );
  });

  it('returns NETWORK_ERROR when fetch itself rejects', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('network down')));

    const result = await apiClient.add(1, 1);

    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error.code).toBe('NETWORK_ERROR');
    }
  });

  it('returns NETWORK_ERROR when the response body is not valid JSON', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        json: () => Promise.reject(new Error('unexpected token')),
      } as unknown as Response),
    );

    const result = await apiClient.sqrt(4);

    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error.code).toBe('NETWORK_ERROR');
    }
  });
});
