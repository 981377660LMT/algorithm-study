import { add as noSignAdd } from './62. 实现大整数加法(BigInt)'
import { subtract as noSignSubtract } from './75. 实现大整数减法(BigInt)'

/**
 * @param {string} num1
 * @param {string} num2
 * @return {string}
 */
function add(num1: string, num2: string): string {
  const isNum1Negative = num1[0] === '-'
  const isNum2Negative = num2[0] === '-'
  const isSameOperator = !(Number(isNum1Negative) ^ Number(isNum2Negative))
  num1 = num1.replace(/^[+|-]/, '')
  num2 = num2.replace(/^[+|-]/, '')

  let sign: string
  let res: string
  if (isSameOperator) {
    sign = isNum1Negative ? '-' : ''
    res = noSignAdd(num1, num2)
  } else if (compareAbsolute(num1, num2) >= 0) {
    sign = isNum1Negative ? '-' : ''
    res = noSignSubtract(num1, num2)
  } else {
    sign = isNum2Negative ? '-' : ''
    res = noSignSubtract(num2, num1)
  }

  return `${sign}${res}`

  function compareAbsolute(str1: string, str2: string): -1 | 0 | 1 {
    if (str1.length !== str2.length) return str1.length > str2.length ? 1 : -1

    for (let i = 0; i < str1.length; i++) {
      if (str1[i] !== str2[i]) return str1[i] > str2[i] ? 1 : -1
    }

    return 0
  }
}

// '-1000000000000000000'
console.log(add('-999999999999999999', '-1'))

// '-999999999999999998'
console.log(add('-999999999999999999', '+1'))

export {}
