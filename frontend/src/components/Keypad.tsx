import { OPERATOR_SYMBOLS, type BinaryOperatorKey } from '../hooks/useCalculator';
import { Button, type ButtonVariant } from './Button';
import styles from './Keypad.module.css';

interface KeypadProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperator: (operator: BinaryOperatorKey) => void;
  onSqrt: () => void;
  onEquals: () => void;
  onClear: () => void;
  disabled: boolean;
}

interface KeyDefinition {
  label: string;
  variant: ButtonVariant;
  onPress: () => void;
}

// Keypad is purely presentational: it only knows the 20-key layout and
// which callback each key forwards to. All calculator logic lives in
// useCalculator.
export function Keypad({ onDigit, onDecimal, onOperator, onSqrt, onEquals, onClear, disabled }: KeypadProps) {
  const keys: KeyDefinition[] = [
    { label: 'AC', variant: 'function', onPress: onClear },
    { label: '√', variant: 'function', onPress: onSqrt },
    { label: OPERATOR_SYMBOLS.percentage, variant: 'operator', onPress: () => onOperator('percentage') },
    { label: OPERATOR_SYMBOLS.divide, variant: 'operator', onPress: () => onOperator('divide') },

    { label: '7', variant: 'digit', onPress: () => onDigit('7') },
    { label: '8', variant: 'digit', onPress: () => onDigit('8') },
    { label: '9', variant: 'digit', onPress: () => onDigit('9') },
    { label: OPERATOR_SYMBOLS.multiply, variant: 'operator', onPress: () => onOperator('multiply') },

    { label: '4', variant: 'digit', onPress: () => onDigit('4') },
    { label: '5', variant: 'digit', onPress: () => onDigit('5') },
    { label: '6', variant: 'digit', onPress: () => onDigit('6') },
    { label: OPERATOR_SYMBOLS.subtract, variant: 'operator', onPress: () => onOperator('subtract') },

    { label: '1', variant: 'digit', onPress: () => onDigit('1') },
    { label: '2', variant: 'digit', onPress: () => onDigit('2') },
    { label: '3', variant: 'digit', onPress: () => onDigit('3') },
    { label: OPERATOR_SYMBOLS.add, variant: 'operator', onPress: () => onOperator('add') },

    { label: OPERATOR_SYMBOLS.power, variant: 'operator', onPress: () => onOperator('power') },
    { label: '0', variant: 'digit', onPress: () => onDigit('0') },
    { label: '.', variant: 'digit', onPress: onDecimal },
    { label: '=', variant: 'equals', onPress: onEquals },
  ];

  return (
    <div className={styles.keypad}>
      {keys.map((key) => (
        <Button key={key.label} label={key.label} variant={key.variant} onClick={key.onPress} disabled={disabled} />
      ))}
    </div>
  );
}
