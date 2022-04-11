import { IStringHasher } from './BigIntHasher'

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
