// http://www.manongjc.com/detail/28-aalqlatkoigmgre.html

// !一个1~n的排列，进行q次区间升/降序排序操作，问x最终所在的位置。
// !我们要找的这个位置的值满足>=x且不满足>=x+1。
// !将排列变为两个只含0/1的序列，最终满足一个为0一个为1的位置就是答案。
// !区间排序就很容易了，我们只在乎0/1，不在乎取值，0和1各自放一起，直接用线段树染色+查询维护

import * as fs from 'fs'
import { resolve } from 'path'

function useInput(path?: string) {
  let data: string
  if (path) {
    data = fs.readFileSync(resolve(__dirname, path), 'utf8')
  } else {
    data = fs.readFileSync(process.stdin.fd, 'utf8')
  }

  const lines = data.split(/\r\n|\r|\n/)
  let lineId = 0
  const input = (): string => lines[lineId++]

  return {
    input
  }
}

const { input } = useInput()

/**
 * 维护01序列的线段树 更新方式为染色
 */
class SegmentTree2 {
  private readonly _size: number
  private readonly _tree: Uint32Array
  private readonly _lazyValue: Uint8Array
  private readonly _isLazy: Uint8Array

  /**
   * @param nums 01数组
   */
  constructor(nums: ArrayLike<0 | 1>) {
    this._size = nums.length
    const cap = 1 << (32 - Math.clz32(this._size - 1) + 1)
    this._tree = new Uint32Array(cap)
    this._lazyValue = new Uint8Array(cap)
    this._isLazy = new Uint8Array(cap)
    this._build(1, 1, this._size, nums)
  }

  query(l: number, r: number): number {
    return this._query(1, l, r, 1, this._size)
  }

  update(l: number, r: number, target: 0 | 1): void {
    this._update(1, l, r, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
  }

  private _build(rt: number, l: number, r: number, nums: ArrayLike<0 | 1>): void {
    if (l === r) {
      this._tree[rt] = nums[l - 1]
      return
    }
    const mid = (l + r) >>> 1
    this._build(rt << 1, l, mid, nums)
    this._build((rt << 1) | 1, mid + 1, r, nums)
    this._pushUp(rt)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this._tree[rt]

    const mid = (l + r) >>> 1
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

    const mid = (l + r) >>> 1
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

const [n, q, x] = input().split(' ').map(Number)
const nums = input().split(' ').map(Number)
const Q: [type: number, left: number, right: number][] = []
for (let i = 0; i < q; i++) {
  const [type, left, right] = input().split(' ').map(Number)
  Q.push([type, left, right])
}

function getOrder(arr: number[]): number[] {
  const tree = new SegmentTree2(arr as ArrayLike<0 | 1>)
  for (let i = 0; i < q; i++) {
    const [type, left, right] = Q[i]
    const count = tree.query(left, right)

    if (type === 1) {
      // 升序排列，左边全变成0，右边全变成1
      if (right - count + 1 <= right) {
        tree.update(right - count + 1, right, 1)
      }

      if (left <= right - count) {
        tree.update(left, right - count, 0)
      }
    } else {
      // 降序排列，左边全变成1，右边全变成0
      if (left <= left + count - 1) {
        tree.update(left, left + count - 1, 1)
      }

      if (left + count <= right) {
        tree.update(left + count, right, 0)
      }
    }
  }

  const res: number[] = []
  for (let i = 1; i <= n; i++) {
    res.push(tree.query(i, i))
  }
  return res
}

const order1 = getOrder(nums.map(v => (v < x ? 0 : 1)))
const order2 = getOrder(nums.map(v => (v < x + 1 ? 0 : 1)))
for (let i = 0; i < n; i++) {
  if (order1[i] !== order2[i]) {
    // 不一样的位置即为答案
    console.log(i + 1)
    break
  }
}

export {}
