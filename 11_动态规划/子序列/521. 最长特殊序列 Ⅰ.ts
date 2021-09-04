/**
 * @param {string} a
 * @param {string} b
 * @return {number}
 * 「最长特殊序列」定义如下：该序列为某字符串独有的最长子序列（即不能是其他字符串的子序列）。
 */
var findLUSlength = function (a: string, b: string): number {
  if (a === b) return -1
  else return Math.max(a.length, b.length)
}

console.log(findLUSlength('aba', 'cdc'))

export default 1
// s1 = 'ab',s2 = 'a',因为ab是s1独有，所以最长子序列为ab，
// s1 = 'ab', s2 = 'ab', 因为ab是两个串都有，ab排除，
// a也是两个串都有，排除，b也是两个串都有，排除。所以最长特殊序列不存在，返回-1
// 通过以上分析，我们可以得出结论，如果：两个串相等（不仅长度相等，内容也相等），那么他们的最长特殊序列不存在。返回-1
// 如果两个串长度不一样，那么长的串 永远也不可能是 短串的子序列，即len(s1) > len(s2),则最长特殊序列为s1,返回长度大的数 `
