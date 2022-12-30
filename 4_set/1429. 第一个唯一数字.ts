class FirstUnique {
  private readonly _one = new Set<number>() // 现在出现频率为1的元素
  private readonly _removed = new Set<number>() // 永久删除

  constructor(nums: number[]) {
    nums.forEach(num => this.add(num))
  }

  // 返回队列中的 第一个唯一 整数的值。如果没有唯一整数，返回 -1
  showFirstUnique(): number {
    if (!this._one.size) return -1
    const first = this._one.keys().next().value
    return first
  }

  add(value: number): void {
    if (this._removed.has(value)) return
    if (this._one.has(value)) {
      this._one.delete(value)
      this._removed.add(value)
      return
    }
    this._one.add(value)
  }
}

const fu = new FirstUnique([2, 3, 5])
console.log(fu.showFirstUnique())
console.log(fu.add(5))
console.log(fu.showFirstUnique())
console.log(fu.add(2))
console.log(fu.showFirstUnique())
console.log(fu.add(3))
console.log(fu.showFirstUnique())

export {}
// JS的set是LinkeHashdSet

// 给定一系列整数，插入一个队列中，找出队列中第一个唯一整数

// 没思路的时候尽量往哈希堆的思路去想
