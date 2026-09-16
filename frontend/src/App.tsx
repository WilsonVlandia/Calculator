import styles from './App.module.css';
import { Calculator } from './components/Calculator';
import { useCalculator } from './hooks/useCalculator';
import { useKeyboardInput } from './hooks/useKeyboardInput';

// App is the only place that connects the controller (useCalculator)
// to the view (Calculator). Calculator itself never touches the hook.
// It's also the one place the physical keyboard is wired up, reusing
// the exact same handlers passed to Calculator/Keypad below — no
// calculator logic is duplicated for the keyboard path.
function App() {
  const calculator = useCalculator();

  useKeyboardInput({
    inputDigit: calculator.inputDigit,
    inputDecimal: calculator.inputDecimal,
    chooseOperator: calculator.chooseOperator,
    equals: calculator.equals,
    clear: calculator.clear,
    backspace: calculator.backspace,
  });

  return (
    <div className={styles.app}>
      <Calculator
        display={calculator.display}
        pendingOperator={calculator.pendingOperator}
        pendingValue={calculator.pendingValue}
        errorMessage={calculator.errorMessage}
        isLoading={calculator.isLoading}
        onDigit={calculator.inputDigit}
        onDecimal={calculator.inputDecimal}
        onOperator={calculator.chooseOperator}
        onSqrt={calculator.sqrt}
        onEquals={calculator.equals}
        onClear={calculator.clear}
        onBackspace={calculator.backspace}
      />
    </div>
  );
}

export default App;
