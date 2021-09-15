// 使用两个数组的栈方法（push, pop） 实现队列
// pop的时候注意一下

class Queue {
  private input: unknown[]
  private output: unknown[]

  constructor() {
    this.input = []
    this.output = []
  }

  push(x: number) {
    this.input.push(x)
  }

  shift() {
    const size = this.output.length
    if (size) {
      return this.output.pop()
    }
    while (this.input.length) {
      this.output.push(this.input.pop())
    }
    return this.output.pop()
  }
}

const queue = new Queue()
queue.push(1)
queue.push(2)
queue.push(3)

console.log(queue.shift())
console.log(queue.shift())
console.log(queue.shift())
console.log(queue.shift())
