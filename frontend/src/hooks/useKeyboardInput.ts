// useKeyboardInput wires the physical keyboard to the exact same
// handlers the Keypad's buttons already call — it never reimplements
// calculator behavior, only translates a key press into a call to one
// of useCalculator's own methods.
import { useEffect } from 'react';
import type { BinaryOperatorKey } from './useCalculator';

interface KeyboardHandlers {
  inputDigit: (digit: string) => void;
  inputDecimal: () => void;
  chooseOperator: (operator: BinaryOperatorKey) => void;
  equals: () => void;
  clear: () => void;
  backspace: () => void;
}

// xʸ (power) and √ (sqrt) are deliberately not mapped: neither has a
// single, unambiguous standard keyboard key, so they stay click-only.
export function useKeyboardInput(handlers: KeyboardHandlers) {
  const { inputDigit, inputDecimal, chooseOperator, equals, clear, backspace } = handlers;

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      const { key } = event;

      if (/^[0-9]$/.test(key)) {
        event.preventDefault();
        inputDigit(key);
        return;
      }

      switch (key) {
        case '.':
          event.preventDefault();
          inputDecimal();
          return;
        case '+':
          event.preventDefault();
          chooseOperator('add');
          return;
        case '-':
          event.preventDefault();
          chooseOperator('subtract');
          return;
        case '*':
        case 'x':
        case 'X':
          event.preventDefault();
          chooseOperator('multiply');
          return;
        case '/':
          event.preventDefault();
          chooseOperator('divide');
          return;
        case '%':
          event.preventDefault();
          chooseOperator('percentage');
          return;
        case 'Enter':
        case '=':
          event.preventDefault();
          equals();
          return;
        case 'Backspace':
          event.preventDefault();
          backspace();
          return;
        case 'Escape':
          event.preventDefault();
          clear();
          return;
        default:
          return;
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [inputDigit, inputDecimal, chooseOperator, equals, clear, backspace]);
}
