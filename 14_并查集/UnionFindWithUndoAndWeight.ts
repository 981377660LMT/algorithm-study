/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// UnionFindWithUndoAndWeight/UnionFindArrayWithUndoAndWeight
// https://hitonanode.github.io/cplib-cpp/unionfind/undo_monoid_unionfind.hpp
// 可撤销并查集 / 维护 `满足可换律的monoid` 权值 的并查集
// SetGroup: 将下标为index元素`所在集合`的权值置为value.
// GetGroup: 获取下标为index元素`所在集合`的权值.
// !Undo: 撤销上一次合并(Union)或者修改权值(Set)操作，没合并成功也要撤销
// Reset: 撤销所有操作

/**
 * 可撤销并查集, 维护连通分量的权值.
 */
class UnionFindArrayWithUndoAndWeight<S> {
  private readonly _ranks: Uint32Array
  private readonly _parents: Uint32Array
  private readonly _weights: S[]
  private readonly _history: { root: number; rank: number; weight: S }[] = []
  private readonly _op: (s1: S, s2: S) => S
  private _part: number

  constructor(initWeights: S[], op: (s1: S, s2: S) => S) {
    const n = initWeights.length
    this._ranks = new Uint32Array(n)
    this._parents = new Uint32Array(n)
    this._weights = Array(n)
    for (let i = 0; i < n; ++i) {
      this._ranks[i] = 1
      this._parents[i] = i
      this._weights[i] = initWeights[i]
    }
    this._op = op
    this._part = n
  }

  /**
   * 将下标为index元素`所在集合`的权值置为value.
   */
  setGroupWeight(index: number, value: S): void {
    index = this.find(index)
    this._history.push({ root: index, rank: this._ranks[index], weight: this._weights[index] })
    this._weights[index] = value
  }

  /**
   * 获取下标为index元素`所在集合`的权值.
   */
  getGroupWeight(index: number): S {
    return this._weights[this.find(index)]
  }

  /**
   * 撤销上一次合并(Union)或者修改权值(Set)操作.
   * 没合并成功也要撤销.
   */
  undo(): void {
    if (!this._history.length) return
    const { root, rank, weight } = this._history.pop()!
    const ps = this._parents[root]
    this._weights[ps] = weight
    this._ranks[ps] = rank
    if (ps !== root) {
      this._parents[root] = root
      this._part++
    }
  }

  /**
   * 撤销所有操作.
   */
  reset(): void {
    while (this._history.length) {
      this.undo()
    }
  }

  find(x: number): number {
    if (this._parents[x] === x) return x
    return this.find(this._parents[x])
  }

  union(x: number, y: number): boolean {
    x = this.find(x)
    y = this.find(y)
    if (this._ranks[x] < this._ranks[y]) {
      x ^= y
      y ^= x
      x ^= y
    }
    this._history.push({ root: y, rank: this._ranks[x], weight: this._weights[x] })
    if (x !== y) {
      this._parents[y] = x
      this._ranks[x] += this._ranks[y]
      this._weights[x] = this._op(this._weights[x], this._weights[y])
      this._part--
      return true
    }
    return false
  }

  isConnected(x: number, y: number): boolean {
    return this.find(x) === this.find(y)
  }

  getSize(x: number): number {
    return this._ranks[this.find(x)]
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < this._parents.length; ++i) {
      const root = this.find(i)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(i)
    }
    return groups
  }

  get part(): number {
    return this._part
  }
}

export { UnionFindArrayWithUndoAndWeight }

if (require.main === module) {
  const uf = new UnionFindArrayWithUndoAndWeight([1, 2, 3, 4, 5], (a, b) => a + b)
  uf.union(0, 1)
  uf.union(1, 2)
  uf.union(3, 4)
  console.log(uf.getGroups())
  console.log(
    uf.getGroupWeight(0),
    uf.getGroupWeight(1),
    uf.getGroupWeight(2),
    uf.getGroupWeight(3),
    uf.getGroupWeight(4)
  )
  uf.undo()
  console.log(
    uf.getGroupWeight(0),
    uf.getGroupWeight(1),
    uf.getGroupWeight(2),
    uf.getGroupWeight(3),
    uf.getGroupWeight(4)
  )
  uf.reset()
  console.log(
    uf.getGroupWeight(0),
    uf.getGroupWeight(1),
    uf.getGroupWeight(2),
    uf.getGroupWeight(3),
    uf.getGroupWeight(4)
  )
  console.log(uf.getGroups())
  uf.union(0, 1)
  console.log(uf.getGroups())
  uf.setGroupWeight(0, 2)
  console.log(
    uf.getGroupWeight(0),
    uf.getGroupWeight(1),
    uf.getGroupWeight(2),
    uf.getGroupWeight(3),
    uf.getGroupWeight(4)
  )
}
