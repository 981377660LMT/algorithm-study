/* eslint-disable no-param-reassign */
// !单点修改+区间查询

/**
 * S 线段树维护的值的类型
 *
 * F 更新操作的值的类型/懒标记的值的类型
 *
 * @see {@link https://atcoder.jp/contests/typical90/submissions/23947318}
 */
abstract class AbstractSegmentTree<S, F> {
  /**
   * {@link data} 的初始值
   */
  protected abstract e(): S

  /**
   * 左右结点的合并方式
   */
  protected abstract op(data1: S, data2: S): S

  /**
   * 数组值的更新方式
   * @param f 更新操作的值
   * @param data 线段树维护的值
   */
  protected abstract mapping(f: F, data: S): S

  private readonly _n: number
  private readonly _data: S[]
  private _log: number

  constructor(sizeOrArray: number | S[]) {
    this._n = typeof sizeOrArray === 'number' ? sizeOrArray : sizeOrArray.length
    this._log = 32 - Math.clz32(this._n - 1)
    while (this._log < this._n) {
      this._log *= 2
    }

    const size = this._log * 2
    this._data = Array(size)
    this._data[0] = this.e()

    if (Array.isArray(sizeOrArray)) {
      for (let i = 0; i < this._n; i++) {
        this._data[this._log + i] = sizeOrArray[i]
      }

      for (let i = this._n; i < this._log; i++) {
        this._data[this._log + i] = this.e()
      }

      for (let i = this._log - 1; i >= 1; i--) {
        this._data[i] = this.op(this._data[i * 2], this._data[i * 2 + 1])
      }
    } else {
      for (let i = 1; i < size; i++) {
        this._data[i] = this.e()
      }
    }
  }

  /**
   * @param index 数组下标 0 <= index < n
   * @param value 更新的值 更新方式取决于 {@link MAPPING}
   */
  update(index: number, value: F): void {
    let k = index + this._log
    this._data[k] = this.mapping(value, this._data[k])
    while (k > 1) {
      k >>= 1
      this._data[k] = this.op(this._data[2 * k], this._data[2 * k + 1])
    }
  }

  /**
   * @returns  切片`[left:right]`内的信息 0 <= left <= right <=n
   */
  query(left: number, right: number): S {
    let res = this.e()
    for (left += this._log, right += this._log; left < right; left >>= 1, right >>= 1) {
      if (left & 1) {
        res = this.op(res, this._data[left++])
      }

      if (right & 1) {
        res = this.op(res, this._data[--right])
      }
    }

    return res
  }

  get length() {
    return this._n
  }
}

const INF = 2e15

/**
 * @description 线段树RMQ最大值(快速版)
 */
class MaxSegmentTree2 extends AbstractSegmentTree<number, number> {
  protected e(): number {
    return -INF
  }

  protected op(data1: number, data2: number): number {
    return Math.max(data1, data2)
  }

  protected mapping(f: number, data: number): number {
    return data + f
    return data
  }
}

/**
 * @description 线段树RMQ最小值(快速版)
 */
class MinSegmentTree2 extends AbstractSegmentTree<number, number> {
  protected e(): number {
    return INF
  }

  protected op(data1: number, data2: number): number {
    return Math.min(data1, data2)
  }

  protected mapping(f: number, data: number): number {
    return data + f
    return data
  }
}

if (require.main === module) {
  const tree = new MaxSegmentTree2([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(tree.length)
  console.log(tree.query(0, 1))
  tree.update(0, 3)
  console.log(tree.query(0, 1))

  const n = 2e5
  const arr = Array(n)
    .fill(0)
    .map((_, i) => i)
  const tree2 = new MaxSegmentTree2(arr)
  console.time('query')
  for (let i = 0; i < n; i++) {
    tree2.query(0, n)
  }
  console.timeEnd('query')
}

export {}
