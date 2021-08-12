/**
 * @param {string} s 4 <= S.length <= 12.
 * @return {string[]}
 * @description
 * 我们有一些二维坐标，如 "(1, 3)" 或 "(2, 0.5)"，然后我们移除所有逗号，小数点和空格，得到一个字符串S
 * 返回所有可能的原始字符串到一个列表中。
 */
const ambiguousCoordinates = function (s: string): string[] {
  const res: string[] = []
  const num = s.slice(1, -1)

  const genNum = (str: string) => {
    const res: string[] = []
    for (let i = 1; i < str.length; i++) {
      res.push(str.slice(0, i) + '.' + str.slice(i))
    }
    res.push(str)
    return res
  }

  for (let i = 1; i < s.length; i++) {
    const left = genNum(num.slice(0, i))
    const right = genNum(num.slice(i))
    for (const l of left) {
      for (const r of right) {
        res.push('(' + l + ', ' + r + ')')
      }
    }
    // console.log(left, right)
  }

  return res
}

console.log(ambiguousCoordinates('(123)'))
// 输出: ["(1, 23)", "(12, 3)", "(1.2, 3)", "(1, 2.3)"]
