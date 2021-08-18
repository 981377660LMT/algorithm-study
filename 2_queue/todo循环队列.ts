// 环形缓冲器
// 在一个普通队列里，一旦一个队列满了，
// 我们就不能插入下一个元素，即使在队列前面仍有空间。
// 但是使用循环队列，
// 我们能使用这些空间去存储新的值。

// 循环队列中的元素是[head,tail]这一段,[0,-1]表示没有元素size为0
// 入队tail++ 出队head++
class CircularQueue {
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

  // 入队tail后移
  enQueue(value: number): boolean {
    if (this.isFull()) return false
    this.tail = (this.tail + 1) % this.maxSize
    this.data[this.tail] = value
    this.size++
    return true
  }

  // 出队head前移
  deQueue(): boolean {
    if (this.isEmpty()) return false
    if (this.head === this.tail) {
      this.head = 0
      this.tail = -1
    } else {
      this.head = (this.head + 1) % this.maxSize
    }
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

const cq = new CircularQueue({ maxSize: 5 })

cq.enQueue(1)
cq.enQueue(1)
cq.enQueue(1)
cq.enQueue(1)
cq.enQueue(1)
cq.enQueue(2)
cq.enQueue(1)
cq.enQueue(1)
cq.deQueue()
cq.deQueue()
cq.deQueue()
cq.deQueue()
cq.enQueue(2)
cq.deQueue()
cq.enQueue(5)
cq.deQueue()
cq.enQueue(44)
cq.enQueue(43)
cq.enQueue(43)
cq.enQueue(43)

console.log(cq)

export default 1
