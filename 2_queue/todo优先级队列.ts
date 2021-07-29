import { MinHeap } from './minheap'

// 内部使用了堆
class PriotiryQueue<Item> {
  private minHeap: MinHeap<Item>
  constructor() {
    this.minHeap = new MinHeap()
  }

  get length() {
    return this.minHeap.size
  }

  // O(log(h))
  push(val: Item) {
    this.minHeap.push(val)
    return this
  }

  // O(log(h))
  shift() {
    return this.minHeap.shift()
  }
}

const pq = new PriotiryQueue()
pq.push(5).push(2).push(3)
console.log(pq)
console.log(pq.shift())
console.log(pq)
console.log(pq.shift())
console.log(pq)
console.log(pq.shift())
console.log(pq.shift())
console.log(pq.shift())
// 优先队列的优势:不需要一次性知道所有数据/数据流读取数据时/极大规模数据(1T)
