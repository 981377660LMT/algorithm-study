// 求数组交集
// 关键思想是存储大数组 遍历小数组

type Index = number
type Value = number

class SparseVector {
  private noneZero: Map<Index, Value>

  constructor(private nums: number[]) {
    this.noneZero = new Map()
    nums.forEach((value, index) => {
      value && this.noneZero.set(index, value)
    })
  }

  // Return the dotProduct of two sparse vectors
  dotProduct(vec: SparseVector): number {
    return this.noneZero.size > vec.noneZero.size ? this.count(vec, this) : this.count(this, vec)
  }

  private count(small: SparseVector, big: SparseVector): number {
    let res = 0
    for (const [index, value] of small.noneZero.entries()) {
      if (big.noneZero.has(index)) res += value * big.noneZero.get(index)!
    }

    return res
  }
}

export {}
