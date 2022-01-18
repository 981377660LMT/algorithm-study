/**
 * @param {string} s
 * @return {number}
 * 给你一个只包含 '(' 和 ')' 的字符串，找出最长有效（格式正确且连续）括号子串的长度。
 */
const longestValidParentheses = function (s: string): number {
  // 用栈模拟一遍，将所有无法匹配的括号的位置全部置1
  // 例如: "()(()"的mark为[0, 0, 1, 0, 0]
  // 再例如: ")()((())"的mark为[1, 0, 0, 1, 0, 0, 0, 0]
  // 经过这样的处理后, 此题就变成了寻找最长的连续的0的长度
  const n = s.length
  const stack: number[] = []
  const mark = Array<number>(n).fill(0)
  for (let i = 0; i < n; i++) {
    if (s[i] === '(') {
      stack.push(i)
    } else {
      if (stack.length === 0) {
        mark[i] = 1
      } else {
        stack.pop()
      }
    }
  }

  // 多余的'('
  for (const remain of stack) {
    mark[remain] = 1
  }

  // 寻找最长的连续的0的长度
  return mark
    .join('')
    .split('1')
    .reduce((pre, cur) => Math.max(pre, cur.length), 0)
  // return
  // console.log(mark)
  // const markString = mark.join('')
  // const match = markString.match(/(0)\1*/g)
  // if (!match) return 0
  // return Math.max.apply(
  //   null,
  //   match.map(zeros => zeros.length)
  // )
}

console.log(longestValidParentheses('(()))()())(()'))

export default 1
