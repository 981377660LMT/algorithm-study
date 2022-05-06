// from datastructures-js
import { Queue } from 'datastructures-js'

// 每次删除头部元素并不是真的移除，而是标记其已经被移除
// 即用数组+左右两个指针 替代stack、queue
const queue = new Queue<number>()
queue.enqueue(1)
queue.enqueue(2)
console.log(queue.dequeue())

// 用数组替代stack、queue
class SimpleQueue<T = number> {
  private readonly data: T[]
  private head = 0
  private tail = 0

  constructor(size: number) {
    this.data = Array(size).fill(undefined)
  }

  get length(): number {
    return this.tail - this.head
  }

  push(element: T): SimpleQueue<T> {
    this.data[this.tail++] = element
    return this
  }

  shift(): T | undefined {
    if (this.length === 0) {
      return undefined
    }
    const front = this.data[this.head++]
    return front
  }

  at(index: number): T | undefined {
    return this.data[this.head + index]
  }
}
