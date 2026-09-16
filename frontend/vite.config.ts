import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

// Vitest reads its `test` options from this same Vite config, so the
// dev server, the production build, and the test runner all share one
// source of truth.
export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts'],
    globals: false,
  },
});
