/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */
/* eslint-disable no-inner-declarations */
// 珂朵莉树(ODT)/Intervals
// !noneValue不使用symbol,而是自定义的哨兵值,更加灵活.

const INF = 2e15

/**
 * 珂朵莉树，基于数据随机的颜色段均摊。
 * `VanEmdeBoasTree`实现.
 */
class ODTVan<S> {
  private readonly _noneValue: S
  private _len = 0
  private _count = 0
  private readonly _data: Map<number, S> = new Map()
  private readonly _leftLimit = -INF
  private readonly _rightLimit = INF
  private readonly _fs: VanEmdeBoasTree = new VanEmdeBoasTree()

  /**
   * 指定哨兵值建立一个ODT.初始时,所有位置的值为 {@link noneValue}.
   * @param noneValue 表示空值的哨兵值.
   */
  constructor(noneValue: S) {
    this._noneValue = noneValue
  }

  /**
   * 返回包含`x`的区间的信息.
   */
  get(x: number, erase = false): [start: number, end: number, value: S] | undefined {
    const start = this._fs.prev(x)
    const end = this._fs.next(x + 1)
    const value = this._getOrNone(start)
    if (erase && value !== this._noneValue) {
      this._len--
      this._count -= end - start
      this._data.set(start, this._noneValue)
      this._mergeAt(start)
      this._mergeAt(end)
    }
    return [start, end, value]
  }

  set(start: number, end: number, value: S): void {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    this.enumerateRange(start, end, () => {}, true) // remove
    this._fs.insert(start)
    this._data.set(start, value)
    if (value !== this._noneValue) {
      this._len++
      this._count += end - start
    }
    this._mergeAt(start)
    this._mergeAt(end)
  }

  enumerateAll(f: (start: number, end: number, value: S) => void): void {
    this.enumerateRange(this._leftLimit, this._rightLimit, f, false)
  }

  /**
   * 遍历范围`[start, end)`内的所有区间.
   */
  enumerateRange(
    start: number,
    end: number,
    f: (start: number, end: number, value: S) => void,
    erase = false
  ): void {
    if (start < this._leftLimit) start = this._leftLimit
    if (end > this._rightLimit) end = this._rightLimit
    if (start >= end) return

    const none = this._noneValue
    if (!erase) {
      let left = this._fs.prev(start)
      while (left < end) {
        const right = this._fs.next(left + 1)
        f(Math.max(left, start), Math.min(right, end), this._getOrNone(left))
        left = right
      }
      return
    }

    let p = this._fs.prev(start)
    if (p < start) {
      this._fs.insert(start)
      const v = this._getOrNone(p)
      this._data.set(start, v)
      if (v !== none) this._len++
    }

    p = this._fs.next(end)
    if (end < p) {
      const v = this._getOrNone(this._fs.prev(end))
      this._data.set(end, v)
      this._fs.insert(end)
      if (v !== none) this._len++
    }

    p = start
    while (p < end) {
      const q = this._fs.next(p + 1)
      const x = this._getOrNone(p)
      f(p, q, x)
      if (x !== none) {
        this._len--
        this._count -= q - p
      }
      this._fs.erase(p)
      p = q
    }

    this._fs.insert(start)
    this._data.set(start, none)
  }

  toString(): string {
    const sb: string[] = [`ODT(${this.length}) {`]
    this.enumerateAll((start, end, value) => {
      const v = value === this._noneValue ? 'null' : value
      sb.push(`  [${start},${end}):${v}`)
    })
    sb.push('}')
    return sb.join('\n')
  }

  /**
   * 区间个数.
   */
  get length(): number {
    return this._len
  }

  /**
   * 区间内元素个数之和.
   */
  get count(): number {
    return this._count
  }

  private _mergeAt(p: number): void {
    if (p <= 0 || this._rightLimit <= p) return
    const q = this._fs.prev(p - 1)
    const dataP = this._getOrNone(p)
    const dataQ = this._getOrNone(q)
    if (dataP === dataQ) {
      if (dataP !== this._noneValue) this._len--
      this._fs.erase(p)
    }
  }

  private _getOrNone(x: number): S {
    const res = this._data.get(x)
    return res === void 0 ? this._noneValue : res
  }
}

/**
 * Van Tree.
 * 梵峨眉大悲寺树.
 */
class VanEmdeBoasTree {
  private readonly _root: _VNode
  private _size = 0

  /**
   * @param depth 树的深度.默认为16.一般取16或32.
   */
  constructor(depth = 16) {
    this._root = new _VNode(depth)
  }

  has(x: number): boolean {
    return this._root.has(x)
  }

  insert(x: number): boolean {
    if (this.has(x)) return false
    this._size++
    this._root.insert(x)
    return true
  }

  erase(x: number): boolean {
    if (!this.has(x)) return false
    this._size--
    this._root.erase(x)
    return true
  }

  /**
   * 返回小于等于i的最大元素.如果不存在,返回-INF.
   */
  prev(x: number): number {
    return this._root.prev(x)
  }

  /**
   * 返回大于等于i的最小元素.如果不存在,返回INF.
   */
  next(x: number): number {
    return this._root.next(x)
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x >= end) break
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(-INF, INF, v => sb.push(v.toString()))
    return `VanEmdeBoasTree(${this.size}){${sb.join(', ')}}`
  }

  /**
   * 如果没有元素,返回INF.
   */
  get min(): number {
    return this._root.min
  }

  /**
   * 如果没有元素,返回-INF.
   */
  get max(): number {
    return this._root.max
  }

  get size(): number {
    return this._size
  }
}

class _VNode {
  dep: number
  min = INF
  max = -INF
  aux: _VNode | undefined = undefined
  son: Map<number, _VNode> = new Map()

  constructor(dep: number) {
    this.dep = dep
  }

  has(x: number): boolean {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x === vMin || x === vMax) return true
    if (!vDep || x < vMin || x > vMax) return false
    const i = x >>> vDep
    const soni = this.son.get(i)
    if (!soni) return false
    return soni.has(x - (i << vDep))
  }

  insert(x: number): void {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (vMin > vMax) {
      this.min = x
      this.max = x
      return
    }
    if (vMin === vMax) {
      if (x < vMin) {
        this.min = x
        return
      }
      if (x > vMax) {
        this.max = x
        return
      }
    }
    if (x < vMin) {
      const tmp = x
      x = vMin
      this.min = tmp
    }
    if (x > vMax) {
      const tmp = x
      x = vMax
      this.max = tmp
    }
    const i = x >>> vDep
    let soni = this.son.get(i)
    if (!soni) {
      soni = new _VNode(vDep >>> 1)
      this.son.set(i, soni)
    }
    if (soni.empty()) {
      if (!this.aux) this.aux = new _VNode(vDep >>> 1)
      this.aux.insert(i)
    }
    soni.insert(x - (i << vDep))
  }

  erase(x: number): void {
    const { min: vMin, max: vMax, dep: vDep, aux: vAux } = this
    if (vMin === x && vMax === x) {
      this.min = INF
      this.max = -INF
      return
    }
    if (x === vMin) {
      if (!vAux || vAux.empty()) {
        this.min = vMax
        return
      }
      const auxMin = vAux.min
      x = (auxMin << vDep) + this.son.get(auxMin)!.min
      this.min = x
    }
    if (x === vMax) {
      if (!vAux || vAux.empty()) {
        this.max = vMin
        return
      }
      const auxMax = vAux.max
      x = (auxMax << vDep) + this.son.get(auxMax)!.max
      this.max = x
    }
    const i = x >>> vDep
    const soni = this.son.get(i)!
    soni.erase(x - (i << vDep))
    if (soni.empty()) vAux!.erase(i)
  }

  prev(x: number): number {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x < vMin) return -INF
    if (x >= vMax) return vMax
    const i = x >>> vDep
    const hi = i << vDep
    const lo = x - hi
    const soni = this.son.get(i)
    if (soni && lo >= soni.min) return hi + soni.prev(lo)
    let y = -INF
    if (this.aux && i > 0) y = this.aux.prev(i - 1)
    if (y === -INF) return vMin
    return (y << vDep) + this.son.get(y)!.max
  }

  next(x: number): number {
    const { min: vMin, max: vMax, dep: vDep } = this
    if (x <= vMin) return vMin
    if (x > vMax) return INF
    const i = x >>> vDep
    const hi = i << vDep
    const lo = x - hi
    const soni = this.son.get(i)
    if (soni && lo <= soni.max) return hi + soni.next(lo)
    let y = INF
    if (this.aux) y = this.aux.next(i + 1)
    if (y === INF) return vMax
    return (y << vDep) + this.son.get(y)!.min
  }

  empty(): boolean {
    return this.min > this.max
  }
}

export { ODTVan }

if (require.main === module) {
  const INF = 2e15
  const van = new ODTVan(INF)
  console.log(van.toString())
  van.set(0, 10, 1)
  van.set(2, 5, 2)
  console.log(van.get(8))
  console.log(van.toString())
  van.enumerateRange(
    1,
    7,
    (start, end, value) => {
      console.log(start, end, value)
    },
    true
  )
  console.log(van.toString(), van.length)

  // 352. 将数据流变为多个不相交区间
  // https://leetcode.cn/problems/data-stream-as-disjoint-intervals/
  class SummaryRanges {
    private readonly _odt = new ODTVan(-1)

    addNum(value: number): void {
      this._odt.set(value, value + 1, 0)
    }

    getIntervals(): number[][] {
      const res: number[][] = []
      this._odt.enumerateAll((start, end, value) => {
        if (value === 0) res.push([start, end - 1])
      })
      return res
    }
  }
}
