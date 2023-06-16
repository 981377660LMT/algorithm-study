/**
 * 判断`shorter`是否是`longer`的子序列.
 * 如果需要多次匹配，使用`子序列自动机`.
 * @complexity O(n1 + n2)
 */
function isSubsequnceNaive(longer: string, shorter: string): boolean {
  const n1 = longer.length
  const n2 = shorter.length
  if (!n2) return true
  if (n1 < n2) return false

  let i = 0
  let j = 0
  while (i < n1 && j < n2) {
    if (longer[i] === shorter[j]) j++
    if (j === n2) return true
    i++
  }
  return false
}

export { isSubsequnceNaive }

if (require.main === module) {
  // https://leetcode.cn/problems/is-subsequence/
  // eslint-disable-next-line no-inner-declarations
  function isSubsequence(s: string, t: string): boolean {
    return isSubsequnceNaive(t, s)
  }
}
