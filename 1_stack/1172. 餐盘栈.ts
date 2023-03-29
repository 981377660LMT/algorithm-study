/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 1172. 餐盘栈

class DinnerPlates {
  private static readonly _N = 2e5
  private readonly _stack: number[][]
  private readonly _stackCap: number
  private readonly _bit: BIT
  private _size = 0

  // 我们把无限数量 ∞ 的栈排成一行，按从左到右的次序从 0 开始编号
  // 每个栈的的最大容量 capacity 都相同。
  constructor(capacity: number) {
    this._stackCap = capacity
    this._bit = new BIT(DinnerPlates._N)
    this._stack = Array(DinnerPlates._N)
      .fill(0)
      .map(() => [])
  }

  // 将给出的正整数 val 推入 从左往右第一个 没有满的栈。
  push(val: number): void {
    const index = this._getPushIndex()
    this._stack[index].push(val)
    this._size++
    this._bit.add(index, 1)
  }

  // 返回 从右往左第一个 非空栈顶部的值，并将其从栈中删除；如果所有的栈都是空的，请返回 -1
  // 空的栈会占用他的位置
  pop(): number {
    const index = this._getPopIndex()
    return this.popAtStack(index - 1)
  }

  // 返回编号 index 的栈顶部的值，并将其从栈中删除；如果编号 index 的栈是空的，请返回 -1。
  popAtStack(index: number): number {
    index++
    if (!this._stack[index].length) return -1
    const top = this._stack[index].pop()!
    this._size--
    this._bit.add(index, -1)
    return top
  }

  private _getPopIndex(): number {
    let left = 1
    let right = DinnerPlates._N
    while (left <= right) {
      const mid = (left + right) >> 1
      // 等于时 右边都不是 要左移
      if (this._bit.query(mid) >= this._size) right = mid - 1
      else left = mid + 1
    }
    return left
  }

  private _getPushIndex(): number {
    let left = 1
    let right = DinnerPlates._N
    while (left <= right) {
      const mid = (left + right) >> 1
      // 等于时 左边都不是 要右移
      if (this._bit.query(mid) < this._stackCap * mid) right = mid - 1
      else left = mid + 1
    }
    return left
  }
}

class BIT {
  private readonly _size: number
  private readonly _tree: Array<number>

  constructor(size: number) {
    this._size = size
    this._tree = Array(size + 1).fill(0)
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw RangeError(`add索引 ${index} 应为正整数`)
    for (let i = index; i <= this._size; i += i & -i) {
      this._tree[i] += delta
    }
  }

  query(index: number): number {
    if (index > this._size) index = this._size
    let res = 0
    for (let i = index; i > 0; i -= i & -i) {
      res += this._tree[i]
    }
    return res
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
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
