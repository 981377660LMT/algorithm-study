class ArrayDeque<T = number> {
  private readonly capacity: number
  private readonly data: T[]
  private head: number
  private tail: number
  public length: number;

  *[Symbol.iterator]() {
    while (this.length) {
      yield this.shift()!
    }
  }

  constructor(capacity: number) {
    this.capacity = capacity
    this.data = []
    this.head = 0 // 从-1开始'向前'存
    this.tail = -1 // 从0开始向后存
    this.length = 0
  }

  // head前移
  unshift(value: T): boolean {
    if (this.isFull()) return false
    this.head = (this.head - 1 + this.capacity) % this.capacity
    this.data[this.head] = value
    this.length++
    return true
  }

  // tail后移
  push(value: T): boolean {
    if (this.isFull()) return false
    this.tail = (this.tail + 1 + this.capacity) % this.capacity
    this.data[this.tail] = value
    this.length++
    return true
  }

  // head后移
  shift(): T | undefined {
    if (this.isEmpty()) return undefined
    const front = this.front()
    this.head = (this.head + 1 + this.capacity) % this.capacity
    this.length--
    return front
  }

  // tail前移
  pop(): T | undefined {
    if (this.isEmpty()) return undefined
    const rear = this.rear()
    this.tail = (this.tail - 1 + this.capacity) % this.capacity
    this.length--
    return rear
  }

  forEach(callback: (value: T, index: number, array: T[]) => void): void {
    let head = this.head
    const times = this.length
    for (let i = 0; i < times; i++) {
      callback(this.data[head], i, this.data)
      head = (head + 1 + this.capacity) % this.capacity
    }
  }

  front(): T | undefined {
    return this.isEmpty() ? undefined : this.data[(this.head + this.capacity) % this.capacity]
  }

  // fix:deque 注意所有下标要加容量后取模
  rear(): T | undefined {
    return this.isEmpty() ? undefined : this.data[(this.tail + this.capacity) % this.capacity]
  }

  isEmpty(): boolean {
    return this.length === 0
  }

  isFull(): boolean {
    return this.length === this.capacity
  }
}

export { ArrayDeque }

if (require.main === module) {
  const deque = new ArrayDeque<number>(4)
  deque.push(1)
  deque.unshift(2)
  // console.log(deque)
  // console.log(deque.shift())
  // console.log(deque)
  // console.log(deque.pop())
  const arr = Array.from(deque)
}
