import { MinHeap } from './minheap'

// 内部使用了堆
class PriotiryQueue {
  private minHeap: MinHeap
  constructor() {
    this.minHeap = new MinHeap([])
  }

  get size() {
    return this.minHeap.size
  }

  // O(log(h))
  push(val: number) {
    this.minHeap.insert(val)
    return this
  }

  // O(log(h))
  shift() {
    this.minHeap.pop()
    return this
  }
}

const pq = new PriotiryQueue()

// 优先队列的优势:不需要一次性知道所有数据/数据流读取数据时/极大规模数据(1T)
