/**
 * 可撤销完全背包，用于求解方案数/可行性.
 */
class UnboundedKnapsackRemovable {
  private readonly _dp: number[]
  private readonly _maxWeight: number
  private readonly _mod?: number | undefined

  constructor(maxWeight: number, mod?: number, dp?: number[]) {
    if (dp != undefined) {
      this._dp = dp
    } else {
      this._dp = Array(maxWeight + 1).fill(0)
      this._dp[0] = 1
    }
    this._maxWeight = maxWeight
    this._mod = mod
  }

  /**
   * 添加一个重量为weight的物品.
   * @complexity O(maxWeight)
   */
  add(weight: number): void {
    if (this._mod == undefined) {
      for (let i = weight; i <= this._maxWeight; i++) {
        this._dp[i] += this._dp[i - weight]
      }
    } else {
      for (let i = weight; i <= this._maxWeight; i++) {
        this._dp[i] = (this._dp[i] + this._dp[i - weight]) % this._mod
      }
    }
  }

  /**
   * 移除一个重量为weight的物品.需要保证weight物品存在.
   * @complexity O(maxWeight)
   */
  remove(weight: number): void {
    if (this._mod == undefined) {
      for (let i = this._maxWeight; i >= weight; i--) {
        this._dp[i] -= this._dp[i - weight]
      }
    } else {
      for (let i = this._maxWeight; i >= weight; i--) {
        this._dp[i] = (this._dp[i] - this._dp[i - weight]) % this._mod
      }
    }
  }

  /**
   * 查询组成重量为weight的物品有多少种方案.
   * @complexity O(1)
   */
  query(weight: number): number {
    if (weight < 0 || weight > this._maxWeight) return 0
    if (this._mod == undefined) return this._dp[weight]
    if (this._dp[weight] < 0) this._dp[weight] += this._mod
    return this._dp[weight]
  }

  copy(): UnboundedKnapsackRemovable {
    return new UnboundedKnapsackRemovable(this._maxWeight, this._mod, this._dp.slice())
  }
}

export { UnboundedKnapsackRemovable }

if (require.main === module) {
  const knapsack = new UnboundedKnapsackRemovable(10)
  knapsack.add(1)
  knapsack.add(2)
  knapsack.remove(2)
  console.log(knapsack.query(10))
}
