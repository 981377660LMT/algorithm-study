// Floyd-Warshall 算法求有向图的传递闭包。通俗地讲，就是可达性问题。
// O(n^3/w)
//
// new TransitiveClosure(n) 构造一个n*n的传递闭包.
// addDirectedEdge(from, to) 添加一条有向边.
// build() 构造传递闭包.
// canReach(from, to) 判断是否可达.

/**
 * Floyd-Warshall 算法求有向图的闭包传递问题，也就是可达性问题。
 * O(n^3/32).
 * - 3000*3000 => 930ms
 * - 4000*4000 => 2.3s
 * - 5000*5000 => 4.2s
 */
class TransitiveClosure {
  private readonly _n: number
  private readonly _canReach: Uint32Array[] // bitset[]
  private _hasBuilt = false

  constructor(n: number) {
    const canReach = Array(n)
    for (let i = 0; i < n; i++) canReach[i] = new Uint32Array((n >>> 5) + 1)
    this._n = n
    this._canReach = canReach
  }

  addDirectedEdge(from: number, to: number): void {
    if (this._hasBuilt) throw new Error("can't add edge after build")
    this._canReach[from][to >>> 5] |= 1 << (to & 31)
  }

  build(): void {
    if (this._hasBuilt) throw new Error("can't build twice")

    const n = this._n
    const canReach = this._canReach
    for (let k = 0; k < n; k++) {
      const canReachK = canReach[k]
      for (let i = 0; i < n; i++) {
        const canReachI = canReach[i]
        if (canReachI[k >>> 5] & (1 << (k & 31))) {
          for (let j = 0; j < canReachK.length; j++) {
            canReachI[j] |= canReachK[j]
          }
        }
      }
    }

    this._hasBuilt = true
  }

  canReach(from: number, to: number): boolean {
    if (!this._hasBuilt) this.build()
    return !!(this._canReach[from][to >>> 5] & (1 << (to & 31)))
  }
}

export { TransitiveClosure }

if (require.main === module) {
  const n = 5000
  const T = new TransitiveClosure(n)
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < n; j++) {
      T.addDirectedEdge(i, j)
    }
  }
  console.time('build')
  T.build()
  console.timeEnd('build')
}
