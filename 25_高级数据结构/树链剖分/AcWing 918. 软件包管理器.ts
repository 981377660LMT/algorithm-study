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
    if (l < 1 || r > this._size)
      throw new RangeError(`[${l}, ${r}] out of range: [1, ${this._size}]`)
  }
}
// #endregion

// 结点编号 1-n
function useSoftwareManager(n: number, adjMap: Map<number, Set<number>>) {
  const tree = new SegmentTree(n + 10)
  const depths = new Uint32Array(n + 10)
  const parents = new Uint32Array(n + 10)
  const subsizes = new Uint32Array(n + 10)
  const heavysons = new Uint32Array(n + 10) // 每个点的重儿子
  const heavyTops = new Uint32Array(n + 10) // 每个点的重链顶点
  const dfsIds = new Uint32Array(n + 10)
  let id = 1

  dfs1(1, 0, 1)
  dfs2(1, 1)

  // 把根节点到x路径上的点全部置为1 输出这次操作安装状态被改变的软件包数量。
  function install(id: number): number {
    let res = 0
    while (id) {
      res += dfsIds[id] - dfsIds[heavyTops[id]] + 1 - tree.query(dfsIds[heavyTops[id]], dfsIds[id])
      tree.update(dfsIds[heavyTops[id]], dfsIds[id], 1)
      id = parents[heavyTops[id]]
    }
    return res
  }

  // 把以x为根的子树中的点全部置为0 输出这次操作安装状态被改变的软件包数量。
  function uninstall(id: number): number {
    let res = tree.query(dfsIds[id], dfsIds[id] + subsizes[id] - 1)
    tree.update(dfsIds[id], dfsIds[id] + subsizes[id] - 1, 0)
    return res
  }

  function dfs1(cur: number, pre: number, depth: number): void {
    depths[cur] = depth
    parents[cur] = pre
    subsizes[cur] = 1
    for (const next of adjMap.get(cur) ?? []) {
      if (next === pre) continue
      dfs1(next, cur, depth + 1)
      subsizes[cur] += subsizes[next]
      if (subsizes[next] > subsizes[heavysons[cur]]) {
        heavysons[cur] = next
      }
    }
  }

  function dfs2(cur: number, heavyStart: number): void {
    dfsIds[cur] = id++
    heavyTops[cur] = heavyStart
    if (heavysons[cur] === 0) return
    dfs2(heavysons[cur], heavyStart)
    for (const next of adjMap.get(cur) ?? []) {
      if (next === heavysons[cur] || next === parents[cur]) continue
      dfs2(next, next)
    }
  }

  return {
    install,
    uninstall,
  }
}

if (require.main === module) {
  const n = 7
  const deps = [0, 0, 0, 1, 1, 5]
  const adjMap = new Map<number, Set<number>>()
  for (let [cur, pre] of deps.entries()) {
    ;[cur, pre] = [cur + 2, pre + 1]
    !adjMap.has(cur) && adjMap.set(cur, new Set())
    !adjMap.has(pre) && adjMap.set(pre, new Set())
    adjMap.get(cur)!.add(pre)
    adjMap.get(pre)!.add(cur)
  }

  const { install, uninstall } = useSoftwareManager(n, adjMap)
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
