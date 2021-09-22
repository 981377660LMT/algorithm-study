class Deque<T> {
  private capacity: number
  private data: T[]
  private head: number
  private tail: number
  private size: number

  constructor(capacity: number) {
    this.capacity = capacity
    this.data = []
    this.head = 0 // 从-1开始'向前'存
    this.tail = -1 // 从0开始向后存
    this.size = 0
  }

  // head前移
  unshift(value: T): boolean {
    if (this.isFull()) return false
    this.head = (this.head - 1 + this.capacity) % this.capacity
    this.data[this.head] = value
    this.size++
    return true
  }

  // tail后移
  push(value: T): boolean {
    if (this.isFull()) return false
    this.tail = (this.tail + 1 + this.capacity) % this.capacity
    this.data[this.tail] = value
    this.size++
    return true
  }

  // head后移
  shift(): T | undefined {
    if (this.isEmpty()) return undefined
    const front = this.front()
    this.head = (this.head + 1 + this.capacity) % this.capacity
    this.size--
    return front
  }

  // tail前移
  pop(): T | undefined {
    if (this.isEmpty()) return undefined
    const rear = this.rear()
    this.tail = (this.tail - 1 + this.capacity) % this.capacity
    this.size--
    return rear
  }

  forEach(callback: (value: T, index: number, array: T[]) => void): void {
    let head = this.head
    const times = this.size
    for (let i = 0; i < times; i++) {
      callback(this.data[head], i, this.data)
      head = (head + 1 + this.capacity) % this.capacity
    }
  }

  private front(): T | undefined {
    return this.isEmpty() ? undefined : this.data[this.head]
  }

  private rear(): T | undefined {
    return this.isEmpty() ? undefined : this.data[this.tail]
  }

  private isEmpty(): boolean {
    return this.size === 0
  }

  private isFull(): boolean {
    return this.size === this.capacity
  }
}

export { Deque }

if (require.main === module) {
  const deque = new Deque<number>(4)
  deque.push(1)
  deque.unshift(2)
  console.log(deque)
  console.log(deque.shift())
  console.log(deque)
  console.log(deque.pop())
}
