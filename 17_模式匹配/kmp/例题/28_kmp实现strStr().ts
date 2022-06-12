import { getNext } from '../kmp'

/**
 * @param {string} long
 * @param {string} short
 * @return {number}
 * @description 在 pattern 字符串中找出 needle 字符串出现的第一个位置（下标从 0 开始）。
 * 如果不存在，则返回  -1 。如果needle是空字符串，则返回0。
 * 相当于实现indexOf
 * v8引擎中,indexOf使用了kmp和bm两种算法,在主串长度小于7时使用kmp,大于7的时候使用bm
 * @summary kmp比暴力解法好
 */
function indexofAll(long: string, short: string): number[] {
  if (short.length === 0) return [0]
  if (long.length < short.length) return []

  const res: number[] = []
  const next = getNext(short)
  let hitJ = 0
  for (let i = 0; i < long.length; i++) {
    while (hitJ > 0 && long[i] !== short[hitJ]) {
      hitJ = next[hitJ - 1]
    }

    if (long[i] === short[hitJ]) hitJ++

    // 找到头了
    if (hitJ === short.length) {
      res.push(i - short.length + 1)
      hitJ = next[hitJ - 1] // 不允许重叠时 hitJ = 0
    }
  }

  return res
}

if (require.main === module) {
  console.log(indexofAll('abcdaabcdfabcdababcdg', 'abcdab'))
  console.log(indexofAll('abcdaabcdfabcdababcdg', 'ab'))
}

// 10
export {}
