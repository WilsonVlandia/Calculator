import { cleanup } from '@testing-library/react';
import { afterEach } from 'vitest';
import '@testing-library/jest-dom/vitest';

// globals: false in vite.config.ts means Testing Library's automatic
// cleanup (which relies on a global afterEach) never registers itself,
// so it is wired up explicitly here.
afterEach(() => {
  cleanup();
});
