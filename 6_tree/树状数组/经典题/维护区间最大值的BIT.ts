const INF = 2e15

// TODO
// !Not Verified

/**
 * 维护最大值的BIT.
 * !1-indexed.
 */
class BITMax {
  private readonly _tree: number[]
  private readonly _max: number[]
  private readonly _n: number

  constructor(n: number) {
    this._tree = Array(n + 1).fill(0)
    this._max = Array(n + 1).fill(0)
    this._n = n
  }

  // O(log^2 n)
  set(index: number, value: number): void {
    this._max[index] = value
    for (; index <= this._n; index += index & -index) {
      this._tree[index] = value
      for (let i = 1; i < (index & -index); i <<= 1) {
        this._tree[index] = Math.max(this._tree[index], this._tree[index - i])
      }
    }
  }

  // O(log^2 n)
  queryMax(left: number, right: number): number {
    let x = right
    let res = -INF
    while (x >= left) {
      if (x - (x & -x) >= left - 1) {
        res = Math.max(res, this._tree[x])
        x -= x & -x
      } else {
        res = Math.max(res, this._max[x])
        x--
      }
    }
    return res
  }
}

export {}

if (require.main === module) {
  const bit = new BITMax(10)

  console.log(bit.toString())
  bit.set(1, 3)
  bit.set(2, 4)
  bit.set(3, 5)
  console.log(bit.toString())

  console.log(bit.toString())
  console.log(bit.queryMax(1, 3))
  console.log(bit.toString())
}
