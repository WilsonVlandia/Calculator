import { OPERATOR_SYMBOLS, type BinaryOperatorKey } from '../hooks/useCalculator';
import { formatDisplayNumber } from '../model/displayFormat';
import styles from './Calculator.module.css';
import { Display } from './Display';
import { Keypad } from './Keypad';

interface CalculatorProps {
  display: string;
  pendingOperator: BinaryOperatorKey | null;
  pendingValue: number | null;
  errorMessage: string | null;
  isLoading: boolean;
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperator: (operator: BinaryOperatorKey) => void;
  onSqrt: () => void;
  onEquals: () => void;
  onClear: () => void;
  onBackspace: () => void;
}

// Calculator assembles Display and Keypad from its props and is the
// one place that turns raw state into display-ready strings (thousands
// separators, scientific notation, the "1,000 +" pending-operation
// line) via the pure formatDisplayNumber helper. It still holds no
// state and calls no API itself.
export function Calculator({
  display,
  pendingOperator,
  pendingValue,
  errorMessage,
  isLoading,
  onDigit,
  onDecimal,
  onOperator,
  onSqrt,
  onEquals,
  onClear,
  onBackspace,
}: CalculatorProps) {
  const formattedValue = formatDisplayNumber(display);
  const pendingLine =
    pendingOperator !== null && pendingValue !== null
      ? `${formatDisplayNumber(String(pendingValue))} ${OPERATOR_SYMBOLS[pendingOperator]}`
      : null;

  return (
    <div className={styles.calculator}>
      <Display
        value={formattedValue}
        pendingLine={pendingLine}
        errorMessage={errorMessage}
        isLoading={isLoading}
        onBackspace={onBackspace}
      />
      <Keypad
        onDigit={onDigit}
        onDecimal={onDecimal}
        onOperator={onOperator}
        onSqrt={onSqrt}
        onEquals={onEquals}
        onClear={onClear}
        disabled={isLoading}
      />
    </div>
  );
}
