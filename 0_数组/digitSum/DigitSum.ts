/* eslint-disable no-inner-declarations */

// FastDigitSum
// 计算一个数字各位digit之和
// !可以加速计算

class DigitSum {
  private readonly _mod: number
  /** 长为10^step的数组, dp[x]表示x的各位数字之和. */
  private readonly _dp: Uint8Array

  constructor(step = 6) {
    step = Math.max(4, Math.min(7, step))
    this._mod = 10 ** step
    this._dp = new Uint8Array(this._mod)
    for (let x = 1; x < this._mod; x++) {
      this._dp[x] = this._dp[Math.floor(x / 10)] + (x % 10)
    }
  }

  sum(x: number): number {
    let res = 0
    const dp = this._dp
    const mod = this._mod
    while (x > 0) {
      res += dp[x % mod]
      x = Math.floor(x / mod)
    }
    return res
  }
}

function digitSumNaive(num: number): number {
  let sum = 0
  while (num > 0) {
    sum += num % 10
    num = Math.floor(num / 10)
  }
  return sum
}

function digitSumNavie32(num: number): number {
  let sum = 0
  while (num > 0) {
    sum += num % 10
    num = (num / 10) | 0
  }
  return sum
}

export { DigitSum, digitSumNaive, digitSumNavie32 }

if (require.main === module) {
  function test() {
    console.time('digitSumNaive')
    for (let i = 0; i < 1e8; i++) {
      // digitSumNaive(i) // digitSumNaive: 4.333s
      digitSumNavie32(i) // digitSumNaive: 1.536s
    }
    console.timeEnd('digitSumNaive')

    console.time('DigitSum')
    const ds = new DigitSum()
    for (let i = 0; i < 1e8; i++) {
      ds.sum(i) // DigitSum: 1.247s
    }
    console.timeEnd('DigitSum')
  }

  test()
}
