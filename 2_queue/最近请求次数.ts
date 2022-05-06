// 用数组替代stack、queue
class Queue<T = number> {
  private readonly data: T[]
  private head = 0
  private tail = 0

  constructor(size: number) {
    this.data = Array(size).fill(undefined)
  }

  get length(): number {
    return this.tail - this.head
  }

  push(element: T): Queue<T> {
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

class RecentCounter {
  private readonly queue = new Queue(1e4 + 10)

  ping(time: number) {
    this.queue.push(time)
    while (this.queue.length && this.queue.at(0)! + 3000 < time) {
      this.queue.shift()
    }

    return this.queue.length
  }
}

const counter = new RecentCounter()
console.log(counter.ping(1))
console.log(counter.ping(2))
console.log(counter.ping(3001))
console.log(counter.ping(3002))

export {}
