import { MinHeap } from '../../../2_queue/minheap'

type Price = number
type Timestamp = number

interface Stock {
  price: Price
  timestamp: Timestamp
}

class StockPrice {
  private minQueue: MinHeap<Stock>
  private maxQueue: MinHeap<Stock>
  private record: Map<Timestamp, Price>
  private curTime: number

  constructor() {
    this.minQueue = new MinHeap((a, b) => a.price - b.price)
    this.maxQueue = new MinHeap((a, b) => b.price - a.price)
    this.record = new Map()
    this.curTime = 0
  }

  update(timestamp: number, price: number): void {
    if (timestamp >= this.curTime) this.curTime = timestamp
    this.record.set(timestamp, price)
    this.minQueue.push({ price, timestamp })
    this.maxQueue.push({ price, timestamp })
  }

  current(): number {
    return this.record.get(this.curTime) || -1
  }

  maximum(): number {
    let res = this.maxQueue.peek()
    while (res.price !== this.record.get(res.timestamp)) {
      this.maxQueue.shift()
      res = this.maxQueue.peek()
    }
    return res.price
  }

  minimum(): number {
    let res = this.minQueue.peek()
    while (res.price !== this.record.get(res.timestamp)) {
      this.minQueue.shift()
      res = this.minQueue.peek()
    }
    return res.price
  }
}
