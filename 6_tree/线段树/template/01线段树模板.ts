/* eslint-disable no-inner-declarations */

/**
 * 维护01序列的线段树 更新方式为`叠加`
 */
class SegmentTree1 {
  private readonly _tree: number[]
  private readonly _lazyValue: number[]
  private readonly _size: number

  /**
   * @param nums 01数组
   */
  constructor(nums: (0 | 1)[]) {
    this._size = nums.length
    this._tree = Array(this._size << 2).fill(0)
    this._lazyValue = Array(this._size << 2).fill(0)
    this._build(1, 1, this._size, nums)
  }

  /**
   * @param k 树上二分查询第k个1的位置 k>=1
   * @complexity O(logn)
   */
  getKthPos(k: number): number {
    return this._getKthPos(1, 1, this._size, k)
  }

  query(left: number, right: number): number {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return 0 // !超出范围返回0
    return this._query(1, left, right, 1, this._size)
  }

  update(left: number, right: number, delta: number): void {
    if (left < 1) left = 1
    if (right > this._size) right = this._size
    if (left > right) return
    this._update(1, left, right, 1, this._size, delta)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: (0 | 1)[]): void {
    if (l === r) {
      this._tree[rt] = nums[l - 1] // 维护01序列
      return
    }
    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _getKthPos(rt: number, l: number, r: number, k: number): number {
    if (l === r) return l
    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    const leftValue = this._tree[rt << 1]
    if (leftValue >= k) return this._getKthPos(rt << 1, l, mid, k)
    return this._getKthPos((rt << 1) | 1, mid + 1, r, k - leftValue)
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
}

/**
 * 维护01序列的线段树 更新方式为`染色` 使用isLazy数组
 */
class SegmentTree2 {
  private readonly _tree: Uint32Array
  private readonly _lazyValue: Uint8Array
  private readonly _isLazy: Uint8Array
  private readonly _size: number

  /**
   * @param nums 01数组
   */
  constructor(nums: (0 | 1)[]) {
    this._size = nums.length
    this._tree = new Uint32Array(this._size << 2)
    this._lazyValue = new Uint8Array(this._size << 2)
    this._isLazy = new Uint8Array(this._size << 2)
    this._build(1, 1, this._size, nums)
  }

  /**
   * @param k 树上二分查询第k个1的位置 k>=1
   * @complexity O(logn)
   */
  getKthPos(k: number): number {
    return this._getKthPos(1, 1, this._size, k)
  }

  query(l: number, r: number): number {
    if (l < 1) l = 1
    if (r > this._size) r = this._size
    if (l > r) return 0 // !超出范围返回0

    return this._query(1, l, r, 1, this._size)
  }

  update(l: number, r: number, target: 0 | 1): void {
    if (l < 1) l = 1
    if (r > this._size) r = this._size
    if (l > r) return

    this._update(1, l, r, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: (0 | 1)[]): void {
    if (l === r) {
      this._tree[rt] = nums[l - 1] // 维护01序列
      return
    }
    const mid = Math.floor((l + r) / 2)
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _getKthPos(rt: number, l: number, r: number, k: number): number {
    if (l === r) return l
    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    const leftValue = this._tree[rt << 1]
    if (leftValue >= k) return this._getKthPos(rt << 1, l, mid, k)
    return this._getKthPos((rt << 1) | 1, mid + 1, r, k - leftValue)
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

  private _update(rt: number, L: number, R: number, l: number, r: number, target: 0 | 1): void {
    if (L <= l && r <= R) {
      this._isLazy[rt] = 1
      this._lazyValue[rt] = target
      this._tree[rt] = target === 1 ? r - l + 1 : 0
      return
    }

    const mid = Math.floor((l + r) / 2)
    this._pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this._pushUp(rt)
  }

  private _pushUp(rt: number): void {
    this._tree[rt] = this._tree[rt << 1] + this._tree[(rt << 1) | 1]
  }

  private _pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this._isLazy[rt]) {
      const target = this._lazyValue[rt]
      this._lazyValue[rt << 1] = target
      this._lazyValue[(rt << 1) | 1] = target
      this._tree[rt << 1] = target === 1 ? mid - l + 1 : 0
      this._tree[(rt << 1) | 1] = target === 1 ? r - mid : 0
      this._isLazy[rt << 1] = 1
      this._isLazy[(rt << 1) | 1] = 1

      this._lazyValue[rt] = 0
      this._isLazy[rt] = 0
    }
  }
}

if (require.main === module) {
  // 身高体重重建队列
  function reconstructQueue(people: number[][]): number[][] {
    const n = people.length
    people.sort((a, b) => a[0] - b[0] || -(a[1] - b[1]))

    const tree = new SegmentTree1(Array(n).fill(1))
    const res = Array.from<unknown, [height: number, preCount: number]>({ length: n }, () => [0, 0])
    people.forEach(([height, preCount]) => {
      const pos = tree.getKthPos(preCount + 1)
      res[pos - 1] = [height, preCount]
      tree.update(pos, pos, -1) // !更新方式为叠加
    })

    return res
  }
}
export {}
