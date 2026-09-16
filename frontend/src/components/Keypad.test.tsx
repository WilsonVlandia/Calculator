import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import type { ComponentProps } from 'react';
import { describe, expect, it, vi } from 'vitest';
import { Keypad } from './Keypad';

function renderKeypad(overrides: Partial<ComponentProps<typeof Keypad>> = {}) {
  const props = {
    onDigit: vi.fn(),
    onDecimal: vi.fn(),
    onOperator: vi.fn(),
    onSqrt: vi.fn(),
    onEquals: vi.fn(),
    onClear: vi.fn(),
    disabled: false,
    ...overrides,
  };
  render(<Keypad {...props} />);
  return props;
}

describe('Keypad', () => {
  it('renders all 20 keys', () => {
    renderKeypad();
    expect(screen.getAllByRole('button')).toHaveLength(20);
  });

  it('forwards a digit press with the pressed digit', async () => {
    const user = userEvent.setup();
    const props = renderKeypad();

    await user.click(screen.getByRole('button', { name: '7' }));

    expect(props.onDigit).toHaveBeenCalledWith('7');
  });

  it('forwards an operator press with its operator key', async () => {
    const user = userEvent.setup();
    const props = renderKeypad();

    await user.click(screen.getByRole('button', { name: '÷' }));

    expect(props.onOperator).toHaveBeenCalledWith('divide');
  });

  it('forwards the power operator', async () => {
    const user = userEvent.setup();
    const props = renderKeypad();

    await user.click(screen.getByRole('button', { name: 'xʸ' }));

    expect(props.onOperator).toHaveBeenCalledWith('power');
  });

  it('forwards sqrt, equals, decimal and clear presses', async () => {
    const user = userEvent.setup();
    const props = renderKeypad();

    await user.click(screen.getByRole('button', { name: '√' }));
    await user.click(screen.getByRole('button', { name: '=' }));
    await user.click(screen.getByRole('button', { name: '.' }));
    await user.click(screen.getByRole('button', { name: 'AC' }));

    expect(props.onSqrt).toHaveBeenCalledTimes(1);
    expect(props.onEquals).toHaveBeenCalledTimes(1);
    expect(props.onDecimal).toHaveBeenCalledTimes(1);
    expect(props.onClear).toHaveBeenCalledTimes(1);
  });

  it('disables every key when disabled is true', () => {
    renderKeypad({ disabled: true });

    for (const button of screen.getAllByRole('button')) {
      expect(button).toBeDisabled();
    }
  });
});
