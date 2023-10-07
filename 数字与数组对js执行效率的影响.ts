// https://leetcode.cn/problems/online-stock-span/description/?envType=daily-question&envId=2023-10-07
// https://leetcode.cn/problems/online-stock-span/submissions/472094704/?envType=daily-question&envId=2023-10-07

// !1.INF对js执行效率的影响(数字大小对js执行效率的影响 - int32)
// 同一份代码，令 INF = 2**31 需要 8500ms，
// 令 INF = 2**31 - 1 需要 3500ms
// 正好是int32与int64的区别
// 这是因为v8对数字的优化方案:...
// !2.数组与对象的性能差异(对象快于数组)
// INF = 2**31 - 1 时
// type E = [min: number, max: number]  // 6500ms, 138MB
// type E = { min: number; max: number } // 3500ms, 108MB
// !3.number与对象的性能差异
// type E = { min: number; max: number } // 3500ms, 108MB
// type E = number // 800ms, 82MB
// !4.将普通数组number[]换成Float64Array
// number[] // 800ms, 82MB
// Float64Array // 550ms, 70MB

// 得出优化结论:
// 1.尽量使用number而不是对象，尽量使用对象存pair而不是数组.
// 2.如果数组只存整数,不超过uint32使用Uint32Array，超过uint32使用Float64Array.
// 3.如果数据范围在int32内，尽量使用2e9作为INF.

const INF = 2e9
class StockSpanner {
  private readonly _Q: RightMostLeftMostQuerySegmentTree
  private _ptr = 0

  constructor() {
    this._Q = new RightMostLeftMostQuerySegmentTree(Array(1e5 + 10).fill(0))
  }

  next(price: number): number {
    const pos = this._ptr++
    this._Q.set(pos, price)
    const higherPos = this._Q.leftNearestHigher(pos)
    return higherPos === -1 ? pos + 1 : pos - higherPos
  }
}

type E = number
type Id = number

class RightMostLeftMostQuerySegmentTree {
  private readonly _n: number
  private readonly _rangeAddRangeMinMax: _SegmentTreeRangeUpdateRangeQuery<E, Id>

  constructor(arr: ArrayLike<number>) {
    this._n = arr.length
    const leaves: E[] = Array(this._n)
    for (let i = 0; i < this._n; i++) leaves[i] = arr[i]
    this._rangeAddRangeMinMax = new _SegmentTreeRangeUpdateRangeQuery<E, Id>(leaves, {
      e: () => 0,
      id: () => 0,
      op: (a, b) => Math.max(a, b),
      mapping: (f, x) => f + x,
      composition: (f, g) => f + g
    })
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    this._rangeAddRangeMinMax.set(index, value)
  }

  /**
   * 查询`index`左侧最近的下标`j`，使得 `nums[j] > nums[index]`.
   * 如果不存在，返回`-1`.
   */
  leftNearestHigher(index: number): number {
    const cur = this._rangeAddRangeMinMax.get(index)
    const cand = this._rangeAddRangeMinMax.minLeft(index, e => e <= cur) - 1
    return cand === -1 ? -1 : cand
  }
}

class _SegmentTreeRangeUpdateRangeQuery<E = number, Id = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _data: E[]
  private readonly _lazy: Id[]
  private readonly _e: () => E
  private readonly _id: () => Id
  private readonly _op: (a: E, b: E) => E
  private readonly _mapping: (id: Id, data: E) => E
  private readonly _composition: (id1: Id, id2: Id) => Id

  /**
   * 区间修改区间查询的懒标记线段树.维护幺半群.
   * @param nOrLeaves 大小或叶子节点的值.
   * @param operations 线段树的操作.
   */
  constructor(
    nOrLeaves: number | ArrayLike<E>,
    operations: {
      e: () => E
      id: () => Id
      op: (e1: E, e2: E) => E
      mapping: (lazy: Id, data: E) => E
      composition: (f: Id, g: Id) => Id
      equalsId?: (id1: Id, id2: Id) => boolean
    } & ThisType<void>
  ) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    const { e, id, op, mapping, composition } = operations

    let height = 32 - Math.clz32(n - 1)
    let size = 1 << height

    const data = new Float32Array(size << 1) as any
    for (let i = 0; i < data.length; i++) data[i] = e()
    const lazy = new Float32Array(size) as any
    for (let i = 0; i < lazy.length; i++) lazy[i] = 0
    this._n = n
    this._size = size
    this._height = height
    this._data = data
    this._lazy = lazy
    this._e = e
    this._id = id
    this._op = op
    this._mapping = mapping
    this._composition = composition

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    this._data[index] = value
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return this._data[index]
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (value: E) => boolean): number {
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let res = this._e()
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      if (!predicate(this._op(this._data[end], res))) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          if (predicate(this._op(this._data[end], res))) {
            res = this._op(this._data[end], res)
            end--
          }
        }
        return end + 1 - this._size
      }
      res = this._op(this._data[end], res)
      if ((end & -end) === end) break
    }
    return 0
  }

  build(leaves: ArrayLike<E>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) this._data[this._size + i] = leaves[i]
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  private _pushUp(index: number): void {
    this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (lazy === 0) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = 0 as Id
  }

  private _propagate(index: number, lazy: Id): void {
    this._data[index] = this._mapping(lazy, this._data[index])
    if (index < this._size) this._lazy[index] = this._composition(lazy, this._lazy[index])
  }
}
