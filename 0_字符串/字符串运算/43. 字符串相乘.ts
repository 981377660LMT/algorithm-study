/**
 * @link https://leetcode-cn.com/problems/multiply-strings/solution/si-lu-qing-xi-by-lllllliuji-2/
 * @param {string} num1  num1 和 num2 的长度小于110
 * @param {string} num2
 * @return {string}
 * @summary
 * 1.把两个数用数组 a, b 来存储，并且反转（从个位开始乘）
 * 2.对于 a 的第 i 位 和 b 的第 j 位相乘的结果存储在 c[i + j] 上，即 c[i + j] += a[i] * b[j];
 * 3. 最后，从 c 的低位向高位整理，c[i + 1] = c[i] / 10, c[i] %= 10
 */
const multiply = (num1: string, num2: string): string => {
  if (num1 === '0' || num2 === '0') return '0'

  /**
   * @param arr   4   13  28  27  18
   * @description 整理： c[i + 1] += c[i] / 10, c[i] %= 10, 从低位开始。
   */
  const handleCarry = (arr: number[]): void => {
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

  const arr1 = num1.split('').map(Number)
  const arr2 = num2.split('').map(Number)
  const res: number[] = Array(arr1.length + arr2.length - 1).fill(0)
  for (let i = 0; i < arr1.length; i++) {
    for (let j = 0; j < arr2.length; j++) {
      res[i + j] += arr1[i] * arr2[j]
    }
  }

  handleCarry(res)
  return res.join('')
}

console.log(323 * 456)
console.log(multiply('123', '456'))
console.log(multiply('323', '456'))

export {}
