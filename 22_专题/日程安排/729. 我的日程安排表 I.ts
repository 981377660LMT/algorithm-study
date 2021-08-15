// 两个区间交叉:[a,b]与[c,d] 等价于a<d&&b>c
class MyCalendar {
  constructor(private record: [number, number][] = []) {}

  /**
   * @description 区间不重叠时才能插入
   * @summary O(n),splice方法
   */
  book(start: number, end: number): boolean {
    let l = 0
    let r = this.record.length - 1

    while (l <= r) {
      const mid = ~~((l + r) / 2)
      const [small, big] = this.record[mid]
      if (small < end && big > start) return false
      if (start >= big) {
        l = mid + 1
      } else {
        r = mid - 1
      }
    }

    this.record.splice(l, 0, [start, end])
    return true
  }
}

const myCalendar = new MyCalendar()
console.log(myCalendar.book(10, 20))
console.log(myCalendar.book(15, 20))
console.log(myCalendar.book(20, 30))

export {}
