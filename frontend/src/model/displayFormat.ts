// Pure, presentation-agnostic formatting for numbers shown on the
// calculator display. No React here — this is testable in isolation,
// the same way the backend's model has no dependency on net/http.

// MAX_DISPLAY_DIGITS is the largest digit count (integer + decimal
// digits combined, sign/point/separators excluded) for which EVERY
// value of that length is guaranteed to be an exactly representable
// double: 10^15 - 1 = 999,999,999,999,999 is below
// Number.MAX_SAFE_INTEGER (2^53 - 1 = 9,007,199,254,740,991, 16
// digits), so every 15-digit integer is exact. At 16 digits that
// guarantee breaks (10^16 - 1 exceeds MAX_SAFE_INTEGER), so 16 is the
// first length that can lose precision. 15 is therefore the last full
// "safe" digit count, not merely "close to 16".
const MAX_DISPLAY_DIGITS = 15;

// Significant digits shown in the mantissa once a number switches to
// scientific notation (matches the fixed-width scientific mode of a
// physical calculator).
const SCIENTIFIC_NOTATION_SIGNIFICANT_DIGITS = 6;

// formatDisplayNumber adds thousands separators to the integer part of
// a number string, or switches to scientific notation once the number
// is too long to show clearly. It takes the exact string as stored in
// calculator state (not a number) so it can preserve a decimal point
// the user just typed ("1234." while typing) or trailing zeros already
// present ("1000.50") without re-deriving them from a parsed float,
// which would silently drop that information.
export function formatDisplayNumber(rawValue: string): string {
  const isNegative = rawValue.startsWith('-');
  const unsigned = isNegative ? rawValue.slice(1) : rawValue;
  const [integerPart, decimalPart] = unsigned.split('.');

  const digitCount = integerPart.length + (decimalPart?.length ?? 0);
  if (digitCount > MAX_DISPLAY_DIGITS) {
    return Number(rawValue).toExponential(SCIENTIFIC_NOTATION_SIGNIFICANT_DIGITS - 1);
  }

  const groupedIntegerPart = integerPart.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  const sign = isNegative ? '-' : '';

  if (decimalPart === undefined) {
    return sign + groupedIntegerPart;
  }
  return `${sign}${groupedIntegerPart}.${decimalPart}`;
}
