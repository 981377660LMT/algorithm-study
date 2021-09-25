class CQueue {
  private input: number[]
  private output: number[]
  constructor() {
    this.input = []
    this.output = []
  }

  appendTail(value: number): void {
    this.input.push(value)
  }
  // 若队列中没有元素，deleteHead 操作返回 -1
  deleteHead(): number {
    const size = this.output.length
    if (size) return this.output.pop()!
    while (this.input.length) {
      this.output.push(this.input.pop()!)
    }
    return this.output.pop() || -1
  }
}

// push(x: number) {
//   this.input.push(x)
// }

// pop() {
//   const size = this.output.length
//   if (size) {
//     return this.output.shift()
//   }
//   while (this.input.length) {
//     this.output.push(this.input.pop())
//   }
//   return this.output.shift()
// }
