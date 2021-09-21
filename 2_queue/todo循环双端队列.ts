// 环形缓冲器
// 在一个普通队列里，一旦一个队列满了，
// 我们就不能插入下一个元素，即使在队列前面仍有空间。
// 但是使用循环队列，
// 我们能使用这些空间去存储新的值。
class MyCircularDeque {
  private maxSize: number
  private data: Array<number>
  private head: number
  private tail: number
  private size: number

  constructor(maxSize: number) {
    this.maxSize = maxSize
    this.data = []
    this.head = 0 // 从-1开始向前存
    this.tail = -1 // 从0开始向后存
    this.size = 0
  }

  // head前移
  insertFront(value: number): boolean {
    if (this.isFull()) return false
    this.head = (this.head - 1 + this.maxSize) % this.maxSize
    this.data[this.head] = value
    this.size++
    return true
  }

  // tail后移
  insertLast(value: number): boolean {
    if (this.isFull()) return false
    this.tail = (this.tail + 1 + this.maxSize) % this.maxSize
    this.data[this.tail] = value
    this.size++
    return true
  }

  // head后移
  deleteFront(): boolean {
    if (this.isEmpty()) return false
    this.head = (this.head + 1 + this.maxSize) % this.maxSize
    this.size--
    return true
  }

  // tail前移
  deleteLast(): boolean {
    if (this.isEmpty()) return false
    this.tail = (this.tail - 1 + this.maxSize) % this.maxSize
    this.size--
    return true
  }

  getFront(): number {
    return this.isEmpty() ? -1 : this.data[(this.head + this.maxSize) % this.maxSize]
  }

  getRear(): number {
    return this.isEmpty() ? -1 : this.data[(this.tail + this.maxSize) % this.maxSize]
  }

  isEmpty(): boolean {
    return this.size === 0
  }

  isFull(): boolean {
    return this.size === this.maxSize
  }
}

const cq = new MyCircularDeque(3)

console.log(cq.insertFront(9))
console.log(cq)
console.log(cq.getRear())
console.log(cq.insertFront(9))

// console.log(cq)

export default 1
