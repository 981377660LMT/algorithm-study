interface IStringHasher<R extends number | bigint> {
  getHashOfSlice(left: number, right: number): R
}

/**
 * @description
 * 哈希值计算方法：
 * hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
 * 越靠左字符权重越大
 */
class BigIntHasher implements IStringHasher<bigint> {
  private static BASE = 131n
  private static OFFSET = 97n
  private static MOD = BigInt(2 ** 64)
  private readonly input: string
  private readonly prefix: BigUint64Array
  private readonly base: BigUint64Array

  static setBASE(base: number): void {
    BigIntHasher.BASE = BigInt(base)
  }

  static setMOD(mod: number): void {
    BigIntHasher.MOD = BigInt(mod)
  }

  static setOFFSET(offset: number): void {
    BigIntHasher.OFFSET = BigInt(offset)
  }

  constructor(input: string) {
    this.input = input
    this.prefix = new BigUint64Array(input.length + 1)
    this.base = new BigUint64Array(input.length + 1)
    this.prefix[0] = 0n
    this.base[0] = 1n

    for (let i = 1; i <= this.input.length; i++) {
      this.prefix[i] =
        this.prefix[i - 1] * BigIntHasher.BASE +
        BigInt(input[i - 1].codePointAt(0)!) -
        BigIntHasher.OFFSET
      this.prefix[i] %= BigIntHasher.MOD
      this.base[i] = this.base[i - 1] * BigIntHasher.BASE
      this.base[i] %= BigIntHasher.MOD
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
    const mod = BigIntHasher.MOD
    const upper = this.prefix[right]
    const lower = this.prefix[left - 1] * this.base[right - (left - 1)]
    return (upper - (lower % mod) + mod) % mod
  }

  private checkRange(left: number, right: number): void {
    if (0 <= left && left <= right && right <= this.input.length) return
    throw new RangeError('left or right out of range')
  }
}

export { BigIntHasher, IStringHasher }

if (require.main === module) {
  const stringHasher = new BigIntHasher('abcdefg')
  console.log(stringHasher.getHashOfSlice(1, 2))
  console.log(stringHasher.getHashOfSlice(2, 2))
}
