/**
 * 所有给定的字符串长度不会超过 10 。
 * 给定字符串列表的长度将在 [2, 50 ] 之间。
 * @return {number}
 * 「最长特殊序列」定义如下：该序列为某字符串独有的最长子序列（即不能是其他字符串的子序列）。
 * @summary
 * 暴力解法即可
 */
function findLUSlength(strs: string[]): number {
  strs.sort((a, b) => b.length - a.length)
  for (const cur of strs) {
    const count = strs.filter(other => isSubsequence(other, cur)).length
    if (count === 1) return cur.length
  }

  return -1

  function isSubsequence(longer: string, shorter: string): boolean {
    let i = 0
    let j = 0

    while (i < longer.length && j < shorter.length) {
      if (longer[i] === shorter[j]) j++
      if (j === shorter.length) return true
      i++
    }

    return false
  }
}

console.log(findLUSlength(['aba', 'cdc', 'eae']))

export default 1
// s1 = 'ab',s2 = 'a',因为ab是s1独有，所以最长子序列为ab，
// s1 = 'ab', s2 = 'ab', 因为ab是两个串都有，ab排除，
// a也是两个串都有，排除，b也是两个串都有，排除。所以最长特殊序列不存在，返回-1
// 通过以上分析，我们可以得出结论，如果：两个串相等（不仅长度相等，内容也相等），那么他们的最长特殊序列不存在。返回-1
// 如果两个串长度不一样，那么长的串 永远也不可能是 短串的子序列，即len(s1) > len(s2),则最长特殊序列为s1,返回长度大的数 `
