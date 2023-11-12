/* eslint-disable no-param-reassign */

/**
 * 树上区间并查集.
 */
class UnionFindRangeOnTree {
  private readonly _n: number
  private readonly _data: Int32Array
  private readonly _treeParents: ArrayLike<number>
  private _part: number

  constructor(n: number, treeParents: ArrayLike<number>) {
    this._n = n
    this._data = new Int32Array(n).fill(-1)
    this._treeParents = treeParents
    this._part = n
  }

  /**
   * 将child结点合并到parent结点上,返回是否合并成功.
   */
  union(
    parent: number,
    child: number,
    f?: (parentRoot: number, childRoot: number) => void
  ): boolean {
    parent = this.find(parent)
    child = this.find(child)
    if (parent === child) return false
    this._data[parent] += this._data[child]
    this._data[child] = parent
    this._part--
    f && f(parent, child)
    return true
  }

  /**
   * 定向合并从祖先ancestor到子孙child路径上的所有节点,返回合并次数.
   * 合并方向为child->ancestor.
   */
  unionRange(
    ancestor: number,
    child: number,
    f?: (ancestorRoot: number, childRoot: number) => void
  ): number {
    const target = this.find(ancestor)
    let mergeCount = 0
    while (true) {
      child = this.find(child)
      if (child === target) break
      this.union(this._treeParents[child], child, f)
      mergeCount++
    }
    return mergeCount
  }

  find(x: number): number {
    if (this._data[x] < 0) return x
    // eslint-disable-next-line no-return-assign
    return (this._data[x] = this.find(this._data[x]))
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getSize(x: number): number {
    return -this._data[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const group = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const root = this.find(i)
      if (!group.has(root)) group.set(root, [])
      group.get(root)!.push(i)
    }
    return group
  }

  get part(): number {
    return this._part
  }
}

export { UnionFindRangeOnTree }

if (require.main === module) {
  const n = 5
  const parents = [-1, 0, 1, 2, 3]
  const uf = new UnionFindRangeOnTree(n, parents)
  console.log(uf.unionRange(0, 4))
  console.log(uf.getGroups())
}
