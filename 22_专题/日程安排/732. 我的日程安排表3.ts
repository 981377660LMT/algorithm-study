// 如果要添加的时间内不会导致K重预订时，则可以存储这个新的日程安排。
class MyCalendar {
  private count: Map<number, number>
  private overlapCount: number

  constructor({
    count = new Map(),
    overlapCount = 1,
  }: {
    count?: Map<number, number>
    overlapCount?: number
  }) {
    this.count = count
    this.overlapCount = overlapCount
  }

  book(start: number, end: number): boolean {
    this.count.set(start, (this.count.get(start) || 0) + 1)
    this.count.set(end, (this.count.get(end) || 0) - 1)

    let tmp = 0
    let max = 0
    for (const [_, count] of [...this.count.entries()].sort((a, b) => a[0] - b[0])) {
      tmp += count
      max = Math.max(max, tmp)
      if (max >= this.overlapCount) {
        this.count.set(start, this.count.get(start)! - 1)
        this.count.set(end, this.count.get(end)! + 1)
        return false
      }
    }

    return true
  }
}

const myCalendar = new MyCalendar({ overlapCount: 3 })
console.log(myCalendar.book(10, 20))
console.log(myCalendar.book(50, 60))
console.log(myCalendar.book(10, 40))
console.log(myCalendar.book(5, 15)) // false
console.log(myCalendar.book(5, 10))
console.log(myCalendar.book(25, 55))

export {}
