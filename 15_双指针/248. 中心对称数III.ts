import { findStrobogrammatic } from './247. 中心对称数 II'

/**
 * @param {number} n
 * @return {string[]}
 * @description
 * 来计算范围在 [low, high] 之间中心对称数的个数。
 * 由于范围可能很大，所以 low 和 high 都用字符串表示。
 */
function findStrobogrammatic2(low: string, high: string): string[] {
  const [n1, n2] = [low.length, high.length]
  const [minLength, maxLength] = [Math.min(n1, n2), Math.max(n1, n2)]
  const res: string[] = []

  for (let n = minLength; n <= maxLength; n++) {
    for (const str of findStrobogrammatic(n)) {
      const num = parseInt(str)
      if (parseInt(low) <= num && num <= parseInt(high)) res.push(str)
    }
  }

  return res
}

console.log(findStrobogrammatic2('50', '100'))
