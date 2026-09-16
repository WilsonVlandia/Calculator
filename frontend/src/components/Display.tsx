import styles from './Display.module.css';

interface DisplayProps {
  value: string;
  pendingLine: string | null;
  errorMessage: string | null;
  isLoading: boolean;
  onBackspace: () => void;
}

// Length thresholds (in already-formatted characters) at which the
// main value switches to a smaller font tier. A formatted 15-digit
// integer with thousands separators ("123,456,789,012,345") is 19
// characters — longer than the worst-case scientific notation string
// ("-1.23457e-308", 13 characters) — so both normal and
// scientific-notation output need this, not just long numbers.
const MEDIUM_VALUE_LENGTH = 7;
const SMALL_VALUE_LENGTH = 12;

function valueSizeClass(text: string): string {
  if (text.length > SMALL_VALUE_LENGTH) return styles.valueTextSmall;
  if (text.length > MEDIUM_VALUE_LENGTH) return styles.valueTextMedium;
  return '';
}

// Display only renders what it is given: the current (already
// formatted) value, the pending-operation line above it, or the error
// message in the value's place when one is present. It decides
// nothing about calculator behavior or number formatting — those are
// composed by Calculator and the model layer respectively; picking a
// smaller font tier for a longer string is a presentation-only concern
// that belongs here. The ⌫ button just forwards a click event up, the
// same as every Keypad button does, so it stays reachable on touch
// devices that have no physical Backspace key.
export function Display({ value, pendingLine, errorMessage, isLoading, onBackspace }: DisplayProps) {
  const isError = errorMessage !== null;

  return (
    <div className={styles.display} data-loading={isLoading}>
      <div className={styles.topRow}>
        <span className={styles.pendingLine}>{!isError && pendingLine ? pendingLine : ''}</span>
        <button
          type="button"
          className={styles.backspaceButton}
          onClick={onBackspace}
          disabled={isLoading}
          aria-label="Backspace"
        >
          ⌫
        </button>
      </div>
      <span className={isError ? styles.errorText : `${styles.valueText} ${valueSizeClass(value)}`}>
        {isError ? errorMessage : value}
      </span>
    </div>
  );
}
