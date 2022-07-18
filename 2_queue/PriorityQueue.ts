import { MinHeap } from '../8_heap/MinHeap'

// 内部使用了堆
class PriorityQueue<Item = number> {
  private minHeap: MinHeap<Item>

  static createPriorityQueue<Item = number>({
    comparator = MinHeap.defaultComparator,
    volumn = Infinity,
    heap = [],
  }: {
    comparator?: (a: Item, b: Item) => number
    volumn?: number
    heap?: Item[]
  }) {
    return new PriorityQueue(comparator, volumn, heap)
  }

  constructor(
    comparator: (a: Item, b: Item) => number = MinHeap.defaultComparator,
    volumn: number = Infinity,
    heap: Item[] = []
  ) {
    this.minHeap = new MinHeap<Item>(comparator, volumn, heap)
  }

  get length(): number {
    return this.minHeap.size
  }

  push(...elements: Item[]): number {
    elements.forEach(v => this.minHeap.heappush(v))
    return this.length
  }

  shift(): Item | undefined {
    return this.minHeap.heappop()
  }

  peek(): Item | undefined {
    return this.minHeap.peek()
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
