// https://ei1333.github.io/library/structure/others/binary-indexed-tree.hpp

// API:
// Add(i, v): a[i] += v
// Query(r): sum(a[0..r))
// Query(l, r): sum(a[l..r))
// LowerBound(x): min(i) where sum(a[0..i]) >= x
// UpperBound(x): min(i) where sum(a[0..i]) > x

/**
 * 支持树上二分的树状数组.
 */
class FenwickTree {
  private readonly _n: number
  private readonly _log: number
  private readonly _data: number[]

  constructor(nOrNums: number | ArrayLike<number>) {
    if (typeof nOrNums === 'number') {
      this._n = nOrNums
      this._log = 32 - Math.clz32(this._n)
      this._data = Array(this._n + 1).fill(0)
    } else {
      this._n = nOrNums.length
      this._log = 32 - Math.clz32(this._n)
      this._data = Array(this._n + 1)
      this.build(nOrNums)
    }
  }

  add(i: number, v: number): void {
    for (i++; i <= this._n; i += i & -i) this._data[i] += v
  }

  // [0, r).
  query(r: number): number {
    if (r > this._n) r = this._n
    let res = 0
    for (; r > 0; r &= r - 1) res += this._data[r]
    return res
  }

  // [l, r).
  queryRange(l: number, r: number): number {
    return this.query(r) - this.query(l)
  }

  // [0,k]闭区间的和大于等于x的最小k.如果不存在,返回n.
  lowerBound(x: number): number {
    let i = 0
    for (let k = 1 << this._log; k > 0; k >>= 1) {
      if (i + k <= this._n && this._data[i + k] < x) {
        x -= this._data[i + k]
        i += k
      }
    }
    return i
  }

  // [0,k]闭区间的和大于x的最小k.如果不存在,返回n.
  upperBound(x: number): number {
    let i = 0
    for (let k = 1 << this._log; k > 0; k >>= 1) {
      if (i + k <= this._n && this._data[i + k] <= x) {
        x -= this._data[i + k]
        i += k
      }
    }
    return i
  }

  build(arr: ArrayLike<number>): void {
    if (this._n !== arr.length) throw new Error(`The length of arr must be equal to ${this._n}.`)
    for (let i = 1; i <= this._n; i++) this._data[i] = arr[i - 1]
    for (let i = 1; i <= this._n; i++) {
      const j = i + (i & -i)
      if (j <= this._n) this._data[j] += this._data[i]
    }
  }

  toString(): string {
    const sb = []
    for (let i = 0; i < this._n; i++) sb.push(this.queryRange(i, i + 1))
    return `FenwickTree: [${sb.join(', ')}]`
  }
}

export { FenwickTree }

if (require.main === module) {
  const tree = new FenwickTree(10)
  console.log(tree.toString())
  tree.add(0, 1)
  tree.add(1, 2)
  tree.add(9, 1)
  console.log(tree.toString())
  console.log(tree.lowerBound(3))
  console.log(tree.upperBound(3))
}
