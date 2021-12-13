import { HashHeap } from '../../../8_heap/HashHeap'

type Start = number
type End = number

class Interval {
  start: number
  end: number
  interval: [Start, End]

  constructor(start: number, end: number) {
    this.start = start
    this.end = end
    this.interval = [start, end]
  }
}

class ExamRoom {
  /**
   * @description 根据线段长度从小到大存放所有线段
   */
  private treeSet: HashHeap<Interval> //
  /**
   * @description 将端点 p 映射到以 p 为左端点的线段
   */
  private startMap: Map<Start, Interval> //
  /**
   * @description 将端点 p 映射到以 p 为右端点的线段
   */
  private endMap: Map<End, Interval>
  private size: number
  constructor(n: number) {
    this.size = n
    // 如果长度相同，就比较索引
    this.treeSet = new HashHeap(
      (a, b) => this.getDistance(b) - this.getDistance(a) || a.start - b.start
    )
    this.startMap = new Map()
    this.endMap = new Map()
    // 虚拟线段
    this.addInterval(new Interval(-1, n))
  }

  /**
   * 来了一名考生，返回你给他分配的座位
   */
  seat(): number {
    const longest = this.treeSet.shift()!
    const { start, end } = longest
    console.log(longest)
    let seat: number
    if (start === -1) {
      seat = 0
    } else if (end === this.size) {
      seat = this.size - 1
    } else {
      seat = (start + end) >> 1
    }

    // 将最长的线段分成两段
    const left = new Interval(start, seat)
    const right = new Interval(seat, end)
    this.removeInterval(longest)
    this.addInterval(left)
    this.addInterval(right)

    return seat
  }

  /**
   *
   * @param p 坐在 p 位置的考生离开了 可以认为 p 位置一定坐有考生
   */
  leave(p: number): void {
    const left = this.endMap.get(p)!
    const right = this.startMap.get(p)!
    const merged = new Interval(left.start, right.end)
    this.removeInterval(left)
    this.removeInterval(right)
    this.addInterval(merged)
  }

  private addInterval(interval: Interval): void {
    this.treeSet.push(interval)
    this.startMap.set(interval.start, interval)
    this.endMap.set(interval.end, interval)
  }

  private removeInterval(interval: Interval): void {
    this.treeSet.remove(interval)
    this.startMap.delete(interval.start)
    this.endMap.delete(interval.end)
  }

  /**
   *
   * @param interval 计算线段长度
   * 由于 如果有多个这样的座位，安排到他到索引最小的那个座位
   * 所以不能简单地让它计算一个线段两个端点间的长度，而是让它计算该线段中点和端点之间的长度
   *  例如[0,4]和[4,9]应该选取前一段的中点2而不是后一段的6
   */
  private getDistance(interval: Interval): number {
    // 注意这里的思路
    if (interval.start === -1) return interval.end
    if (interval.end === this.size) return interval.end - 1 - interval.start
    return (interval.end - interval.start) >> 1
  }

  static main() {
    const examRoom = new ExamRoom(8)
    // console.dir(examRoom, { depth: null })

    console.log(examRoom.seat())
    // console.dir(examRoom, { depth: null })
    console.log(examRoom.seat())
    console.log(examRoom.seat())
    console.dir(examRoom, { depth: null })
    examRoom.leave(0)
    examRoom.leave(7)
    console.dir(examRoom, { depth: null })
    console.log(examRoom.seat())
    // console.log(examRoom.seat())
    // console.log(examRoom.seat())
    // console.log(examRoom.seat())
    // console.log(examRoom.seat())
  }
}

ExamRoom.main()

export {}
// 这道题可以理解为男生进男厕所如何挑选小便池子的问题

// 当学生进入考场后，他必须坐在能够使他与离他最近的人之间的距离达到最大化的座位上。
// 如果有多个这样的座位，安排到他到索引最小的那个座位
// 如果考场里没有人，那么学生就坐在 0 号座位上。
// https://labuladong.gitbook.io/algo/mu-lu-ye-4/zuo-wei-tiao-du

// 但凡遇到在动态过程中取最值的要求，肯定要使用有序数据结构，我们常用的数据结构就是二叉堆和平衡二叉搜索树了。
// 二叉堆实现的优先级队列取最值的时间复杂度是 O(logN)，但是只能删除最大值。
// 平衡二叉树也可以取最值，也可以修改、删除任意一个值，而且时间复杂度都是 O(logN)。
// 综上，二叉堆不能满足 leave 操作，应该使用平衡二叉树。
// 所以这里我们会用到 Java 的一种数据结构 TreeSet，这是一种有序数据结构，底层由红黑树维护有序性。
