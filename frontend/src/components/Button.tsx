import styles from './Button.module.css';

export type ButtonVariant = 'digit' | 'operator' | 'function' | 'equals';

interface ButtonProps {
  label: string;
  onClick: () => void;
  variant: ButtonVariant;
  disabled?: boolean;
}

// Button is purely presentational: it renders a label with a visual
// variant and forwards clicks. It holds no calculator logic of its own.
export function Button({ label, onClick, variant, disabled = false }: ButtonProps) {
  return (
    <button
      type="button"
      className={`${styles.button} ${styles[variant]}`}
      onClick={onClick}
      disabled={disabled}
      aria-label={label}
    >
      {label}
    </button>
  );
}
