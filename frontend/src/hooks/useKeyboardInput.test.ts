import { renderHook } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useKeyboardInput } from './useKeyboardInput';

function pressKey(key: string) {
  window.dispatchEvent(new KeyboardEvent('keydown', { key, cancelable: true }));
}

function renderHandlers() {
  const handlers = {
    inputDigit: vi.fn(),
    inputDecimal: vi.fn(),
    chooseOperator: vi.fn(),
    equals: vi.fn(),
    clear: vi.fn(),
    backspace: vi.fn(),
  };
  renderHook(() => useKeyboardInput(handlers));
  return handlers;
}

describe('useKeyboardInput', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('forwards a digit key to inputDigit', () => {
    const handlers = renderHandlers();
    pressKey('7');
    expect(handlers.inputDigit).toHaveBeenCalledWith('7');
  });

  it('forwards "." to inputDecimal', () => {
    const handlers = renderHandlers();
    pressKey('.');
    expect(handlers.inputDecimal).toHaveBeenCalledTimes(1);
  });

  it('maps +, -, *, x and / to their operators', () => {
    const handlers = renderHandlers();
    pressKey('+');
    pressKey('-');
    pressKey('*');
    pressKey('x');
    pressKey('/');
    expect(handlers.chooseOperator.mock.calls).toEqual([
      ['add'],
      ['subtract'],
      ['multiply'],
      ['multiply'],
      ['divide'],
    ]);
  });

  it('maps % to the percentage operator', () => {
    const handlers = renderHandlers();
    pressKey('%');
    expect(handlers.chooseOperator).toHaveBeenCalledWith('percentage');
  });

  it('treats Enter the same as clicking =', () => {
    const handlers = renderHandlers();
    pressKey('Enter');
    expect(handlers.equals).toHaveBeenCalledTimes(1);
  });

  it('treats "=" the same as Enter', () => {
    const handlers = renderHandlers();
    pressKey('=');
    expect(handlers.equals).toHaveBeenCalledTimes(1);
  });

  it('clears on Escape', () => {
    const handlers = renderHandlers();
    pressKey('Escape');
    expect(handlers.clear).toHaveBeenCalledTimes(1);
  });

  it('erases the last character on Backspace', () => {
    const handlers = renderHandlers();
    pressKey('Backspace');
    expect(handlers.backspace).toHaveBeenCalledTimes(1);
  });

  it('does not map power or sqrt to any key', () => {
    const handlers = renderHandlers();
    pressKey('p');
    pressKey('s');
    pressKey('^');
    expect(handlers.chooseOperator).not.toHaveBeenCalled();
  });

  it('calls preventDefault only for keys it handles', () => {
    renderHandlers();

    const handledEvent = new KeyboardEvent('keydown', { key: '/', cancelable: true });
    window.dispatchEvent(handledEvent);
    expect(handledEvent.defaultPrevented).toBe(true);

    const unhandledEvent = new KeyboardEvent('keydown', { key: 'F5', cancelable: true });
    window.dispatchEvent(unhandledEvent);
    expect(unhandledEvent.defaultPrevented).toBe(false);
  });

  it('removes the listener on unmount', () => {
    const handlers = {
      inputDigit: vi.fn(),
      inputDecimal: vi.fn(),
      chooseOperator: vi.fn(),
      equals: vi.fn(),
      clear: vi.fn(),
      backspace: vi.fn(),
    };
    const { unmount } = renderHook(() => useKeyboardInput(handlers));

    unmount();
    pressKey('5');

    expect(handlers.inputDigit).not.toHaveBeenCalled();
  });
});
