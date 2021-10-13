import { add as noSignAdd } from './62. 实现大整数加法(BigInt)'

/**
 * @param {string} a
 * @param {string} b
 * @return {string} 注意结果需要四舍五入，去掉小数部。
 */
function divide(num1: string, num2: string): string {
  if (num2 === '0') throw new Error('invalid divider')

  const isNum1Negative = num1[0] === '-'
  const isNum2Negative = num2[0] === '-'
  const isSameOperator = !(Number(isNum1Negative) ^ Number(isNum2Negative))
  num1 = num1.replace(/^[+|-]/, '')
  num2 = num2.replace(/^[+|-]/, '')

  let divider = num2
  let res = 1
  while (compareAbsolute(divider, num1) <= 0) {
    divider = noSignAdd(divider, num2)
    res++
  }
  res--
  return isSameOperator ? `${res}` : `-${res}`

  function compareAbsolute(str1: string, str2: string): -1 | 0 | 1 {
    if (str1.length !== str2.length) return str1.length > str2.length ? 1 : -1

    for (let i = 0; i < str1.length; i++) {
      if (str1[i] !== str2[i]) return str1[i] > str2[i] ? 1 : -1
    }

    return 0
  }
}

console.log(divide('1123456787654323456789', '1234567887654323456'))
console.log(divide('-1123456787654323456789', '1234567887654323456'))
console.log(divide('5', '2'))
console.log(divide('-3', '2'))

export {}
