type TimeStamp = number
type Count = number

class HitCounter {
  private record: Map<TimeStamp, Count>
  private curCount: number

  constructor() {
    this.record = new Map()
    this.curCount = 0
  }

  // 最早的时间戳从1开始，且都是按照时间顺序对系统进行调用
  hit(timestamp: number): void {
    this.record.set(timestamp, (this.record.get(timestamp) || 0) + 1)
    this.curCount++
  }

  getHits(timestamp: number): number {
    for (const [preTime, count] of this.record.entries()) {
      if (preTime + 300 <= timestamp) {
        this.curCount -= count
        this.record.delete(preTime)
      } else {
        break
      }
    }

    return this.curCount
  }
}

export {}
// 设计一个敲击计数器，使它可以统计在过去5分钟(300秒)内被敲击次数。
// 如果每秒的敲击次数是一个很大的数字，你的计数器可以应对吗

// LRU？类似于1797. 设计一个验证系统
// 这类题的关键是所有函数调用中，currentTime 的值 严格递增
