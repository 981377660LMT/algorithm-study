/**
 * @type {import('ts-jest/dist/types').InitialOptionsTsJest}
 */
module.exports = {
  // preset:'',
  // roots: ['<rootDir>/前端笔记'],
  testMatch: ['**/__tests__/*.[jt]s?(x)', '**/*.test.[jt]s?(x)'],
  transform: {
    '^.+\\.tsx?$': 'ts-jest',
  },
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  transformIgnorePatterns: ['node_modules'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  collectCoverage: true,
  coveragePathIgnorePatterns: ['node_modules'],
  testEnvironment: 'node',
  modulePathIgnorePatterns: ['node_modules'],
}
