class ProductOfNumbers {
  private pre: number[]

  constructor() {
    this.pre = [1]
  }

  // 将数字 num 添加到当前数字列表的最后面
  add(num: number): void {
    if (num === 0) this.pre = [1]
    else {
      this.pre.push(num * this.pre[this.pre.length - 1])
    }
  }

  // 返回当前数字列表中，最后 k 个数字的乘积。
  // 你可以假设当前列表中始终 至少 包含 k 个数字
  getProduct(k: number): number {
    const len = this.pre.length
    if (len - k >= 1) {
      return this.pre[len - 1] / this.pre[len - k - 1]
    }
    // pre 数组的长度小于 k，说明末尾 k 个数里肯定有 0，直接输出 0 即可
    return 0
  }
}

export {}

// 如果是求任意区间怎么办
class AdvancedProductOfNumbers {
  private pre: number[]
  private zeroPre: number[] // 前缀和记录0的个数

  constructor() {
    this.pre = [1]
    this.zeroPre = [0]
  }

  // 将数字 num 添加到当前数字列表的最后面
  add(num: number): void {
    if (num === 0) {
      this.pre = [1]
      this.zeroPre.push(this.zeroPre[this.zeroPre.length - 1] + 1)
    } else {
      this.pre.push(num * this.pre[this.pre.length - 1])
      this.zeroPre.push(this.zeroPre[this.zeroPre.length - 1])
    }
  }

  // 返回当前数字列表中，第s个元素到第e个元素间的乘积个数字的乘积。
  // 你可以假设当前列表中始终 至少 包含 k 个数字
  getProduct(start: number, end: number): number {
    if (this.zeroPre[start - 1] !== this.zeroPre[end]) return 0
    return this.pre[end] / this.pre[start - 1]
  }
}
