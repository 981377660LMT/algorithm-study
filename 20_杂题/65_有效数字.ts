/* eslint-disable @typescript-eslint/ban-ts-comment */
// 有效数字

/**
 * @param {string} s
 * @return {boolean}
 */
function isNumber(str: string) {
  if (str.includes('Infinity')) return false
  // @ts-ignore
  // eslint-disable-next-line no-restricted-globals
  return str !== '' && !isNaN(str)
}

export {}
