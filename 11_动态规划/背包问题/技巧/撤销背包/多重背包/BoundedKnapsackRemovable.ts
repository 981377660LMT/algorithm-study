/**
 * 可撤销多重背包，用于求解方案数/可行性.
 */
class BoundedKnapsackRemovable {
  private readonly _dp: number[]
  private readonly _mod?: number | undefined
  private readonly _maxValue: number
  private _countSum = 0

  constructor(maxValue: number, mod?: number, dp?: number[]) {
    if (dp != undefined) {
      this._dp = dp
    } else {
      this._dp = Array(maxValue + 1).fill(0)
      this._dp[0] = 1
    }
    this._maxValue = maxValue
    this._mod = mod
  }

  /**
   * 加入一个价值为value(value>0)的物品，数量为count.
   * @complexity O(maxValue)
   */
  add(value: number, count: number): void {
    if (value <= 0) throw new Error(`value must be positive, but got ${value}`)
    this._countSum += count * value
    const maxJ = Math.min(this._countSum, this._maxValue)
    if (this._mod == undefined) {
      for (let j = value; j <= maxJ; j++) {
        this._dp[j] += this._dp[j - value]
      }
      for (let j = maxJ; j >= value * (count + 1); j--) {
        this._dp[j] -= this._dp[j - value * (count + 1)]
      }
    } else {
      for (let j = value; j <= maxJ; j++) {
        this._dp[j] = (this._dp[j] + this._dp[j - value]) % this._mod
      }
      for (let j = maxJ; j >= value * (count + 1); j--) {
        this._dp[j] = (this._dp[j] - this._dp[j - value * (count + 1)]) % this._mod
        if (this._dp[j] < 0) this._dp[j] += this._mod
      }
    }
  }

  /**
   * 移除一个价值为value(value>0)的物品，数量为count.需要保证value物品存在.
   * @complexity O(maxValue)
   */
  remove(value: number, count: number): void {
    const maxJ = Math.min(this._countSum, this._maxValue)
    if (this._mod == undefined) {
      for (let i = (count + 1) * value; i <= maxJ; i++) {
        this._dp[i] += this._dp[i - (count + 1) * value]
      }
      for (let i = maxJ; i >= value; i--) {
        this._dp[i] -= this._dp[i - value]
      }
    } else {
      for (let i = (count + 1) * value; i <= maxJ; i++) {
        this._dp[i] = (this._dp[i] + this._dp[i - (count + 1) * value]) % this._mod
      }
      for (let i = maxJ; i >= value; i--) {
        this._dp[i] = (this._dp[i] - this._dp[i - value]) % this._mod
        if (this._dp[i] < 0) this._dp[i] += this._mod
      }
    }

    this._countSum -= count * value
  }

  query(value: number): number {
    if (value < 0 || value > this._maxValue) return 0
    if (this._mod == undefined) return this._dp[value]
    if (this._dp[value] < 0) this._dp[value] += this._mod
    return this._dp[value]
  }

  copy(): BoundedKnapsackRemovable {
    const res = new BoundedKnapsackRemovable(this._maxValue, this._mod, this._dp.slice())
    res._countSum = this._countSum
    return res
  }
}

export { BoundedKnapsackRemovable }

if (require.main === module) {
  const k = new BoundedKnapsackRemovable(10)
  k.add(2, 5)
  k.add(1, 10)
  console.log(k.query(10))
  k.remove(1, 10)
  console.log(k.query(10))

  const tmp = k.copy()
  tmp.add(1, 10)
  console.log(tmp.query(10))
  console.log(k.query(10))
}
