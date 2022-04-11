// 滚动hash+二分答案

import { IStringHasher } from '../../BigIntHasher'

/**
 * @description
 * 哈希值计算方法：
 * hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
 * 越靠左字符权重越大
 */
class ArrayHasher implements IStringHasher<bigint> {
  private static BASE = 131n
  private static OFFSET = 97n

  private static MOD = BigInt(2 ** 64)
  private readonly input: number[]
  private readonly prefix: BigUint64Array
  private readonly base: BigUint64Array

  static setBASE(base: number): void {
    ArrayHasher.BASE = BigInt(base)
  }

  static setMOD(mod: number): void {
    ArrayHasher.MOD = BigInt(mod)
  }

  static setOFFSET(offset: number): void {
    ArrayHasher.OFFSET = BigInt(offset)
  }

  constructor(input: number[]) {
    this.input = input
    this.prefix = new BigUint64Array(input.length + 1)
    this.base = new BigUint64Array(input.length + 1)
    this.prefix[0] = 0n
    this.base[0] = 1n

    for (let i = 1; i <= this.input.length; i++) {
      this.prefix[i] =
        this.prefix[i - 1] * ArrayHasher.BASE + BigInt(input[i - 1]) - ArrayHasher.OFFSET
      this.prefix[i] %= ArrayHasher.MOD
      this.base[i] = this.base[i - 1] * ArrayHasher.BASE
      this.base[i] %= ArrayHasher.MOD
    }
  }

  /**
   *
   * @param left
   * @param right
   * @returns 切片 [left:right] 的哈希值
   */
  getHashOfSlice(left: number, right: number): bigint {
    if (left === right) return 0n
    left += 1
    this.checkRange(left, right)
    const mod = ArrayHasher.MOD
    const upper = this.prefix[right]
    const lower = this.prefix[left - 1] * this.base[right - (left - 1)]
    return (upper - (lower % mod) + mod) % mod
  }

  private checkRange(left: number, right: number): void {
    if (0 <= left && left <= right && right <= this.input.length) return
    throw new RangeError('left or right out of range')
  }
}

/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}  dp[i][j] 为 以 A[i], B[j] 结尾的两个数组中公共的、长度最长的子数组的长度
 */
function findLength(nums1: number[], nums2: number[]): number {
  const hasher1 = new ArrayHasher(nums1)
  const hasher2 = new ArrayHasher(nums2)
  let left = 0
  let right = Math.min(nums1.length, nums2.length)

  while (left <= right) {
    const mid = (left + right) >> 1
    if (isExist(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  function isExist(len: number): boolean {
    if (len === 0) return true
    const visited = new Set<bigint>()

    for (let left = 0; left + len <= nums1.length; left++) {
      const hash = hasher1.getHashOfSlice(left, left + len)
      console.log(len, left, hash)
      visited.add(hash)
    }

    for (let left = 0; left + len <= nums2.length; left++) {
      const hash = hasher2.getHashOfSlice(left, left + len)
      if (visited.has(hash)) return true
    }

    return false
  }
}

// console.log(findLength([1, 2, 3, 2, 1], [3, 2, 1, 4, 7]))
console.log(findLength([70, 39, 25, 40, 7], [52, 20, 67, 5, 31])) // 0
// console.log(findLength([1, 0, 1, 0, 0, 0, 0, 0, 1, 1], [1, 1, 0, 1, 1, 0, 0, 0, 0, 0])) // 6
// console.log(findLength([0, 1, 1, 0, 1, 1, 1, 0, 1, 0], [1, 0, 0, 0, 1, 0, 0, 1, 1, 0])) // 4
// export {}
