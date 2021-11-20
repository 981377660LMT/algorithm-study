// 使用两个数组的栈方法（push, pop） 实现队列
// pop的时候注意一下

class MyQueue {
  private inputStack: unknown[]
  private outputStack: unknown[]

  constructor() {
    this.inputStack = []
    this.outputStack = []
  }

  push(x: number) {
    this.inputStack.push(x)
  }

  shift() {
    if (this.outputStack.length > 0) {
      return this.outputStack.pop()
    }
    while (this.inputStack.length > 0) {
      this.outputStack.push(this.inputStack.pop())
    }
    return this.outputStack.pop()
  }

  empty() {
    return this.inputStack.length === 0 && this.outputStack.length === 0
  }

  peek() {
    if (this.outputStack.length > 0) return this.outputStack[this.outputStack.length - 1]
    while (this.inputStack.length > 0) {
      this.outputStack.push(this.inputStack.pop())
    }
    return this.outputStack[this.outputStack.length - 1]
  }
}

const queue = new MyQueue()
queue.push(1)
queue.push(2)
queue.push(3)

console.log(queue.shift())
console.log(queue.shift())
console.log(queue.shift())
console.log(queue.shift())
