interface StackItem {
  val: number
  max: number
}

class MaxStack {
  private stack: StackItem[]

  constructor() {
    this.stack = []
  }

  push(val: number): void {
    this.stack.push({
      val,
      max: this.stack.length === 0 ? val : Math.max(val, this.stack[this.stack.length - 1].max),
    })
  }

  pop(): number {
    return this.stack.pop()!.val
  }

  top(): number {
    return this.stack[this.stack.length - 1].val
  }

  peekMax(): number {
    return this.stack[this.stack.length - 1].max
  }

  // 检索并返回栈中最大元素，并将其移除。如果有多个最大元素，只要移除 最靠近栈顶 的那个。
  // popmax 就是倒出来，弹出最大后，再回去 注意回去时调用的MaxStack.push方法
  popMax(): number {
    const res = this.peekMax()
    const tmpStack: number[] = []

    while (this.top() !== res) {
      tmpStack.push(this.stack.pop()!.val)
    }

    this.stack.pop()

    while (tmpStack.length > 0) {
      this.push(tmpStack.pop()!)
    }

    return res
  }
}
/**
 * Your MaxStack object will be instantiated and called as such:
 * var obj = new MaxStack()
 * obj.push(x)
 * var param_2 = obj.pop()
 * var param_3 = obj.top()
 * var param_4 = obj.peekMax()
 * var param_5 = obj.popMax()
 */
const stk = new MaxStack()
stk.push(5)
stk.push(1)
console.log(stk.popMax())
console.log(stk)
console.log(stk.peekMax())

export {}
