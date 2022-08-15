// 406. 根据身高重建队列-线段树树上二分 O(nlogn)

function reconstructQueue(people: number[][]): number[][] {
  const n = people.length
  people.sort((a, b) => a[0] - b[0] || -(a[1] - b[1]))

  const tree = new SegmentTree(n)
  const res = Array.from<unknown, [height: number, preCount: number]>({ length: n }, () => [0, 0])
  people.forEach(([height, preCount]) => {
    const pos = tree.getPos(preCount + 1)
    res[pos - 1] = [height, preCount]
    tree.update(pos, pos, -1)
  })

  return res
}

class SegmentTree {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this._size = size
    this._tree = Array(size << 2).fill(0)
    this._lazyValue = Array(size << 2).fill(0)
    this._build(1, 1, size)
  }

  /**
   * @param k 树上二分查询第k个1的位置  k>=1
   * @complexity O(logn)
   */
  getPos(k: number): number {
    return this._getPos(1, 1, this._size, k)
  }

  query(l: number, r: number): number {
    // this.checkRange(l, r)
    return this._query(1, l, r, 1, this._size)
  }

  update(l: number, r: number, delta: number): void {
    // this.checkRange(l, r)
    this._update(1, l, r, 1, this._size, delta)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number): void {
    if (l === r) {
      this._tree[rt] = 1 // 维护01序列
      return
    }
    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid)
    this._build((rt << 1) | 1, mid + 1, r)
    this._pushUp(rt)
  }

  private _getPos(rt: number, l: number, r: number, k: number): number {
    if (l === r) return l
    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    const leftValue = this._tree[rt << 1]
    if (leftValue >= k) return this._getPos(rt << 1, l, mid, k)
    return this._getPos((rt << 1) | 1, mid + 1, r, k - leftValue)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
    if (L <= l && r <= R) {
      this._lazyValue[rt] += delta
      this._tree[rt] += delta * (r - l + 1)
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, delta)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, delta)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = this._tree[rt << 1] + this._tree[(rt << 1) | 1]
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._lazyValue[rt]) {
      const delta = this._lazyValue[rt]
      this._lazyValue[rt << 1] += delta
      this._lazyValue[(rt << 1) | 1] += delta
      this._tree[rt << 1] += delta * (mid - l + 1)
      this._tree[(rt << 1) | 1] += delta * (r - mid)
      this._lazyValue[rt] = 0
    }
  }

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this._size) {
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this._size}]`)
    }
  }
}

if (require.main === module) {
  console.log(
    reconstructQueue([
      [7, 0],
      [4, 4],
      [7, 1],
      [5, 0],
      [6, 1],
      [5, 2]
    ])
  )
}

export {}
