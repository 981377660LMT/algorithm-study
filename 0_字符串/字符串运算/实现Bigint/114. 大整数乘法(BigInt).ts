/**
 * @param {string} a
 * @param {string} b
 * @return {string}
 */
function multiply(num1: string, num2: string): string {
  if (num1 === '0' || num2 === '0') return '0'
  const isNum1Negative = num1[0] === '-'
  const isNum2Negative = num2[0] === '-'
  const isSameOperator = !(Number(isNum1Negative) ^ Number(isNum2Negative))
  num1 = num1.replace(/^[+|-]/, '')
  num2 = num2.replace(/^[+|-]/, '')

  const arr1 = num1.split('').map(Number)
  const arr2 = num2.split('').map(Number)
  // console.log(arr1, arr2)
  const res = Array(arr1.length + arr2.length - 1).fill(0)
  for (let i = 0; i < arr1.length; i++) {
    for (let j = 0; j < arr2.length; j++) {
      res[i + j] += arr1[i] * arr2[j]
    }
  }

  handleCarry(res)
  return (isSameOperator ? '' : '-') + res.join('')

  /**
   * @param arr   4   13  28  27  18
   * @description 整理： c[i + 1] += c[i] / 10, c[i] %= 10, 从低位开始。
   */
  function handleCarry(arr: number[]) {
    let carry = 0
    for (let i = arr.length - 1; i >= 1; i--) {
      carry = ~~(arr[i] / 10)
      arr[i - 1] += carry
      arr[i] %= 10
    }

    // 单独处理第一位
    if (arr[0] >= 10) {
      carry = ~~(arr[0] / 10)
      arr[0] %= 10
      arr.unshift(carry)
    }
  }
}

// console.log(multiply('-99999999999999999', '-99999999999999999'))
console.log(multiply('1123456787654323456789', '1234567887654323456'))

export {}
