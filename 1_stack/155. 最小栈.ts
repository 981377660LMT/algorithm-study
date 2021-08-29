interface IStack {
  val: number
  min: number
}
class MinStack {
  private stack: IStack[]
  constructor() {
    this.stack = []
  }

  push(val: number): void {
    this.stack.push({ val, min: this.stack.length === 0 ? val : Math.min(val, this.getMin()) })
  }

  pop(): void {
    this.stack.pop()
  }

  top(): number {
    return this.stack[this.stack.length - 1].val
  }

  getMin(): number {
    return this.stack[this.stack.length - 1].min
  }
}

var obj = new MinStack()
obj.push(1)
obj.pop()
var param_3 = obj.top()
var param_4 = obj.getMin()

export {}
