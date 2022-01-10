import assert from 'assert'

/**
 * @description 循环数组
 */
class ArrayDeque<T = number> {
  private readonly capacity: number
  private readonly data: T[]
  private head: number
  private tail: number
  public length: number

  constructor(capacity: number) {
    this.capacity = capacity
    this.data = []
    this.head = 0 // 从-1开始'向前'存
    this.tail = -1 // 从0开始向后存
    this.length = 0
  }

  *[Symbol.iterator]() {
    let head = this.head
    const times = this.length
    for (let i = 0; i < times; i++) {
      yield this.data[head]
      head = (head + 1 + this.capacity) % this.capacity
    }
  }

  at(index: number): T | undefined {
    if (index < 0) index += this.length
    const pos = (this.head + index + this.capacity) % this.capacity
    return this.data[pos]
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

  // fix:deque 注意所有下标要加容量后取模
  private front(): T | undefined {
    return this.isEmpty() ? undefined : this.data[(this.head + this.capacity) % this.capacity]
  }

  private rear(): T | undefined {
    return this.isEmpty() ? undefined : this.data[(this.tail + this.capacity) % this.capacity]
  }

  private isEmpty(): boolean {
    return this.length === 0
  }

  private isFull(): boolean {
    return this.length === this.capacity
  }
}

export { ArrayDeque }

if (require.main === module) {
  const deque = new ArrayDeque<number>(4)

  deque.push(1)
  deque.unshift(2)

  for (const num of deque) {
    console.log(num)
  }

  assert.strictEqual(deque.at(0), 2)
  assert.strictEqual(deque.at(1), 1)
  assert.strictEqual(deque.at(2), void 0)
  assert.strictEqual(deque.at(-1), 1)
  assert.strictEqual(deque.at(-2), 2)
  assert.strictEqual(deque.at(-3), void 0)
}
