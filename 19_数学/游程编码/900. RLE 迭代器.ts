class RLEIterator {
  private encoding: number[]
  private index: number
  /**
   *
   * @param encoding
   * @description
   * 对于所有偶数 i，A[i] 告诉我们在序列中重复非负整数值 A[i + 1] 的次数。
   */
  constructor(encoding: number[]) {
    this.encoding = encoding
    this.index = 0
  }

  /**
   *
   * @param n
   * @description
   * 它耗尽接下来的  n 个元素（n >= 1）并返回以这种方式耗去的最后一个元素。
   * 如果没有剩余的元素可供耗尽，则  next 返回 -1 。
   */
  next(n: number): number {
    while (this.index < this.encoding.length && this.encoding[this.index] < n) {
      n -= this.encoding[this.index]
      this.index += 2
    }

    if (this.index >= this.encoding.length) return -1

    this.encoding[this.index] -= n
    return this.encoding[this.index + 1]
  }
}

/**
 * Your RLEIterator object will be instantiated and called as such:
 * var obj = new RLEIterator(encoding)
 * var param_1 = obj.next(n)
 */

const rle = new RLEIterator([3, 8, 0, 9, 2, 5])

console.log(rle.next(2))
console.log(rle.next(1))
console.log(rle.next(1))
console.log(rle.next(2))
// 输出：[8,8,5,-1]
