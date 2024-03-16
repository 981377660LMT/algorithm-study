/* eslint-disable no-inner-declarations */

/**
 * z算法求字符串每个后缀与原串的最长公共前缀长度.
 * @param arr 字符串或ascii码数组.
 * @returns z数组.
 * z[0]=0
 * z[i]是后缀s[i:]与s的最长公共前缀(LCP)的长度 (i>=1).
 */
function zAlgo(arr: ArrayLike<unknown>): Uint32Array {
  const n = arr.length
  const z = new Uint32Array(n)
  let left = 0
  let right = 0
  for (let i = 1; i < n; i++) {
    z[i] = Math.max(Math.min(z[i - left], right - i + 1), 0)
    while (i + z[i] < n && arr[z[i]] === arr[i + z[i]]) {
      left = i
      right = i + z[i]
      z[i]++
    }
  }
  return z
}

/**
 * 求s2的每个后缀与s1的最长公共前缀长度.
 */
function zAlgoTwoString(s1: string, s2: string): Uint32Array {
  const n1 = s1.length
  const n2 = s2.length
  const z = zAlgo(s1 + s2)
  for (let i = n1; i < n1 + n2; i++) {
    z[i] = Math.min(z[i], n1)
  }
  return z.subarray(n1)
}

export { zAlgo, zAlgoTwoString }

if (require.main === module) {
  // 2223. 构造字符串的总得分和
  // https://leetcode.cn/problems/sum-of-scores-of-built-strings/
  function sumScores(s: string): number {
    const z = zAlgo(s)
    const n = s.length
    return z.reduce((pre, cur) => pre + cur, 0) + n
  }

  console.log(zAlgoTwoString('aaaa', 'aabaa'))
}
