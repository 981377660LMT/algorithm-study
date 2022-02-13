interface IStringHasher<R extends number | bigint> {
  getHashOfRange(left: number, right: number): R
}

/**
 * @description
 * 哈希值计算方法：
 * hash(s, p, m) = (val(s[0]) * pk-1 + val(s[1]) * pk-2 + ... + val(s[k-1]) * p0) mod m.
 * 越靠左字符权重越大
 */
class BigIntHasher implements IStringHasher<bigint> {
  private static BASE = 131n
  private static MOD = BigInt(2 ** 64)
  private readonly input: ArrayLike<string>
  private readonly prefix: BigUint64Array
  private readonly base: BigUint64Array

  static setBASE(base: number): void {
    BigIntHasher.BASE = BigInt(base)
  }

  static setMOD(mod: number): void {
    BigIntHasher.MOD = BigInt(mod)
  }

  constructor(input: ArrayLike<string>) {
    this.input = input
    this.prefix = new BigUint64Array(input.length + 1)
    this.base = new BigUint64Array(input.length + 1)
    this.prefix[0] = 0n
    this.base[0] = 1n

    for (let i = 1; i <= this.input.length; i++) {
      this.prefix[i] =
        this.prefix[i - 1] * BigIntHasher.BASE + BigInt(input[i - 1].codePointAt(0)!) - 96n
      this.prefix[i] %= BigIntHasher.MOD
      this.base[i] = this.base[i - 1] * BigIntHasher.BASE
      this.base[i] %= BigIntHasher.MOD
    }
  }

  /**
   *
   * @param left
   * @param right
   * @returns 闭区间 [left,right] 子串的哈希值  left>=1 right<=n
   * @description
   * 注意要`模mod加mod再模mod`
   */
  getHashOfRange(left: number, right: number): bigint {
    this.checkRange(left, right)
    const mod = BigIntHasher.MOD
    const upper = this.prefix[right]
    const lower = this.prefix[left - 1] * this.base[right - (left - 1)]
    return (upper - (lower % mod) + mod) % mod
  }

  private checkRange(left: number, right: number) {
    if (right < left) {
      throw new RangeError('right 不能小于 left')
    }

    if (left < 1) {
      throw new RangeError('left 不能小于1')
    }

    if (right < 1) {
      throw new RangeError('right 不能小于1')
    }

    if (left > this.input.length) {
      throw new RangeError('left 不能 超出边界')
    }

    if (right > this.input.length) {
      throw new RangeError('right 不能 超出边界')
    }
  }
}

export { BigIntHasher }
if (require.main === module) {
  const stringHasher = new BigIntHasher('abcdefg')
  console.log(stringHasher.getHashOfRange(1, 1))
  console.log(stringHasher.getHashOfRange(2, 2))
}
