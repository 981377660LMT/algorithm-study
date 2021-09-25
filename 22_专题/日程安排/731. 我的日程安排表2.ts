// 如果要添加的时间内不会导致三重预订时，则可以存储这个新的日程安排。
class MyCalendar {
  constructor(private count: Map<number, number> = new Map()) {}

  /**
   * @description
   * 我们采用遍历的方式,每次book的时候,start处++,代表时间推进到start的时候,多了一个任务,end处--,代表时间推进到end的时候任务完成一个
   * 我们使用一个 count map 来存储所有的预定，对于每次插入，我们执行count[start] += 1和count[end] -= 1。
   * 用JS的map每次遍历需要排序。如果你使用的是 Java，可以直接使用现成的数据结构 TreeMap。
   * 与HashMap相比，TreeMap是一个能比较元素大小的Map集合，会对传入的key进行了大小排序
   * TreeMap实现了红黑树的结构，形成了一颗二叉树。
   * @summary 此处为O(nlogn) 如果是treemap则为O(n)遍历复杂度
   */
  book(start: number, end: number): boolean {
    this.count.set(start, (this.count.get(start) || 0) + 1)
    this.count.set(end, (this.count.get(end) || 0) - 1)

    let tmp = 0
    let max = 0
    for (const [_, count] of [...this.count.entries()].sort((a, b) => a[0] - b[0])) {
      tmp += count
      max = Math.max(max, tmp)
      if (max >= 3) {
        this.count.set(start, this.count.get(start)! - 1)
        this.count.set(end, this.count.get(end)! + 1)
        return false
      }
    }

    return true
  }
}

const myCalendar = new MyCalendar()
console.log(myCalendar.book(10, 20))
console.log(myCalendar.book(50, 60))
console.log(myCalendar.book(10, 40))
console.log(myCalendar.book(5, 15)) // false
console.log(myCalendar.book(5, 10))
console.log(myCalendar.book(25, 55))

export {}
