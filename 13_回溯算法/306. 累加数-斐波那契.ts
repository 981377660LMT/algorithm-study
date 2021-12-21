/**
 * @link https://leetcode-cn.com/problems/additive-number/solution/xia-biao-zuo-wei-fen-duan-dian-dfs-by-over-lord/
 * @param {string} num
 * @return {boolean}
 * 除了最开始的两个数以外，字符串中的其他数都等于它之前两个数相加的和。
 * 来判断给定输入是否是累加数。
 * @summary 只要固定前面两个数，后面的数都是确定的。所以我们用两个循环找前面两个数，看是否可以引出后面的数。
 * 你如何处理一个溢出的过大的整数输入?
 * @summary 为了处理大正数相加应该使用两字符串相加的程序
 *
 */
const isAdditiveNumber = function (num: string): boolean {
  if (num.length <= 2) return false

  const addStrings = (num1: string, num2: string): string => {
    let i = num1.length - 1
    let j = num2.length - 1
    let carry = 0
    let res = ''

    while (i >= 0 || j >= 0 || carry) {
      const d1 = i < 0 ? 0 : parseInt(num1[i])
      const d2 = j < 0 ? 0 : parseInt(num2[j])
      const sum = d1 + d2 + carry
      res = `${sum % 10}${res}`
      carry = ~~(sum / 10)
      i--
      j--
    }

    return res
  }

  const bt = (num1: string, num2: string, rest: string): boolean => {
    if (rest === '') return true
    const sum = addStrings(num1, num2)
    if (!rest.startsWith(sum)) {
      return false
    }
    return bt(num2, sum, rest.slice(sum.length))
  }

  // 出发点 bt
  for (let i = 0; i < num.length - 1; i++) {
    for (let j = i + 1; j < num.length; j++) {
      const num1 = num.slice(0, i + 1)
      const num2 = num.slice(i + 1, j + 1)
      const rest = num.slice(j + 1)
      if (num1.length > 1 && num1[0] === '0') return false
      // 剪枝
      if (num2.length > 1 && num2[0] === '0') break
      if (rest.length < num1.length || rest.length < num2.length) break
      if (bt(num1, num2, rest)) return true
    }
  }

  return false
}

console.log(isAdditiveNumber('199100199'))

export default 1
