// // https://leetcode-cn.com/problems/fancy-sequence/solution/ru-guo-bu-liao-jie-cheng-fa-ni-yuan-jiu-hao-hao-xu/

import assert from 'assert'

const MOD = BigInt(1e9 + 7)
// 注意1e9+7*(1e9+7) 超出了js的精度了 9e15
// 结点实现超内存了 换数组

class SegmentTree {
  private readonly size: number
  private readonly tree: BigUint64Array
  private readonly lazyAdd: BigUint64Array
  private readonly lazyMul: BigUint64Array

  constructor(size: number) {
    this.size = size
    this.tree = new BigUint64Array(size << 2)
    this.lazyAdd = new BigUint64Array(size << 2)
    this.lazyMul = new BigUint64Array(size << 2)
  }

  query(l: bigint, r: bigint): number {
    return this._query(1, l, r, 1n, BigInt(this.size))
  }

  update(l: bigint, r: bigint, delta: bigint, type: 'ADD' | 'MUL'): void {
    this._update(1, l, r, 1n, BigInt(this.size), delta, type)
  }

  private _query(rt: number, L: bigint, R: bigint, l: bigint, r: bigint): number {
    if (L <= l && r <= R) return Number(this.tree[rt])

    const mid = (l + r) >> 1n
    this._pushDown(rt, l, r, mid)
    let res = 0

    if (L <= mid) {
      res += this._query(rt << 1, L, R, l, mid)
    }

    if (mid < R) {
      res += this._query((rt << 1) | 1, L, R, mid + 1n, r)
    }

    return res
  }

  private _update(
    rt: number,
    L: bigint,
    R: bigint,
    l: bigint,
    r: bigint,
    delta: bigint,
    type: 'ADD' | 'MUL'
  ): void {
    if (L <= l && r <= R) {
      if (type === 'ADD') {
        this.tree[rt] += (l - r + 1n) * delta
        this.lazyAdd[rt] += delta
        this.tree[rt] %= MOD
        this.lazyAdd[rt] %= MOD
      } else {
        this.tree[rt] *= delta
        this.lazyAdd[rt] *= delta
        this.lazyMul[rt] *= delta
        this.tree[rt] %= MOD
        this.lazyAdd[rt] %= MOD
        this.lazyMul[rt] %= MOD
      }

      return
    }

    const mid = (l + r) >> 1n
    this._pushDown(rt, l, r, mid)

    if (L <= mid) {
      this._update(rt << 1, L, R, l, mid, delta, type)
    }

    if (mid < R) {
      this._update((rt << 1) | 1, L, R, mid + 1n, r, delta, type)
    }

    this._pushUp(rt)
  }

  private _pushDown(rt: number, l: bigint, r: bigint, mid: bigint): void {
    if (this.lazyAdd[rt] !== 0n || this.lazyMul[rt] !== 1n) {
      this.tree[rt << 1] = this.tree[rt << 1] * this.lazyMul[rt] + this.lazyAdd[rt]
      this.lazyAdd[rt << 1] = this.lazyAdd[rt << 1] * this.lazyMul[rt] + this.lazyAdd[rt]
      this.lazyMul[rt << 1] *= this.lazyMul[rt]

      this.tree[(rt << 1) | 1] = this.tree[(rt << 1) | 1] * this.lazyMul[rt] + this.lazyAdd[rt]
      this.lazyAdd[(rt << 1) | 1] =
        this.lazyAdd[(rt << 1) | 1] * this.lazyMul[rt] + this.lazyAdd[rt]
      this.lazyMul[(rt << 1) | 1] *= this.lazyMul[rt]

      this.tree[rt << 1] %= MOD
      this.lazyAdd[rt << 1] %= MOD
      this.lazyMul[rt << 1] %= MOD
      this.tree[(rt << 1) | 1] %= MOD
      this.lazyAdd[(rt << 1) | 1] %= MOD
      this.lazyMul[(rt << 1) | 1] %= MOD

      this.lazyAdd[rt] = 0n
      this.lazyMul[rt] = 1n
    }
  }

  // 只有单点查询没必要pushUp
  private _pushUp(rt: number): void {
    // this.tree[rt] = this.tree[rt << 1] + this.tree[(rt << 1) | 1]
  }
}

class Fancy {
  private size: number
  private tree: SegmentTree

  /**
   * 总共最多会有 1e5 次对 append，addAll，multAll 和 getIndex 的调用。
   */
  constructor() {
    this.size = 0
    this.tree = new SegmentTree(1e5 + 10)
  }

  /**
   * @param val 将整数 val 添加在序列末尾
   */
  append(val: number): void {
    this.size += 1
    this.tree.update(BigInt(this.size), BigInt(this.size), BigInt(val), 'ADD')
  }

  /**
   * @param inc 将所有序列中的现有数值都增加 inc
   */
  addAll(inc: number): void {
    if (this.size === 0) return
    this.tree.update(1n, BigInt(this.size), BigInt(inc), 'ADD')
  }

  /**
   * @param m 将序列中的所有现有数值都乘以整数 m
   */
  multAll(m: number): void {
    if (this.size === 0) return
    this.tree.update(1n, BigInt(this.size), BigInt(m), 'MUL')
  }

  /**
   * @param idx 得到下标为 idx 处的数值,并将结果对 109 + 7 取余。
   * 如果下标大于等于序列的长度，请返回 -1
   */
  getIndex(idx: number): number {
    if (idx >= this.size) return -1
    return this.tree.query(BigInt(idx + 1), BigInt(idx + 1))
  }
}

// /**
//  * Your Fancy object will be instantiated and called as such:
//  * var obj = new Fancy()
//  * obj.append(val)
//  * obj.addAll(inc)
//  * obj.multAll(m)
//  * var param_4 = obj.getIndex(idx)
//  */

if (require.main === module) {
  const fancy = new Fancy()
  fancy.append(2)
  fancy.addAll(3)
  assert.strictEqual(fancy.getIndex(0), 5)

  fancy.append(7)
  assert.strictEqual(fancy.getIndex(1), 7) // 0

  fancy.multAll(2)
  assert.strictEqual(fancy.getIndex(0), 10)
  assert.strictEqual(fancy.getIndex(1), 14)
}

// export {}
