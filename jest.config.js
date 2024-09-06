/**
 * @type {import('ts-jest/dist/types').InitialOptionsTsJest}
 */
module.exports = {
  // preset:'',
  // roots: ['<rootDir>/前端笔记'],
  //
  // 默认情况下，testMatch 的值是：
  // Glob 模式中的 ** 表示“任意目录”，* 表示“任意数量的字符”，?(x) 表示“x 是可选的”
  // [
  //   "**/__tests__/**/*.[jt]s?(x)",
  //   "**/?(*.)+(spec|test).[tj]s?(x)"
  // ]
  testMatch: ['**/__tests__/*.[t]s?(x)', '**/*.test.[t]s?(x)', '**/*.spec.[t]s?(x)'],
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1'
  },
  transformIgnorePatterns: ['node_modules'],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  collectCoverage: true,
  coveragePathIgnorePatterns: ['node_modules'],
  testEnvironment: 'node',
  modulePathIgnorePatterns: ['node_modules', 'immutable-js']
}
