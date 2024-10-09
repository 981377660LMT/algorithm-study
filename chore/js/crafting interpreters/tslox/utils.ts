/* eslint-disable implicit-arrow-linebreak */

export const isDigit = (c: string): boolean => c >= '0' && c <= '9'

export const isAlpha = (c: string): boolean =>
  (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c === '_'

export const isAlphaNumeric = (c: string): boolean => isAlpha(c) || isDigit(c)

export const capitalize = (s: string): string => s.charAt(0).toUpperCase() + s.slice(1)

// eslint-disable-next-line @typescript-eslint/no-empty-function
export const noop = (): void => {}
