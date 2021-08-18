// 环形缓冲器
// 在一个普通队列里，一旦一个队列满了，
// 我们就不能插入下一个元素，即使在队列前面仍有空间。
// 但是使用循环队列，
// 我们能使用这些空间去存储新的值。
class CircularDeque {
  private maxSize: number
  private data: Uint16Array
  private head: number
  private tail: number
  private size: number

  constructor({ maxSize }: { maxSize: number }) {
    this.maxSize = maxSize
    this.data = new Uint16Array(maxSize)
    this.head = 0
    this.tail = -1
    this.size = 0
  }

  // head前移
  insertFront(value: number): boolean {
    if (this.isFull()) return false
    this.head = (this.head - 1) % this.maxSize
    this.data[this.head] = value
    this.size++
    return true
  }

  // tail后移
  insertLast(value: number): boolean {
    if (this.isFull()) return false
    this.tail = (this.tail + 1) % this.maxSize
    this.data[this.tail] = value
    this.size++
    return true
  }

  // head后移
  deleteFront(): boolean {
    if (this.isEmpty()) return false
    this.head = (this.head + 1) % this.maxSize
    this.size--
    return true
  }

  // tail前移
  deleteLast(): boolean {
    if (this.isEmpty()) return false
    this.tail = (this.tail - 1) % this.maxSize
    this.size--
    return true
  }

  front(): number {
    return this.isEmpty() ? -1 : this.data[this.head]
  }

  rear(): number {
    return this.isEmpty() ? -1 : this.data[this.tail]
  }

  isEmpty(): boolean {
    return this.size === 0
  }

  isFull(): boolean {
    return this.size === this.maxSize
  }
}

const cq = new CircularDeque({ maxSize: 5 })

cq.insertLast(1)
cq.insertLast(1)
// cq.insertLast(2)
// cq.deleteFront()
// cq.deleteFront()
// cq.insertFront(2)

console.log(cq)

export default 1
