/**
 * @param {string} a
 * @param {string} b
 * @return {string} 注意结果需要四舍五入，去掉小数部。
 */
function divide(num1: string, num2: string): string {
  const isNum1Negative = num1[0] === '-'
  const isNum2Negative = num2[0] === '-'
  const isSameOperator = !(Number(isNum1Negative) ^ Number(isNum2Negative))
  num1 = num1.replace(/^[+|-]/, '')
  num2 = num2.replace(/^[+|-]/, '')
}

console.log(divide('1123456787654323456789', '1234567887654323456'))
console.log(divide('-1123456787654323456789', '1234567887654323456'))
console.log(divide('5', '2'))
console.log(divide('-3', '2'))

export {}
