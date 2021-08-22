/**
 * @param {string} s 4 <= S.length <= 12.
 * @return {string[]}
 * @description
 * 我们有一些二维坐标，如 "(1, 3)" 或 "(2, 0.5)"，然后我们移除所有逗号，小数点和空格，得到一个字符串S
 * 返回所有可能的原始字符串到一个列表中。
 * @description 拆成两个数字，逗号分隔 如果数字大于2位则可以继续加.号
 * O(n^3) 时间复杂度   O(n^2) 空间复杂度
 * @summary 重点在于不合法分割的判断，整数部分不能有leading zero， 小数部分不能有trailing zero。
 * 单独写小函数分别判断这两个逻辑比较清楚。
 */
const ambiguousCoordinates = function (s: string): string[] {
  const res: string[] = []
  const num = s.slice(1, -1)

  const getValidSplits = (str: string) => {
    const res: string[] = []
    const isValidInteger = (str: string) => str.length === 1 || str[0] !== '0'
    const isValidDecimal = (str: string) => str[str.length - 1] !== '0'

    // 不加.
    if (isValidInteger(str)) res.push(str)

    // 加.
    for (let i = 1; i < str.length; i++) {
      const integer = str.slice(0, i)
      const decimal = str.slice(i)
      if (isValidInteger(integer) && isValidDecimal(decimal)) {
        res.push(`${integer}.${decimal}`)
      }
    }

    return res
  }

  for (let i = 1; i < num.length; i++) {
    const left = getValidSplits(num.slice(0, i))
    const right = getValidSplits(num.slice(i))
    for (const l of left) {
      for (const r of right) {
        res.push('(' + l + ', ' + r + ')')
      }
    }
  }

  return res
}

console.log(ambiguousCoordinates('(123)'))
// 输出: ["(1, 23)", "(12, 3)", "(1.2, 3)", "(1, 2.3)"]
console.log(ambiguousCoordinates('(0123)'))
// ["(0, 1.23)","(0, 12.3)","(0, 123)","(0.1, 2.3)","(0.1, 23)","(0.12, 3)"]
