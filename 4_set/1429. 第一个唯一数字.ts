class FirstUnique {
  private only: Set<number> // 现在出现频率为1的元素
  private dead: Set<number> // 永久删除

  constructor(nums: number[]) {
    this.only = new Set()
    this.dead = new Set()

    for (const num of nums) {
      this.add(num)
    }
  }

  // 返回队列中的 第一个唯一 整数的值。如果没有唯一整数，返回 -1
  showFirstUnique(): number {
    if (!this.only.size) return -1
    const first = this.only.keys().next().value
    return first
  }

  add(value: number): void {
    if (this.dead.has(value)) return
    if (this.only.has(value)) {
      this.only.delete(value)
      this.dead.add(value)
    } else this.only.add(value)
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
