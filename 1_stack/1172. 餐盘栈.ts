import { BIT } from '../6_tree/树状数组/BIT'

class DinnerPlates {
  private bit: BIT
  private capacity: number
  private size: number
  private N: number
  private stack: number[][]
  // 我们把无限数量 ∞ 的栈排成一行，按从左到右的次序从 0 开始编号
  // 每个栈的的最大容量 capacity 都相同。
  constructor(capacity: number) {
    this.capacity = capacity
    this.size = 0
    this.N = 2 * 10 ** 5
    this.bit = new BIT(this.N)
    this.stack = Array.from({ length: this.N }, () => [])
  }

  // 将给出的正整数 val 推入 从左往右第一个 没有满的栈。
  push(val: number): void {
    const index = this.getPushIndex()
    this.stack[index].push(val)
    this.size++
    this.bit.add(index, 1)
  }

  // 返回 从右往左第一个 非空栈顶部的值，并将其从栈中删除；如果所有的栈都是空的，请返回 -1
  // 空的栈会占用他的位置
  pop(): number {
    // if (!this.size) return -1
    // const index = this.getPopIndex()
    // const top = this.stack[index].pop()!
    // this.size--
    // this.bit.add(index, -1)
    // return top
    const index = this.getPopIndex()
    return this.popAtStack(index - 1)
  }

  // 返回编号 index 的栈顶部的值，并将其从栈中删除；如果编号 index 的栈是空的，请返回 -1。
  popAtStack(index: number): number {
    const idx = index + 1
    if (!this.stack[idx].length) return -1
    const top = this.stack[idx].pop()!
    this.size--
    this.bit.add(idx, -1)
    return top
  }

  private getPopIndex(): number {
    let l = 1
    let r = this.N
    while (l <= r) {
      const mid = (l + r) >> 1
      // 等于时 右边都不是 要左移
      if (this.bit.query(mid) >= this.size) r = mid - 1
      else l = mid + 1
    }
    return l
  }

  private getPushIndex(): number {
    let l = 1
    let r = this.N
    while (l <= r) {
      const mid = (l + r) >> 1
      // 等于时 左边都不是 要右移
      if (this.bit.query(mid) < this.capacity * mid) r = mid - 1
      else l = mid + 1
    }
    return l
  }
}

const d = new DinnerPlates(2)
d.push(1)
d.push(2)
d.push(3)
d.push(4)
d.push(5)
console.log(d.popAtStack(0)) // 2
export {}

// 从右往左第一个栈未空的
// 从左往右第一个栈未满的
// 如果我们从 nums[] 出发，只能是用线性扫描来找到这个栈，显然对 nums[] 是不能用二分的
// 这时我们想到前缀和 sums[] 是具有单调性的，对于 sums[] 是可以用二分来找的
// 而可以动态维护前缀和的便是树状数组 tree[]，故我们可以只用到 tree[] 便实现了对 nums[] 和 sums[] 的管理
// 二分找未空栈
// 从右往左第一个非空栈，它的前缀和一定等于总元素个数
// 即满足 sums[] >= size
// 故我们只需要找到最小的一个值 x，使它满足 sums[x] >= size 即可
// 对于从左往右第一个非满栈，在它前面的栈一定都是满栈，则对这个非满栈来说
// 满足 sums[x] > x * capacity
// O(logn^2)
