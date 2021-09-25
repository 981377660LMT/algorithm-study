const mod = 10 ** 9 + 7

class Fancy {
  private nums: number[]
  private add: number[]
  private multi: number[]

  constructor() {
    this.nums = []
    this.add = [0] // 前缀
    this.multi = [1] // 前缀
  }

  append(val: number): void {
    this.nums.push(val)
    this.add.push(this.add[this.add.length - 1])
    this.multi.push(this.multi[this.multi.length - 1])
  }

  addAll(inc: number): void {
    this.add[this.add.length - 1] += inc
  }

  multAll(m: number): void {
    this.add[this.add.length - 1] *= m
    this.add[this.add.length - 1] %= mod
    this.multi[this.multi.length - 1] *= m
    this.multi[this.multi.length - 1] %= mod
  }

  getIndex(idx: number): number {
    if (idx > this.nums.length) return -1
    // 逆元
    const inverse = Math.pow(this.multi[idx], mod - 2) % mod
    // https://leetcode-cn.com/problems/fancy-sequence/solution/onlognqian-zhui-he-qian-zhui-ji-shu-zu-by-simpleso/
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
