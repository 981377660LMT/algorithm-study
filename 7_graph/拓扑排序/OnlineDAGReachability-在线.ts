//
// Italiano's dynamic reachability data structure for DAG
//
// Description:
//   It is a data structure that admits the following operations:
//     add_edge(s, t):     insert edge (s,t) to the network if
//                         it does not make a cycle
//
//     is_reachable(s, t): return true if there is a path s --> t
//
// Algorithm:
//   We maintain reachability trees T(u) for all u in V.
//   Then is_reachable(s, t) is solved by checking "t in T(u)".
//   For add_edge(s, t), if is_reachable(s, t) or is_reachable(t, s) then
//   no update is performed. Otherwise, we meld T(s) and T(t).
//
// Complexity:
//   !amortized O(n) per update
//
// Verified:
//   SPOJ 9458: Ghosts having fun
//
// References:
//   Giuseppe F. Italiano (1988):
//   Finding paths and deleting edges in directed acyclic graphs.
//   Information Processing Letters, vol. 28, no. 1, pp. 5--11.
//
// 维护DAG可到达性.
// n<=1e4

/**
 * 动态DAG可到达性.
 */
class OnlineDAGReachability {
  private readonly _n: number
  private readonly _parent: Int16Array
  private readonly _child: number[][]

  /**
   * @param n n<=1e4.
   */
  constructor(n: number) {
    const size = n * n
    const parent = new Int16Array(size)
    const child = Array(size)
    for (let i = 0; i < size; i++) {
      parent[i] = -1
      child[i] = []
    }
    this._n = n
    this._parent = parent
    this._child = child
  }

  /**
   * 判断是否可从`from`到达`to`.
   */
  isReachable(from: number, to: number): boolean {
    return from === to || this._parent[from * this._n + to] >= 0
  }

  /**
   * 添加有向边`from->to`.
   * 如果添加后形成环, 则返回false.
   * @complexity `O(n)`
   */
  addEdge(from: number, to: number): boolean {
    if (this.isReachable(to, from)) return false
    if (this.isReachable(from, to)) return true
    for (let p = 0; p < this._n; p++) {
      if (this.isReachable(p, from) && !this.isReachable(p, to)) {
        this._meld(p, to, from, to)
      }
    }
    return true
  }

  private _meld(root: number, sub: number, u: number, v: number): void {
    this._parent[root * this._n + v] = u
    this._child[root * this._n + u].push(v)
    this._child[sub * this._n + v].forEach(c => {
      if (!this.isReachable(root, c)) {
        this._meld(root, sub, v, c)
      }
    })
  }
}

export { OnlineDAGReachability }

if (require.main === module) {
  const R = new OnlineDAGReachability(5)

  R.addEdge(0, 1)
  console.log(R.isReachable(0, 1))
  console.log(R.isReachable(1, 0))

  R.addEdge(1, 2)
  console.log(R.isReachable(0, 2))
  console.log(R.isReachable(2, 0))

  R.addEdge(2, 3)
  console.log(R.isReachable(0, 3))
  console.log(R.isReachable(3, 0))

  R.addEdge(3, 4)
  console.log(R.isReachable(0, 4))
  console.log(R.isReachable(4, 0))
  R.addEdge(4, 0)
}
