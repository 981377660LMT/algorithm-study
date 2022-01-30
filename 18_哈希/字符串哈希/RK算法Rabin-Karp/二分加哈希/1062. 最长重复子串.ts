// 滚动hash+二分答案

import { BigIntHasher } from '../../BigIntHasher'

/**
 *
 * @param s  1 <= S.length <= 1500   nlogn
 * @description
 * 给定字符串 S，找出最长重复子串的长度。如果不存在重复子串就返回 0。
 */
function longestRepeatingSubstring(s: string): number {
  const hasher = new BigIntHasher(s)
  BigIntHasher.setMOD(1 << 24)

  let left = 0
  let right = s.length
  while (left <= right) {
    const mid = (left + right) >> 1
    if (isExist(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  function isExist(length: number): boolean {
    if (length === 0) return true
    const visited = new Set<bigint>()

    for (let i = 1; i + length - 1 <= s.length; i++) {
      const sliceHash = hasher.getHashOfRange(i, i + length - 1)
      if (visited.has(sliceHash)) return true
      visited.add(sliceHash)
    }

    return false
  }
}

console.log(longestRepeatingSubstring('aabcaabdaab'))
// 输出：3
// 解释：最长的重复子串为 "aab"，出现 3 次。
console.log(longestRepeatingSubstring('abcd')) // 0
console.log(
  longestRepeatingSubstring(
    'bacbabcaacbbccacbacbbbccccbbcccabbabcabbbaaabbaccbcbcbcbacaaabbccbcacabcbccababca'
  )
) // 6
console.log(
  longestRepeatingSubstring('aaabaabbbaaabaabbaabbbabbbaaaabbaaaaaabbbaabbbbbbbbbaaaabbabbaba')
) // 10
