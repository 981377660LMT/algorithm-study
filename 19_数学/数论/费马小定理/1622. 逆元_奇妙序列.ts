const MOD = 1e9 + 7

// 范围乘操作
// 总共最多会有 105 次对 append，addAll，multAll 和 getIndex 的调用。
class Fancy {
  private nums: number[]
  /**
   * @description add 和 mul 都是记录 在 index 位置上的 操作
   */
  private add: number[]
  private mul: number[]

  constructor() {
    this.nums = []
    this.add = [0] // 前缀
    this.mul = [1] // 前缀
  }

  append(val: number): void {
    this.nums.push(val)
    this.add.push(this.add[this.add.length - 1])
    this.mul.push(this.mul[this.mul.length - 1])
  }

  addAll(inc: number): void {
    this.add[this.add.length - 1] += inc
  }

  multAll(m: number): void {
    this.add[this.add.length - 1] = (this.add[this.add.length - 1] * m) % MOD
    this.mul[this.mul.length - 1] = (this.mul[this.mul.length - 1] * m) % MOD
  }

  getIndex(idx: number): number {
    if (idx >= this.nums.length) return -1
    // 逆元
    // mul[-1]/mul[i]化为 mul[-1]*inv(mul[i])
    const 乘的倍数 =
      (this.mul[this.mul.length - 1] * (Math.pow(this.mul[idx], MOD - 2) % MOD)) % MOD
    const 加了多少 = this.add[this.add.length - 1] - this.add[idx] * 乘的倍数
    return (this.nums[idx] * 乘的倍数 + 加了多少) % MOD
  }
}

export {}

/**
 * Your Fancy object will be instantiated and called as such:
 * var obj = new Fancy()
 * obj.append(val)
 * obj.addAll(inc)
 * obj.multAll(m)
 * var param_4 = obj.getIndex(idx)
 */
