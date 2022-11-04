/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */
// # 有 1 ~ n 编号的 n 个软件包，除了 1 号软件包，每个软件包有且只有一个它所依赖的软件包，

// # 且软件包之间的依赖关系不存在环。
// # 每安装一个软件，都必须先安装它所依赖的软件包；
// # 每卸载一个软件，都必须把所有依赖于它的软件包都先卸载掉。

// # 初始时所有软件包都是未安装状态，然后进行 m 次操作，每次操作需要安装或者在一个软件包，
// # 并输出这次操作安装状态被改变的软件包数量。

// # 每次对编号为 x 的节点操作
// # 安装操作需要把根节点到 x 路径上的点全部置为 1
// # 卸载操作需要把以 x 为根的子树中的点全部置为 0
// #region

class SegmentTree {
  private readonly _tree: Uint32Array
  private readonly _lazyValue: Uint8Array
  private readonly _isLazy: Uint8Array
  private readonly _size: number

  /**
   *
   * @param size 区间右边界
   */
  constructor(size: number) {
    this._size = size
    this._tree = new Uint32Array(size << 2)
    this._lazyValue = new Uint8Array(size << 2)
    this._isLazy = new Uint8Array(size << 2)
  }

  query(l: number, r: number): number {
    this._checkRange(l, r)
    return this._query(1, l, r, 1, this._size)
  }

  update(l: number, r: number, target: number): void {
    this._checkRange(l, r)
    this._update(1, l, r, 1, this._size, target)
  }

  queryAll(): number {
    return this._tree[1]
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

  private _update(rt: number, L: number, R: number, l: number, r: number, target: number): void {
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

  private _checkRange(l: number, r: number): void {
    if (l < 1 || r > this._size) {
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this._size}]`)
    }
  }
}
// #endregion

// n个结点
function useSoftwareManager(n: number, adjList: number[][], root: number) {
  const tree = new SegmentTree(n + 10)
  const depths = new Uint32Array(n + 10)
  const parents = new Int32Array(n + 10) // 不存在为-1
  const subSizes = new Uint32Array(n + 10).fill(1)
  const heavySons = new Int32Array(n + 10) // 每个点的重儿子，不存在时为 -1
  const tops = new Uint32Array(n + 10) // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
  const dfsIds = new Uint32Array(n + 10)
  const dfsIdToNode = new Uint32Array(n + 10) // dfs序在树中所对应的节点 dfsIdToNode[dfsIds[i]] = i
  let dfsId = 1

  dfs1(root, -1, 0)
  dfs2(root, root)

  // 寻找重儿子，拆分出重链与轻链
  function dfs1(cur: number, pre: number, depth: number): void {
    let heavySon = -1
    let heavySonSize = 0
    for (let i = 0; i < adjList[cur].length; i++) {
      const next = adjList[cur][i]
      if (next === pre) continue
      dfs1(next, cur, depth + 1)
      subSizes[cur] += subSizes[next]
      if (subSizes[next] > heavySonSize) {
        heavySon = next
        heavySonSize = subSizes[next]
      }
    }
    depths[cur] = depth
    heavySons[cur] = heavySon
    parents[cur] = pre
  }

  // 对这些链进行维护，就要确保每个链上的节点都是连续的
  // 注意在进行重新编号的时候先访问重链，这样可以保证重链内的节点编号连续
  function dfs2(cur: number, top: number): void {
    tops[cur] = top
    dfsId++
    dfsIds[cur] = dfsId
    dfsIdToNode[dfsId] = cur
    if (heavySons[cur] !== -1) {
      // 优先遍历重儿子，保证在同一条重链上的点的 DFS 序是连续的
      dfs2(heavySons[cur], top)
      for (let i = 0; i < adjList[cur].length; i++) {
        const next = adjList[cur][i]
        if (next !== parents[cur] && next !== heavySons[cur]) {
          dfs2(next, next)
        }
      }
    }
  }

  // 把根节点到id路径上的点全部置为1 输出这次操作安装状态被改变的软件包数量。
  function install(id: number): number {
    let res = 0
    while (id !== -1) {
      const top = tops[id]
      const heavyLength = dfsIds[id] - dfsIds[top] + 1
      res += heavyLength - tree.query(dfsIds[top], dfsIds[id])
      tree.update(dfsIds[top], dfsIds[id], 1)
      id = parents[top]
    }
    return res
  }

  // 把以id为根的子树中的点全部置为0 输出这次操作安装状态被改变的软件包数量。
  function uninstall(id: number): number {
    const left = dfsIds[id]
    const right = dfsIds[id] + subSizes[id] - 1
    let res = tree.query(left, right)
    tree.update(left, right, 0)
    return res
  }

  return {
    install,
    uninstall
  }
}

if (require.main === module) {
  const n = 7
  const deps = [0, 0, 0, 1, 1, 5]
  const adjList: number[][] = Array.from({ length: n + 1 }, () => [])
  for (let [cur, pre] of deps.entries()) {
    ;[cur, pre] = [cur + 2, pre + 1]
    adjList[pre].push(cur)
    adjList[cur].push(pre)
  }

  // 根结点为1
  const { install, uninstall } = useSoftwareManager(n, adjList, 1)
  console.log(install(6))
  console.log(install(7))
  console.log(uninstall(2))
  console.log(install(5))
  console.log(uninstall(1))
  3
  1
  3
  2
  3
}

export {}
