// 维护分量和的并查集
//
// Union Find Data Structure
//
// Description:
//   An union-find data structure (aka. disjoint set data structure)
//   maintains a disjoint sets and supports the following operations.
//   - unite(u, v): merge sets containing u and v.
//   - find(u, v) : return true if u and v are in the same set
//   - size(u)    : size of the set containing u.
//
//   The weighted version additionally maintains the values for the
//   elements and supports the following operations:
//   - add(u, a)   : value[u] += a
//   - addSet(u, a): value[v] += a for v in set(u)
//   - get(u)      : return value[u]
//   - getSet(u)   : return sum(value[v] for v in set(u))
//
// Complexity:
//   Amortized O(a(n)) for all operations.
//   Here, a(n) is the inverse Ackermann function, which is
//   less than five for a realistic size input.
//
// Verified:
//   AOJ1330 http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1330
//   and other many problems.
//

/**
 * 维护联通分量和的并查集.
 */
class WeightedUnionFind {
  private readonly _parent: Int32Array
  private readonly _value: Float64Array
  private readonly _delta: Float64Array
  private readonly _total: Float64Array
  private _part: number

  constructor(n: number) {
    this._parent = new Int32Array(n).fill(-1)
    this._value = new Float64Array(n)
    this._delta = new Float64Array(n)
    this._total = new Float64Array(n)
    this._part = n
  }

  /**
   * u的值加上delta.
   */
  addWeight(u: number, delta: number): void {
    this._value[u] += delta
    this._total[this.find(u)] += delta
  }

  /**
   * u所在集合的值加上delta.
   */
  addGroupWeight(u: number, delta: number): void {
    const root = this.find(u)
    this._delta[root] += delta
    this._total[root] -= this._parent[root] * delta
  }

  /**
   * u的值.
   */
  getWeight(u: number): number {
    return this._value[u] + this._find(u).delta
  }

  /**
   * u所在集合的值.
   */
  getGroupWeight(u: number): number {
    return this._total[this.find(u)]
  }

  union(u: number, v: number, f?: (big: number, small: number) => void): boolean {
    u = this.find(u)
    v = this.find(v)
    if (u === v) return false
    if (this._parent[u] > this._parent[v]) {
      u ^= v
      v ^= u
      u ^= v
    }
    this._parent[u] += this._parent[v]
    this._parent[v] = u
    this._delta[v] -= this._delta[u]
    this._total[u] += this._total[v]
    this._part--
    f && f(u, v)
    return true
  }

  find(u: number): number {
    return this._find(u).root
  }

  isConnected(u: number, v: number): boolean {
    return this.find(u) === this.find(v)
  }

  getSize(u: number): number {
    return -this._parent[this.find(u)]
  }

  getGroups(): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < this._parent.length; ++i) {
      const root = this.find(i)
      if (groups.has(root)) {
        groups.get(root)!.push(i)
      } else {
        groups.set(root, [i])
      }
    }
    return groups
  }

  get part(): number {
    return this._part
  }

  private _find(u: number): { root: number; delta: number } {
    if (this._parent[u] < 0) return { root: u, delta: this._delta[u] }
    let { root, delta } = this._find(this._parent[u])
    delta += this._delta[u]
    this._parent[u] = root
    this._delta[u] = delta - this._delta[root]
    return { root, delta }
  }
}

export { WeightedUnionFind, WeightedUnionFind as UnionFindWeighted }

// 2382. 删除操作后的最大子段和
// https://leetcode.cn/problems/maximum-segment-sum-after-removals/
//
function maximumSegmentSum(nums: number[], removeQueries: number[]): number[] {
  const n = nums.length
  const q = removeQueries.length
  const uf = new WeightedUnionFind(n)
  const res: number[] = Array(q).fill(0)
  let maxPartSum = 0
  const visited = new Uint8Array(n)
  for (let qi = q - 1; ~qi; qi--) {
    res[qi] = maxPartSum
    const pos = removeQueries[qi]
    uf.addWeight(pos, nums[pos])
    visited[pos] = 1
    if (pos > 0 && visited[pos - 1]) uf.union(pos, pos - 1)
    if (pos < n - 1 && visited[pos + 1]) uf.union(pos, pos + 1)
    maxPartSum = Math.max(maxPartSum, uf.getGroupWeight(pos))
  }
  return res
}
