import { deserializeNode } from '../../6_tree/力扣加加/构建类/297二叉树的序列化与反序列化'

class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
  }
}

class SegmentTree {
  private tree: number[]
  private lazy: number[]
  private n: number

  constructor(n: number) {
    this.n = n
    this.tree = Array(4 * n).fill(0)
    this.lazy = Array(4 * n).fill(0)
  }

  query(l: number, r: number) {
    return this._query(1, l, r, 1, this.n)
  }

  update(l: number, r: number, delta: number) {
    this._update(1, l, r, 1, this.n, delta)
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = (l + r) >> 1
    this._pushDown(rt, l, r, mid)
    let res = 0

    if (L <= mid) {
      res += this._query(rt << 1, L, R, l, mid)
    }

    if (mid < R) {
      res += this._query((rt << 1) | 1, L, R, mid + 1, r)
    }

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, delta: number) {
    if (L <= l && r <= R) {
      this.lazy[rt] = delta
      this.tree[rt] = delta === 1 ? r - l + 1 : 0
      return
    }

    const mid = (l + r) >> 1
    this._pushDown(rt, l, r, mid)

    if (L <= mid) {
      this._update(rt << 1, L, R, l, mid, delta)
    }

    if (mid < R) {
      this._update((rt << 1) | 1, L, R, mid + 1, r, delta)
    }

    this._pushUp(rt)
  }

  private _pushUp(rt: number) {
    this.tree[rt] = this.tree[rt << 1] + this.tree[(rt << 1) | 1]
  }

  private _pushDown(rt: number, l: number, r: number, mid: number) {
    if (this.lazy[rt]) {
      const delta = this.lazy[rt]
      this.lazy[rt << 1] = delta
      this.lazy[(rt << 1) | 1] = delta
      this.tree[rt << 1] = delta === 1 ? mid - l + 1 : 0
      this.tree[(rt << 1) | 1] = delta === 1 ? r - mid : 0
      this.lazy[rt] = 0
    }
  }
}

function getNumber(root: TreeNode | null, ops: number[][]): number {
  function dfs(root: TreeNode | null) {
    if (!root) return
    values.add(root.val)
    allSet.add(root.val)
    dfs(root.left)
    dfs(root.right)
  }

  const values = new Set<number>()
  const allSet = new Set<number>()
  dfs(root)

  for (const [_, x, y] of ops) {
    allSet.add(x)
    allSet.add(y)
  }

  const allNums = [...allSet].sort((a, b) => a - b)
  const mapping = new Map<number, number>()
  for (const [index, num] of allNums.entries()) {
    mapping.set(num, index + 1)
  }

  const sg = new SegmentTree(allNums.length + 10)
  for (const [opt, x, y] of ops) {
    if (opt === 0) {
      sg.update(mapping.get(x)!, mapping.get(y)!, -1)
    } else {
      sg.update(mapping.get(x)!, mapping.get(y)!, 1)
    }
  }

  let res = 0
  for (let v of values) {
    res += sg.query(mapping.get(v)!, mapping.get(v)!)
  }

  return res
}

const root = deserializeNode([1, null, 2, null, 3, null, 4, null, 5])
console.log(
  getNumber(root, [
    [1, 2, 4],
    [1, 1, 3],
    [0, 3, 5],
  ])
)
export {}
