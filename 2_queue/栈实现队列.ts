// 栈实现队列
// 使用两个数组的栈方法（push, pop） 实现队列
// pop的时候注意一下

class MyQueue<T> {
  private inputStack: T[] = []
  private outputStack: T[] = []

  push(x: T) {
    this.inputStack.push(x)
  }

  shift() {
    if (this.outputStack.length) return this.outputStack.pop()
    while (this.inputStack.length) this.outputStack.push(this.inputStack.pop()!)
    return this.outputStack.pop()
  }

  empty() {
    return !this.inputStack.length && !this.outputStack.length
  }

  peek() {
    if (this.outputStack.length) return this.outputStack[this.outputStack.length - 1]
    while (this.inputStack.length) this.outputStack.push(this.inputStack.pop()!)
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
console.time('SimpleQueue')
const n = 1e7
for (let i = 0; i < n; i++) {
  queue.push(i)
}
for (let i = 0; i < n; i++) {
  queue.peek()
}
for (let i = 0; i < n; i++) {
  queue.shift()
}
console.timeEnd('SimpleQueue') // 600ms
export {}
