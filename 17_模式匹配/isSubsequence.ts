/**
 * 判断`shorter`是否是`longer`的子序列.
 * 如果需要多次匹配，使用`子序列自动机`.
 * @complexity O(n1 + n2)
 */
function isSubsequnceNaive<S extends ArrayLike<unknown>>(longer: S, shorter: S): boolean {
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

/**
 * 返回 {@link longer} 的每个前缀中的子序列匹配 {@link shorter} 的最大长度.
 *
 * @example
 * ```typescript
 * matchSubsequence('aabc', 'abc')
 * // [0, 1, 1, 2, 3]
 * ```
 */
function matchSubsequence<S extends ArrayLike<unknown>>(longer: S, shorter: S): number[] {
  const res = Array(longer.length + 1).fill(0)
  let i = 0
  let j = 0
  while (i < longer.length && j < shorter.length) {
    j += +(longer[i] === shorter[j])
    i++
    res[i] = j
  }
  res.fill(j, i + 1)
  return res
}

export { isSubsequnceNaive, matchSubsequence }

if (require.main === module) {
  // https://leetcode.cn/problems/is-subsequence/
  // eslint-disable-next-line no-inner-declarations
  function isSubsequence(s: string, t: string): boolean {
    return isSubsequnceNaive(t, s)
  }

  console.log(matchSubsequence('aabc', 'abc')) // [0, 1, 1, 2, 3]
}
