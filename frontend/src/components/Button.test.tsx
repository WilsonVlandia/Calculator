import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { Button } from './Button';

describe('Button', () => {
  it('renders its label', () => {
    render(<Button label="7" variant="digit" onClick={() => {}} />);
    expect(screen.getByRole('button', { name: '7' })).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button label="=" variant="equals" onClick={handleClick} />);

    await user.click(screen.getByRole('button', { name: '=' }));

    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('does not call onClick when disabled', async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();
    render(<Button label="+" variant="operator" onClick={handleClick} disabled />);

    await user.click(screen.getByRole('button', { name: '+' }));

    expect(handleClick).not.toHaveBeenCalled();
  });
});
