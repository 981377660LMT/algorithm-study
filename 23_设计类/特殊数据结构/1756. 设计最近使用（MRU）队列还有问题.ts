class BIT {
  public size: number
  private tree: number[]

  constructor(size: number) {
    this.size = size
    this.tree = Array(size + 1).fill(0)
  }

  add(x: number, k: number) {
    if (x <= 0) throw Error('查询索引应为正整数')
    for (let i = x; i <= this.size; i += this.lowbit(i)) {
      this.tree[i] += k
    }
  }

  query(x: number) {
    let res = 0
    for (let i = x; i > 0; i -= this.lowbit(i)) {
      res += this.tree[i]
    }
    return res
  }

  private lowbit(x: number) {
    return x & -x
  }
}

// 利用树状数组维护前缀和，二分查找第K个数字
class MRUQueue {
  private bit: BIT
  private data: number[]
  private tail: number

  constructor(private n: number) {
    this.data = Array.from({ length: n }, (_, i) => i + 1)
    this.bit = new BIT(n + 1)
    for (let i = 1; i <= n; i++) {
      this.bit.add(i, 1)
    }
    this.tail = n + 1
  }

  // 将第 k 个元素（从 1 开始索引）移到队尾，并返回该元素。
  fetch(k: number): number {
    let l = 1
    let r = this.tail
    while (l <= r) {
      const mid = (l + r) >> 1
      if (this.bit.query(mid) < k) l = mid + 1
      else r = mid - 1
    }
    this.data[this.tail] = this.data[l]
    this.data[l] = 0
    this.bit.add(l, -1)
    this.bit.add(this.tail, 1)
    this.tail++
    return this.data[this.tail - 1]
  }
}

const queue = new MRUQueue(8)
console.log(queue)

console.log(queue.fetch(3))
console.log(queue)
console.log(queue.fetch(5))
console.log(queue.fetch(2))
console.log(queue.fetch(8))
// 设计一种类似队列的数据结构，该数据结构将最近使用的元素移到队列尾部。
export {}
