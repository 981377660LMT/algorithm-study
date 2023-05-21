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

class WeightedUnionFind {
  private readonly _parent: Int32Array
  private readonly _value: number[]
  private readonly _delta: number[]
  private readonly _total: number[]
  private _part: number

  constructor(n: number) {
    this._parent = new Int32Array(n).fill(-1)
    this._value = Array(n).fill(0)
    this._delta = Array(n).fill(0)
    this._total = Array(n).fill(0)
    this._part = n
  }

  /**
   * u的值加上delta.
   */
  add(u: number, delta: number): void {
    this._value[u] += delta
    this._total[this.find(u)] += delta
  }

  /**
   * u所在集合的值加上delta.
   */
  addGroup(u: number, delta: number): void {
    const root = this.find(u)
    this._delta[root] += delta
    this._total[root] -= this._parent[root] * delta
  }

  /**
   * u的值.
   */
  get(u: number): number {
    return this._value[u] + this._find(u)[1]
  }

  /**
   * u所在集合的值.
   */
  getGroup(u: number): number {
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
    return this._find(u)[0]
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

  private _find(u: number): [root: number, delta: number] {
    if (this._parent[u] < 0) return [u, this._delta[u]]
    let [a, b] = this._find(this._parent[u])
    b += this._delta[u]
    this._parent[u] = a
    this._delta[u] = b - this._delta[a]
    return [a, b]
  }
}

export { WeightedUnionFind }
