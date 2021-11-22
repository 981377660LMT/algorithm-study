class BigIntHasher<T extends ArrayLike<bigint | boolean | number | string>> {
  private static BASE = 131n
  private static MOD = BigInt(2 ** 64)
  private readonly input: T
  private readonly prefix: BigUint64Array
  private readonly base: BigUint64Array

  static setBASE(base: number): void {
    BigIntHasher.BASE = BigInt(base)
  }

  static setMOD(mod: number): void {
    BigIntHasher.MOD = BigInt(mod)
  }

  constructor(input: T) {
    this.input = input
    this.prefix = new BigUint64Array(input.length + 1)
    this.base = new BigUint64Array(input.length + 1)
    this.prefix[0] = 0n
    this.base[0] = 1n

    for (let i = 1; i <= this.input.length; i++) {
      this.prefix[i] = this.prefix[i - 1] * BigIntHasher.BASE + BigInt(input[i - 1])
      this.base[i] = this.base[i - 1] * BigIntHasher.BASE
    }
  }

  /**
   *
   * @param left
   * @param right
   * @returns 闭区间 [left,right] 子串的哈希值  left right `从1开始`
   * @description
   * 注意要`模mod加mod再模mod`
   */
  getHashOfRange(left: number, right: number): bigint {
    // this.checkRange(left, right)
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
