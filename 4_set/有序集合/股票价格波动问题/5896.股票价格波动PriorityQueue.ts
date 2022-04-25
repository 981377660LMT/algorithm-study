import { TreapMultiSet } from '../js/Treap'

type Time = number
type Price = number

// 很多数据结构设计题都可以用一个map+一个sortedList解决
// map维护数据的 id => value 的对应关系，sortedList 维护数据的有序性
class StockPrice {
  private curTime: Time
  private map: Map<Time, Price>
  private sortedList: TreapMultiSet<Price>

  constructor() {
    this.curTime = 0
    this.map = new Map()
    this.sortedList = new TreapMultiSet<Price>()
  }

  update(timestamp: number, price: number): void {
    this.curTime = Math.max(timestamp, this.curTime)
    this.map.has(timestamp) && this.sortedList.delete(this.map.get(timestamp)!)
    this.sortedList.add(price)
    this.map.set(timestamp, price)
  }

  current(): number {
    return this.map.get(this.curTime)!
  }

  maximum(): number {
    return this.sortedList.last()!
  }

  minimum(): number {
    return this.sortedList.first()!
  }
}
