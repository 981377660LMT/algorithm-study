import { MinHeap } from './minheap'

// 内部使用了堆
class PriorityQueue<Item = number> {
  private minHeap: MinHeap<Item>
  constructor(
    compareFunction: (a: Item, b: Item) => number = MinHeap.defaultCompareFunction,
    volumn: number = Infinity,
    heap: Item[] = []
  ) {
    this.minHeap = new MinHeap<Item>(compareFunction, volumn, heap)
  }

  static createPriorityQueue<Item = number>({
    compareFunction = MinHeap.defaultCompareFunction,
    volumn = Infinity,
    heap = [],
  }: {
    compareFunction?: (a: Item, b: Item) => number
    volumn?: number
    heap?: Item[]
  }) {
    return new PriorityQueue(compareFunction, volumn, heap)
  }

  get length() {
    return this.minHeap.size
  }

  // O(log(h))
  push(...val: Item[]) {
    val.forEach(v => this.minHeap.push(v))
    return this.length
  }

  // O(log(h))
  shift() {
    return this.minHeap.shift()
  }

  peek() {
    return this.minHeap.peek()
  }

  heapify() {
    this.minHeap.heapify()
  }
}

if (require.main === module) {
  const pq = new PriorityQueue()
  pq.push(5, 2, 3)
  // console.log(pq)
  console.log(pq.shift())
  // console.log(pq)
  console.log(pq.shift())
  // console.log(pq)
  console.log(pq.shift())
  console.log(pq.shift())
  console.log(pq.shift())
}

// 优先队列的优势:不需要一次性知道所有数据/数据流读取数据时/极大规模数据(1T)

export { PriorityQueue }
