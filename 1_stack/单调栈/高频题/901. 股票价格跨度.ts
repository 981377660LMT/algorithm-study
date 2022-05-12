class StockSpanner {
  constructor(private stack: [price: number, gap: number][] = []) {}

  /**
   *
   * @param price
   * @returns 股票价格小于或等于今天价格的最大连续日数
   */
  next(price: number): number {
    let gap = 1
    while (this.stack.length > 0 && this.stack[this.stack.length - 1][0] <= price) {
      gap += this.stack.pop()![1]
    }
    this.stack.push([price, gap])
    return gap
  }
}

const ss = new StockSpanner()
console.log(ss.next(100))
console.log(ss.next(80))
console.log(ss.next(60))
console.log(ss.next(70))
console.log(ss.next(60))
console.log(ss.next(75))
console.log(ss.next(85))
