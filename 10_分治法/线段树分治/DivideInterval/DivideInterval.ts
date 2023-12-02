/**
 * 线段树分割区间.
 * 将长度为n的序列搬到长度为offset+n的线段树上, 以实现快速的区间操作.
 */
class DivideInterval {
  /**
   * 线段树中一共有`offset+n`个节点.
   * `offset+i`对应原图的第i个顶点(0<=i<n).
   */
  readonly offset: number
  private readonly _n: number

  /** n 相同的线段树结构相同. */
  constructor(n: number) {
    let offset = 1
    while (offset < n) {
      offset <<= 1
    }
    this.offset = offset
    this._n = n
  }

  /**
   * 获取原下标为rawIndex的元素在树中的(叶子)编号.
   */
  id(rawIndex: number): number {
    if (!(rawIndex >= 0 && rawIndex < this._n)) throw new Error('invalid index')
    return this.offset + rawIndex
  }

  /**
   * O(logn) 顺序遍历`[start,end)`区间对应的线段树节点编号.
   *
   * ```ts
   * const D = new DivideInterval(13)
   * D.enumerateSegment(0, 1, i => console.log(i))  // 16
   * D.enumerateSegment(0, 3, i => console.log(i))  // 8 18
   * ```
   */
  enumerateSegment(start: number, end: number, f: (segmentId: number) => void): void {
    if (!(start >= 0 && start <= end && end <= this._n)) {
      throw new Error('invalid range')
    }
    const ids = this._getSegmentIds(start, end)
    for (let i = 0; i < ids.length; i++) {
      f(ids[i])
    }
  }

  isLeaf(segmentId: number): boolean {
    return segmentId >= this.offset
  }

  /** 线段树结点个数. */
  get size(): number {
    return this.offset + this._n
  }

  /**
   * O(n) 从根向叶子方向push.
   */
  pushDown(f: (parent: number, child: number) => void): void {
    for (let p = 1; p < this.offset; p++) {
      f(p, p << 1)
      f(p, (p << 1) | 1)
    }
  }

  /**
   * O(n) 从叶子向根方向update.
   */
  pushUp(f: (parent: number, child1: number, child2: number) => void): void {
    for (let p = this.offset - 1; p > 0; p--) {
      f(p, p << 1, (p << 1) | 1)
    }
  }

  /** 线段树中两个节点的最近公共祖先(两个二进制数字的最长公共前缀). */
  static lca(u: number, v: number): number {
    if (u === v) return u
    if (u > v) {
      const tmp = u
      u = v
      v = tmp
    }
    const depth1 = this.depth(u)
    const depth2 = this.depth(v)
    const diff = u ^ (v >>> (depth2 - depth1))
    if (diff === 0) return u
    const len = 32 - Math.clz32(diff)
    return u >>> len
  }

  /** 线段树中节点的深度. */
  static depth(u: number): number {
    if (u === 0) return 0
    return 31 - Math.clz32(u)
  }

  private _getSegmentIds(start: number, end: number): number[] {
    if (!(start >= 0 && start <= end && end <= this._n)) {
      throw new Error('invalid range')
    }
    const leftRes: number[] = []
    const rightRes: number[] = []
    for (start += this.offset, end += this.offset; start < end; start >>>= 1, end >>>= 1) {
      if ((start & 1) === 1) {
        leftRes.push(start)
        start++
      }
      if ((end & 1) === 1) {
        end--
        rightRes.push(end)
      }
    }
    for (let i = rightRes.length - 1; i >= 0; i--) {
      leftRes.push(rightRes[i])
    }
    return leftRes
  }
}

export { DivideInterval }

if (require.main === module) {
  const D = new DivideInterval(13)
  D.enumerateSegment(0, 3, i => console.log(i))
}
