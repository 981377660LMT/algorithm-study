// # 给你一个二进制字符串 s 和一个整数 k 。
// # 如果所有长度为 k 的二进制字符串都是 s 的子串，请返回 true ，否则请返回 false 。
// # 1 <= k <= 20
// # 1 <= s.length <= 5 * 105

import { BigIntHasher } from '../BigIntHasher'

function hasAllCodes(s: string, k: number): boolean {
  const leftHasher = new BigIntHasher(s)
  BigIntHasher.setBASE(2)

  const visited = new Set<bigint>()
  for (let left = 1; left + k - 1 <= s.length; left++) {
    const hash = leftHasher.getHashOfRange(left, left + k - 1)
    visited.add(hash)
  }

  return visited.size === 1 << k
}

console.log(hasAllCodes('00110110', 2))
