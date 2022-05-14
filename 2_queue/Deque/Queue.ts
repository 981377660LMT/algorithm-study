// from datastructures-js
// import { Queue } from 'datastructures-js'

class SimpleQueue<T = number> {
  private readonly data: T[] = []
  private head = 0
  private tail = 0

  constructor(iterable?: Iterable<T>) {
    for (const item of iterable ?? []) {
      this.push(item)
    }
  }

  get length(): number {
    return this.tail - this.head
  }

  push(element: T): SimpleQueue<T> {
    this.data[this.tail++] = element
    return this
  }

  shift(): T | undefined {
    if (this.length === 0) return undefined
    const front = this.data[this.head++]
    return front
  }

  at(index: number): T | undefined {
    if (index < 0) index += this.length
    if (index < 0 || index >= this.length) return undefined
    return this.data[this.head + index]
  }
}

if (require.main === module) {
  // 每次删除头部元素并不是真的移除，而是标记其已经被移除
  // 即用数组+左右两个指针 替代stack、queue
  const queue = new SimpleQueue<number>()
  queue.push(1)
  queue.push(2)
  console.log(queue.shift())
  console.log(queue.at(-1))
  console.log(queue.at(-2))
  console.log(queue.at(0))
  console.log(queue.at(1))
}

export { SimpleQueue }
