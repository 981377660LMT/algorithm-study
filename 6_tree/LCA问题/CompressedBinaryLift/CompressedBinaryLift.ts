interface MutableArrayLike<T> {
  readonly length: number
  [n: number]: T
}

/**
 * 空间复杂度`O(n)`的树上倍增.
 *
 * @see
 * - https://taodaling.github.io/blog/2020/03/18/binary-lifting/
 * - https://codeforces.com/blog/entry/74847
 * - https://codeforces.com/blog/entry/100826
 */
class CompressedBinaryLift {
  readonly depth: MutableArrayLike<number>
  readonly parent: MutableArrayLike<number>
  private readonly _jump: MutableArrayLike<number>

  constructor(
    n: number,
    depthOnTree: MutableArrayLike<number>,
    parentOnTree: MutableArrayLike<number>
  )
  constructor(tree: MutableArrayLike<MutableArrayLike<number>>, root?: number)
  constructor(arg0: any, arg1?: any, arg2?: any) {
    if (arguments.length === 3) {
      const n = arg0
      this.depth = arg1
      this.parent = arg2
      this._jump = new Int32Array(n).fill(-1)
      for (let i = 0; i < n; i++) this._consider(i)
    } else {
      const n = arg0.length
      if (arg1 == undefined) arg1 = 0
      this.depth = new Int32Array(n)
      this.parent = new Int32Array(n)
      this.parent[arg1] = -1
      this._jump = new Int32Array(n)
      this._jump[arg1] = arg1
      this._setUp(arg0, arg1)
    }
  }

  firstTrue = (start: number, predicate: (end: number) => boolean): number => {
    while (!predicate(start)) {
      if (predicate(this._jump[start])) {
        start = this.parent[start]
      } else {
        if (start === this._jump[start]) return -1
        start = this._jump[start]
      }
    }
    return start
  }

  lastTrue = (start: number, predicate: (end: number) => boolean): number => {
    if (!predicate(start)) return -1
    while (true) {
      if (predicate(this._jump[start])) {
        if (start === this._jump[start]) return start
        start = this._jump[start]
      } else if (predicate(this.parent[start])) {
        start = this.parent[start]
      } else {
        return start
      }
    }
  }

  upToDepth = (root: number, toDepth: number): number => {
    if (!(0 <= toDepth && toDepth <= this.depth[root])) return -1
    while (this.depth[root] > toDepth) {
      if (this.depth[this._jump[root]] < toDepth) {
        root = this.parent[root]
      } else {
        root = this._jump[root]
      }
    }
    return root
  }

  kthAncestor = (node: number, k: number): number => {
    const targetDepth = this.depth[node] - k
    return this.upToDepth(node, targetDepth)
  }

  lca = (a: number, b: number): number => {
    if (this.depth[a] > this.depth[b]) {
      a = this.kthAncestor(a, this.depth[a] - this.depth[b])
    } else if (this.depth[a] < this.depth[b]) {
      b = this.kthAncestor(b, this.depth[b] - this.depth[a])
    }
    while (a !== b) {
      if (this._jump[a] === this._jump[b]) {
        a = this.parent[a]
        b = this.parent[b]
      } else {
        a = this._jump[a]
        b = this._jump[b]
      }
    }
    return a
  }

  dist = (a: number, b: number): number => {
    return this.depth[a] + this.depth[b] - 2 * this.depth[this.lca(a, b)]
  }

  jump = (start: number, target: number, step: number): number => {
    const lca = this.lca(start, target)
    const dep1 = this.depth[start]
    const dep2 = this.depth[target]
    const deplca = this.depth[lca]
    const dist = dep1 + dep2 - 2 * deplca
    if (step > dist) return -1
    if (step <= dep1 - deplca) return this.kthAncestor(start, step)
    return this.kthAncestor(target, dist - step)
  }

  private _consider = (root: number): void => {
    if (root === -1 || this._jump[root] !== -1) return
    const p = this.parent[root]
    this._consider(p)
    this._addLeaf(root, p)
  }

  private _addLeaf = (leaf: number, parent: number): void => {
    if (parent == -1) {
      this._jump[leaf] = leaf
    } else {
      const tmp = this._jump[parent]
      if (this.depth[parent] - this.depth[tmp] === this.depth[tmp] - this.depth[this._jump[tmp]]) {
        this._jump[leaf] = this._jump[tmp]
      } else {
        this._jump[leaf] = parent
      }
    }
  }

  private _setUp = (tree: MutableArrayLike<MutableArrayLike<number>>, root: number): void => {
    const queue: number[] = [root]
    let head = 0
    while (head < queue.length) {
      const cur = queue[head++]
      const nexts = tree[cur]
      for (let i = 0; i < nexts.length; i++) {
        const next = nexts[i]
        if (next === this.parent[cur]) continue
        this.depth[next] = this.depth[cur] + 1
        this.parent[next] = cur
        queue.push(next)
        this._addLeaf(next, cur)
      }
    }
  }
}

export { CompressedBinaryLift }

if (require.main === module) {
  const tree = [[1, 2], [0, 3, 4], [0], [1], [1]]
  const cbl = new CompressedBinaryLift(tree)
  console.log(cbl.lca(3, 1)) // 0

  // https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/
  class TreeAncestor {
    private readonly bl: CompressedBinaryLift

    constructor(n: number, parent: number[]) {
      const adjList: number[][] = Array(n)
      for (let i = 0; i < n; i++) adjList[i] = []
      parent.forEach((p, i) => {
        if (p !== -1) adjList[p].push(i)
      })
      this.bl = new CompressedBinaryLift(adjList)
    }

    getKthAncestor(node: number, k: number): number {
      return this.bl.kthAncestor(node, k)
    }
  }
}
