/* eslint-disable func-names */

/**
 * @param {string} a
 * @param {string} b
 * @return {number}
 * 寻找重复叠加字符串 a 的最小次数，
 * 使得字符串 b 成为叠加后的字符串 a 的子串，如果不存在则返回 -1。
 *
 * @summary
 * !最坏情况： 一个完整的B可能首部用到A的一部分，尾部用到A的一部分，
 */
function repeatedStringMatch(a: string, b: string): number {
  const count = Math.ceil(b.length / a.length)
  const str = a.repeat(count)
  if (str.includes(b)) return count
  return (str + a).includes(b) ? count + 1 : -1
}

console.log(repeatedStringMatch('abcd', 'cdabcdab')) // 3
console.log(repeatedStringMatch('dab', 'abbd')) // 3

export default 1
