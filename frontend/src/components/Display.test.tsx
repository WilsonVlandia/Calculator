import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { Display } from './Display';

describe('Display', () => {
  it('shows the numeric value when there is no error', () => {
    render(<Display value="42" pendingLine={null} errorMessage={null} isLoading={false} onBackspace={() => {}} />);

    expect(screen.getByText('42')).toBeInTheDocument();
  });

  it('shows the error message instead of the value when present', () => {
    render(
      <Display
        value="0"
        pendingLine={null}
        errorMessage="Cannot divide by zero"
        isLoading={false}
        onBackspace={() => {}}
      />,
    );

    expect(screen.getByText('Cannot divide by zero')).toBeInTheDocument();
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });

  it('shows the pending operation line when given one', () => {
    render(
      <Display value="3" pendingLine="1,000 +" errorMessage={null} isLoading={false} onBackspace={() => {}} />,
    );

    expect(screen.getByText('1,000 +')).toBeInTheDocument();
  });

  it('hides the pending line while an error is shown', () => {
    render(
      <Display
        value="0"
        pendingLine="1,000 +"
        errorMessage="Cannot divide by zero"
        isLoading={false}
        onBackspace={() => {}}
      />,
    );

    expect(screen.queryByText('1,000 +')).not.toBeInTheDocument();
  });

  it('calls onBackspace when the backspace button is clicked', async () => {
    const user = userEvent.setup();
    const handleBackspace = vi.fn();
    render(
      <Display value="12" pendingLine={null} errorMessage={null} isLoading={false} onBackspace={handleBackspace} />,
    );

    await user.click(screen.getByRole('button', { name: 'Backspace' }));

    expect(handleBackspace).toHaveBeenCalledTimes(1);
  });

  it('disables the backspace button while loading', () => {
    render(<Display value="12" pendingLine={null} errorMessage={null} isLoading onBackspace={() => {}} />);

    expect(screen.getByRole('button', { name: 'Backspace' })).toBeDisabled();
  });

  describe('font size tiers for long values', () => {
    it('uses the default (large) size for a short value', () => {
      render(<Display value="1,234" pendingLine={null} errorMessage={null} isLoading={false} onBackspace={() => {}} />);

      const valueEl = screen.getByText('1,234');
      expect(valueEl.className).not.toMatch(/valueTextMedium|valueTextSmall/);
    });

    it('uses the medium size for a value past the medium threshold', () => {
      // "12,345,678" is 10 characters, past the 7-character threshold.
      render(
        <Display value="12,345,678" pendingLine={null} errorMessage={null} isLoading={false} onBackspace={() => {}} />,
      );

      const valueEl = screen.getByText('12,345,678');
      expect(valueEl.className).toMatch(/valueTextMedium/);
    });

    it('uses the small size for a value past the small threshold', () => {
      // A 15-digit number with separators, 19 characters.
      const longValue = '123,456,789,012,345';
      render(<Display value={longValue} pendingLine={null} errorMessage={null} isLoading={false} onBackspace={() => {}} />);

      const valueEl = screen.getByText(longValue);
      expect(valueEl.className).toMatch(/valueTextSmall/);
    });

    it('uses the small size for a scientific notation value', () => {
      const scientificValue = '-1.23457e-308';
      render(
        <Display
          value={scientificValue}
          pendingLine={null}
          errorMessage={null}
          isLoading={false}
          onBackspace={() => {}}
        />,
      );

      const valueEl = screen.getByText(scientificValue);
      expect(valueEl.className).toMatch(/valueTextSmall/);
    });
  });
});
