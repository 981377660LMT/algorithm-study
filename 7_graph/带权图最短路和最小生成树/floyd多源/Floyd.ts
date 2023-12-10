const INF = 2e15

/**
 * 返回一个函数,该函数可以求出从`start`到`target`的最短路径长度.
 * 如果不存在这样的路径,返回`-1`.
 */
function floyd(
  n: number,
  edges: [u: number, v: number, w: number][] | number[][],
  directed = false
): (start: number, target: number) => number {
  const dist: number[] = Array(n * n)
  for (let i = 0; i < n * n; ++i) dist[i] = INF
  for (let i = 0; i < n; ++i) dist[i * n + i] = 0

  if (directed) {
    edges.forEach(([u, v, w]) => {
      dist[u * n + v] = Math.min(dist[u * n + v], w)
    })
  } else {
    edges.forEach(([u, v, w]) => {
      dist[u * n + v] = Math.min(dist[u * n + v], w)
      dist[v * n + u] = Math.min(dist[v * n + u], w)
    })
  }

  for (let k = 0; k < n; ++k) {
    for (let i = 0; i < n; ++i) {
      for (let j = 0; j < n; ++j) {
        const cand = dist[i * n + k] + dist[k * n + j]
        if (dist[i * n + j] > cand) {
          dist[i * n + j] = cand
        }
      }
    }
  }

  return (start: number, target: number) => {
    const res = dist[start * n + target]
    return res === INF ? -1 : res
  }
}

class Floyd {
  private _hasBuilt = false
  private readonly _n: number
  private readonly _dist: number[]

  /**
   * pre[a][b]表示a作为最短路起点，a->b的最短路上b的前驱节点.
   */
  private readonly _pre: Uint32Array
  private readonly _directedEdges: [u: number, v: number, w: number][] = []

  constructor(n: number) {
    const dist = Array(n * n)
    const pre = new Uint32Array(n * n)
    for (let i = 0; i < n; ++i) {
      for (let j = 0; j < n; ++j) {
        const cur = i * n + j
        pre[cur] = i
        dist[cur] = INF
      }
    }
    for (let i = 0; i < n; ++i) dist[i * n + i] = 0

    this._n = n
    this._pre = pre
    this._dist = dist
  }

  /**
   * 添加从`u`到`v`的边权为`w`的无向边.
   */
  addEdge(u: number, v: number, w: number): void {
    this.addDirectedEdge(u, v, w)
    this.addDirectedEdge(v, u, w)
  }

  /**
   * 添加从`u`到`v`的边权为`w`的有向边.
   */
  addDirectedEdge(u: number, v: number, w: number): void {
    if (this._hasBuilt) throw new Error('Can not add edge after build.')
    this._directedEdges.push([u, v, w])
  }

  /**
   * @complexity O(n^3)
   */
  build(): void {
    if (this._hasBuilt) throw new Error('Can not build twice.')
    const n = this._n
    this._directedEdges.forEach(([u, v, w]) => {
      this._dist[u * n + v] = Math.min(this._dist[u * n + v], w)
    })

    for (let k = 0; k < n; ++k) {
      for (let i = 0; i < n; ++i) {
        for (let j = 0; j < n; ++j) {
          const cand = this._dist[i * n + k] + this._dist[k * n + j]
          if (this._dist[i * n + j] > cand) {
            this._dist[i * n + j] = cand
            this._pre[i * n + j] = this._pre[k * n + j]
          }
        }
      }
    }

    this._hasBuilt = true
  }

  /**
   * 求出从`start`到`target`的最短路径长度,如果不存在这样的路径,返回-1.
   * @complexity O(1)
   */
  dist(start: number, target: number): number {
    if (!this._hasBuilt) this.build()
    const res = this._dist[start * this._n + target]
    return res === INF ? -1 : res
  }

  /**
   * 求出从`start`到`target`的最短路径.如果不存在这样的路径,返回空数组.
   */
  getPath(start: number, target: number): number[] {
    if (!this._hasBuilt) this.build()
    if (this.dist(start, target) === -1) return []
    let cur = target
    const path = [target]
    while (cur !== start) {
      cur = this._pre[start * this._n + cur]
      path.push(cur)
    }
    return path.reverse()
  }

  /**
   * 判断是否存在负环.
   * @complexity O(n)
   */
  hasNegativeCycle(): boolean {
    if (!this._hasBuilt) this.build()
    const n = this._n
    for (let i = 0; i < n; ++i) {
      if (this._dist[i * n + i] < 0) return true
    }
    return false
  }
}

/**
 * 动态Floyd算法,支持向图中添加边.
 */
class FloydDynamic {
  private _hasBuilt = false
  private readonly _n: number
  private readonly _dist: number[]
  private readonly _directedEdges: [u: number, v: number, w: number][] = []

  constructor(n: number) {
    const dist = Array(n * n)
    for (let i = 0; i < n * n; ++i) dist[i] = INF
    for (let i = 0; i < n; ++i) dist[i * n + i] = 0
    this._n = n
    this._dist = dist
  }

  /**
   * 添加从`u`到`v`的边权为`w`的无向边.
   */
  addEdge(u: number, v: number, w: number): void {
    this.addDirectedEdge(u, v, w)
    this.addDirectedEdge(v, u, w)
  }

  /**
   * 添加从`u`到`v`的边权为`w`的有向边.
   */
  addDirectedEdge(u: number, v: number, w: number): void {
    if (this._hasBuilt) throw new Error('Can not add edge after build.')
    this._directedEdges.push([u, v, w])
  }

  /**
   * @complexity O(n^3)
   */
  build(): void {
    if (this._hasBuilt) throw new Error('Can not build twice.')
    const n = this._n
    this._directedEdges.forEach(([u, v, w]) => {
      this._dist[u * n + v] = Math.min(this._dist[u * n + v], w)
    })

    for (let k = 0; k < n; ++k) {
      for (let i = 0; i < n; ++i) {
        for (let j = 0; j < n; ++j) {
          const cand = this._dist[i * n + k] + this._dist[k * n + j]
          if (this._dist[i * n + j] > cand) {
            this._dist[i * n + j] = cand
          }
        }
      }
    }

    this._hasBuilt = true
  }

  /**
   * 向边集中添加一条边.
   * @param directed 是否为有向边
   * @complexity O(n^2)
   * 加边时,枚举每个点对,根据是否经过edge来更新最短路.
   */
  updateEdge(u: number, v: number, w: number, directed = true): void {
    if (!this._hasBuilt) this.build()
    if (directed) {
      this._updateDirectedEdge(u, v, w)
    } else {
      this._updateEdge(u, v, w)
    }
  }

  /**
   * 求出从`start`到`target`的最短路径长度,如果不存在这样的路径,返回-1.
   * @complexity O(1)
   */
  dist(start: number, target: number): number {
    if (!this._hasBuilt) this.build()
    const res = this._dist[start * this._n + target]
    return res === INF ? -1 : res
  }

  /**
   * 判断是否存在负环.
   * @complexity O(n)
   */
  hasNegativeCycle(): boolean {
    if (!this._hasBuilt) this.build()
    const n = this._n
    for (let i = 0; i < n; ++i) {
      if (this._dist[i * n + i] < 0) return true
    }
    return false
  }

  private _updateDirectedEdge(u: number, v: number, w: number): void {
    const n = this._n
    for (let i = 0; i < n; ++i) {
      for (let j = 0; j < n; ++j) {
        const cand = this._dist[i * n + u] + w + this._dist[v * n + j]
        if (this._dist[i * n + j] > cand) {
          this._dist[i * n + j] = cand
        }
      }
    }
  }

  private _updateEdge(u: number, v: number, w: number): void {
    const n = this._n
    for (let i = 0; i < n; ++i) {
      for (let j = 0; j < n; ++j) {
        const cand1 = this._dist[i * n + u] + w + this._dist[v * n + j]
        if (this._dist[i * n + j] > cand1) {
          this._dist[i * n + j] = cand1
        }
        const cand2 = this._dist[i * n + v] + w + this._dist[u * n + j]
        if (this._dist[i * n + j] > cand2) {
          this._dist[i * n + j] = cand2
        }
      }
    }
  }
}

export { Floyd, FloydDynamic, floyd }

if (require.main === module) {
  const n = 800
  console.time('floyd')
  const f = new FloydDynamic(n)
  f.build()
  console.timeEnd('floyd')
}
