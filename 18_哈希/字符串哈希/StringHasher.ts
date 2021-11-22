import assert from 'assert'

interface IStringHasher {
  getHashOfRange(left: number, right: number): number
}

/**
 * @description
 * 如果使用Uint32Array并指定MOD为2**32 计算prefix和base时可以不用手动mod了
 * 但是计算区间哈希还是要模mod加mod再模mod
 * 生日悖论里面，如果想从 H 数集里面取 n 个值，那么产生碰撞的概率大约是 1 - math.exp(-(n^2)/(2*H))。
 */
class RabinKarpHasher implements IStringHasher {
  private static BASE = 131
  private static MOD = 1 << 24
  private readonly inputString: string
  private readonly prefix: number[]
  private readonly base: number[]

  static setMOD(mod: number): void {
    RabinKarpHasher.MOD = mod
  }

  static setBASE(base: number): void {
    RabinKarpHasher.BASE = base
  }

  constructor(input: string) {
    this.inputString = input
    this.prefix = Array(input.length + 1)
    this.base = Array(input.length + 1)
    this.prefix[0] = 0
    this.base[0] = 1

    for (let i = 1; i <= this.inputString.length; i++) {
      this.prefix[i] = this.prefix[i - 1] * RabinKarpHasher.BASE + input.charCodeAt(i - 1) - 97
      this.prefix[i] %= RabinKarpHasher.MOD
      this.base[i] = this.base[i - 1] * RabinKarpHasher.BASE
      this.base[i] %= RabinKarpHasher.MOD
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
  getHashOfRange(left: number, right: number): number {
    this.checkRange(left, right)
    const mod = RabinKarpHasher.MOD
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

    if (left > this.inputString.length) {
      throw new RangeError('left 不能 超出边界')
    }

    if (right > this.inputString.length) {
      throw new RangeError('right 不能 超出边界')
    }
  }
}

if (require.main === module) {
  const sh1 = new RabinKarpHasher('abcabc')
  assert.strictEqual(sh1.getHashOfRange(1, 1), sh1.getHashOfRange(4, 4))
  assert.strictEqual(sh1.getHashOfRange(1, 2), sh1.getHashOfRange(4, 5))

  const sh2 = new RabinKarpHasher('AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT'.toLowerCase())
  assert.strictEqual(sh2.getHashOfRange(1, 10), sh2.getHashOfRange(11, 20))
}

export { RabinKarpHasher, IStringHasher }
