import { describe, expect, it } from 'vitest';
import { formatDisplayNumber } from './displayFormat';

describe('formatDisplayNumber', () => {
  it('adds thousands separators to an integer', () => {
    expect(formatDisplayNumber('1234567')).toBe('1,234,567');
  });

  it('does not add separators to a number under 1000', () => {
    expect(formatDisplayNumber('42')).toBe('42');
  });

  it('formats zero as-is', () => {
    expect(formatDisplayNumber('0')).toBe('0');
  });

  it('keeps a trailing decimal point while the user is still typing', () => {
    expect(formatDisplayNumber('1234.')).toBe('1,234.');
  });

  it('preserves decimal digits, including trailing zeros, exactly as given', () => {
    expect(formatDisplayNumber('1000.50')).toBe('1,000.50');
  });

  it('keeps the negative sign attached to the number, unaffected by grouping', () => {
    expect(formatDisplayNumber('-1234567')).toBe('-1,234,567');
  });

  it('formats a negative decimal correctly', () => {
    expect(formatDisplayNumber('-1000.5')).toBe('-1,000.5');
  });

  describe('the scientific notation threshold', () => {
    it('shows a 15-digit integer normally, right at the threshold', () => {
      expect(formatDisplayNumber('123456789012345')).toBe('123,456,789,012,345');
    });

    it('switches a 16-digit integer to scientific notation, right above the threshold', () => {
      const result = formatDisplayNumber('1234567890123456');
      expect(result).toBe(Number('1234567890123456').toExponential(5));
      expect(result).toMatch(/^1\.23457e\+15$/);
    });

    it('shows a mixed integer+decimal number normally when the combined digit count is exactly 15', () => {
      // 10 integer digits + 5 decimal digits = 15
      expect(formatDisplayNumber('1234567890.12345')).toBe('1,234,567,890.12345');
    });

    it('switches a mixed integer+decimal number to scientific notation when the combined digit count is 16', () => {
      // 10 integer digits + 6 decimal digits = 16
      const result = formatDisplayNumber('1234567890.123456');
      expect(result).toBe(Number('1234567890.123456').toExponential(5));
    });

    it('still groups thousands normally for numbers just below the threshold after the switch exists', () => {
      // Regression guard: adding the scientific-notation branch must not
      // affect ordinary formatting for shorter numbers.
      expect(formatDisplayNumber('12345678901234')).toBe('12,345,678,901,234');
    });

    it('formats a negative number that crosses the threshold', () => {
      const result = formatDisplayNumber('-1234567890123456');
      expect(result).toBe(Number('-1234567890123456').toExponential(5));
      expect(result.startsWith('-')).toBe(true);
    });
  });
});
