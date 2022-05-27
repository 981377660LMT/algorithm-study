// ops[i] = [type, x, y] 表示第 i 次操作为：
// type 等于 0 时，将节点值范围在 [x, y] 的节点均染蓝
// type 等于 1 时，将节点值范围在 [x, y] 的节点均染红
// 请返回完成所有染色后，该二叉树中红色节点的数量。

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
  private readonly tree: Uint32Array
  private readonly lazyValue: Uint8Array
  private readonly isLazy: Uint8Array
  private readonly size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this.size = size
    this.tree = new Uint32Array(size << 2)
    this.lazyValue = new Uint8Array(size << 2)
    this.isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    this.checkRange(l, r)
    return this._query(1, l, r, 1, this.size)
  }

  update(l: number, r: number, target: 0 | 1): void {
    this.checkRange(l, r)
    this._update(1, l, r, 1, this.size, target)
  }

  queryAll(): number {
    return this.tree[1]
  }

  private _query(rt: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) return this.tree[rt]

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    let res = 0
    if (L <= mid) res += this._query(rt << 1, L, R, l, mid)
    if (mid < R) res += this._query((rt << 1) | 1, L, R, mid + 1, r)

    return res
  }

  private _update(rt: number, L: number, R: number, l: number, r: number, target: 0 | 1): void {
    if (L <= l && r <= R) {
      this.isLazy[rt] = 1
      this.lazyValue[rt] = target
      this.tree[rt] = target === 1 ? r - l + 1 : 0
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.pushDown(rt, l, r, mid)
    if (L <= mid) this._update(rt << 1, L, R, l, mid, target)
    if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(rt)
  }

  private pushUp(rt: number): void {
    this.tree[rt] = this.tree[rt << 1] + this.tree[(rt << 1) | 1]
  }

  private pushDown(rt: number, l: number, r: number, mid: number): void {
    if (this.isLazy[rt]) {
      const target = this.lazyValue[rt]
      this.lazyValue[rt << 1] = target
      this.lazyValue[(rt << 1) | 1] = target
      this.tree[rt << 1] = target === 1 ? mid - l + 1 : 0
      this.tree[(rt << 1) | 1] = target === 1 ? r - mid : 0
      this.isLazy[rt << 1] = 1
      this.isLazy[(rt << 1) | 1] = 1

      this.lazyValue[rt] = 0
      this.isLazy[rt] = 0
    }
  }

  private checkRange(l: number, r: number): void {
    if (l < 1 || r > this.size) throw new RangeError(`[${l}, ${r}] out of range: [1, ${this.size}]`)
  }
}

function getNumber(root: TreeNode | null, ops: number[][]): number {
  const treeValues = new Set<number>()
  const allValues = new Set<number>()
  dfs(root)

  for (const [_, x, y] of ops) {
    allValues.add(x)
    allValues.add(y)
  }

  const allNums = [...allValues].sort((a, b) => a - b)
  const mapping = new Map<number, number>()
  for (const [index, num] of allNums.entries()) {
    mapping.set(num, index + 1)
  }

  const sg = new SegmentTree(allNums.length + 10)
  for (const [opt, x, y] of ops) {
    if (opt === 0) {
      sg.update(mapping.get(x)!, mapping.get(y)!, 0)
    } else {
      sg.update(mapping.get(x)!, mapping.get(y)!, 1)
    }
  }

  let res = 0
  for (let v of treeValues) {
    res += sg.query(mapping.get(v)!, mapping.get(v)!)
  }

  return res

  function dfs(root: TreeNode | null): void {
    if (!root) return
    treeValues.add(root.val)
    allValues.add(root.val)
    dfs(root.left)
    dfs(root.right)
  }
}

export {}
