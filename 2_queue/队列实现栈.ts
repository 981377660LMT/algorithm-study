// pop的时候注意一下

class MyStack {
  private inputQueue: unknown[]
  private outputQueue: unknown[]

  constructor() {
    this.inputQueue = []
    this.outputQueue = []
  }

  push(x: number) {
    this.inputQueue.push(x)
  }

  pop() {
    while (this.inputQueue.length > 1) {
      this.outputQueue.push(this.inputQueue.shift())
    }
    const res = this.inputQueue.shift()
    ;[this.inputQueue, this.outputQueue] = [this.outputQueue, this.inputQueue]
    return res
  }

  empty() {
    return this.inputQueue.length === 0 && this.outputQueue.length === 0
  }

  top() {
    return this.inputQueue[this.inputQueue.length - 1]
  }
}

const stack = new MyStack()
stack.push(1)

console.log(stack.pop())
console.log(stack)
console.log(stack.empty())
