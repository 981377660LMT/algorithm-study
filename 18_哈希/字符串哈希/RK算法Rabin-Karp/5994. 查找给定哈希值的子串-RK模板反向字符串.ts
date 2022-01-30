// 题目里：hash(s, p, m) = (val(s[0]) * p0 + val(s[1]) * p1 + ... + val(s[k-1]) * pk-1) mod m.
// 注意我们的RK算法里计算哈希值的方法是左边字符权重大，题目是右边权重大
// 所以要把我们的字符串反过来，调api，哈希值相等时返回这一段的reversed

interface IStringHasher {
  getHashOfRange(left: number, right: number): number | bigint
}

class BigIntHasher implements IStringHasher {
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
    this.prefix[0] = 1n
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
   * @returns 闭区间 [left,right] 子串的哈希值  left right `从1开始`
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

function subStrHash(
  s: string,
  power: number,
  modulo: number,
  k: number,
  hashValue: number
): string {
  BigIntHasher.setBASE(power)
  BigIntHasher.setMOD(modulo)
  s = s.split('').reverse().join('')
  const hasher = new BigIntHasher(s)

  let res = 0
  for (let i = 0; i + k <= s.length; i++) {
    const hash = hasher.getHashOfRange(i + 1, i + k)
    if (Number(hash) === hashValue) res = i
  }

  return s
    .slice(res, res + k)
    .split('')
    .reverse()
    .join('')
}

console.log(subStrHash('leetcode', 7, 20, 2, 0))
console.log(subStrHash('fbxzaad', 31, 100, 3, 32))

export {}
