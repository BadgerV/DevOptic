module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  extends: [
    'eslint:recommended',
    '@typescript-eslint/recommended',
    'eslint-plugin-react-hooks/exhaustive-deps',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh', 'import'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
    'no-restricted-imports': [
      'error',
      {
        patterns: [
          {
            group: ['../../../*'],
            message: 'Avoid deep relative imports. Use absolute imports instead.'
          },
          {
            group: ['**/domains/*/components/**'],
            message: 'Cross-domain component imports not allowed. Use shared components.'
          }
        ]
      }
    ],
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
  },
}
