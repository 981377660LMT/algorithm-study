/**
 * @param {string} s
 * @param {string} c
 * @return {number[]}
 * @description
 * 给定一个字符串 S 和一个字符 C。返回一个代表字符串 S 中每个字符到字符串 S 中的字符 C 的最短距离的数组。
 * 题目数据保证 c 在 s 中至少出现一次
 */
const shortestToChar = function (s: string, c: string): number[] {
  const len = s.length
  const res = Array<number>(len).fill(0)
  const stack: number[] = []
  for (let i = len - 1; i >= 0; i--) {
    if (s[i] === c) stack.push(i)
  }

  let last = Infinity
  for (let i = 0; i < len; i++) {
    if (stack.length && stack[stack.length - 1] <= i) {
      last = stack.pop()!
    }
    // 这种技巧要熟练
    const next = stack[stack.length - 1] || Infinity
    res[i] = Math.min(Math.abs(i - last), Math.abs(i - next))
  }

  return res
}

console.log(shortestToChar('loveleetcode', 'e'))
