// 使用主席树可以维护区间[L,R]中每个数值的次数，
// 然后，在这个root[R]-root[L-1]版本的线段树上，
// 查询节点的值大于等于threshold的节点的值。
class SegmentTreeNode {
  left!: SegmentTreeNode
  right!: SegmentTreeNode
  value = 0
}

// 值域上的线段树
class SegmentTree {
  private readonly root = new SegmentTreeNode()
  private readonly size: number

  constructor(size: number) {
    this.size = size
  }

  query(left: number, right: number, threshold: number): number {
    return this._query(left, right, 1, this.size, threshold)
  }

  insert(pre: number, cur: number, num: number): void {
    this._insert(pre, cur, 1, this.size, num)
  }

  private _query(L: number, R: number, l: number, r: number, threshold: number): number {}

  private _insert(L: number, R: number, l: number, r: number, num: number) {
    if (L <= l && r <= R) return
    const mid = Math.floor((l + r) / 2)
    if (num <= mid) this._insert(L, R, l, r, num)
    else {
    }
  }
}

class MajorityChecker {
  constructor(arr: number[]) {}

  query(left: number, right: number, threshold: number): number {}
}

/**
 * Your MajorityChecker object will be instantiated and called as such:
 * var obj = new MajorityChecker(arr)
 * var param_1 = obj.query(left,right,threshold)
 */
export {}
