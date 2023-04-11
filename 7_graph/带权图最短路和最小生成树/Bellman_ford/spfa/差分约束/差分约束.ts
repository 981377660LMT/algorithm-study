// !如果所有边权非负,可以把spfa换成dijkstra

import { Heap } from '../../../../../8_heap/Heap'

const INF = 2e15

/**
 * 差分约束
 */
class DualShortestPath {
  private readonly _n: number
  private readonly _g: [next: number, weight: number][][]
  private readonly _min: boolean
  private _hasNeg = false

  constructor(n: number, min: boolean) {
    this._n = n
    this._g = Array(n).fill(0)
    for (let i = 0; i < n; i++) {
      this._g[i] = []
    }
    this._min = min
  }

  /**
   * f(j) <= f(i) + w
   */
  addEdge(i: number, j: number, w: number): void {
    if (this._min) {
      this._g[i].push([j, w])
    } else {
      this._g[j].push([i, w])
    }
    this._hasNeg = this._hasNeg || w < 0
  }

  run(): [dist: number[], ok: boolean] {
    if (this._min) return this._spfaMin()
    return this._hasNeg ? this._spfaMax() : this._dijkMax()
  }

  private _spfaMin(): [dist: number[], ok: boolean] {
    const dist = Array(this._n).fill(0)
    const count = new Uint32Array(this._n)
    const inStack = new Uint8Array(this._n)
    const stack = new Uint32Array(this._n)
    for (let i = 0; i < this._n; i++) {
      count[i] = 1
      inStack[i] = 1
      stack[i] = i
    }
    let ptr = 0
    let len = this._n

    while (len) {
      const cur = stack[ptr++]
      len--
      inStack[cur] = 0
      for (let i = 0; i < this._g[cur].length; i++) {
        const [next, weight] = this._g[cur][i]
        const cand = dist[cur] + weight
        if (cand < dist[next]) {
          dist[next] = cand
          if (!inStack[next]) {
            count[next]++
            if (count[next] >= this._n) {
              return [[], false]
            }
            inStack[next] = 1
            stack[--ptr] = next
            len++
          }
        }
      }
    }

    return [dist.map(num => -num), true]
  }

  private _spfaMax(): [dist: number[], ok: boolean] {
    const dist = Array(this._n).fill(INF)
    const count = new Uint32Array(this._n)
    const inStack = new Uint8Array(this._n)
    let queue = [0]
    dist[0] = 0
    inStack[0] = 1
    count[0] = 1
    while (queue.length) {
      const nextQueue: number[] = []
      for (let i = 0; i < queue.length; i++) {
        const cur = queue[i]
        inStack[cur] = 0
        for (let j = 0; j < this._g[cur].length; j++) {
          const [next, weight] = this._g[cur][j]
          const cand = dist[cur] + weight
          if (cand < dist[next]) {
            dist[next] = cand
            if (!inStack[next]) {
              count[next]++
              if (count[next] >= this._n) {
                return [[], false]
              }
              inStack[next] = 1
              nextQueue.push(next)
            }
          }
        }
      }

      queue = nextQueue
    }

    return [dist, true]
  }

  private _dijkMax(): [dist: number[], ok: boolean] {
    const dist = Array(this._n).fill(INF)
    dist[0] = 0
    const heap = new Heap<[dist: number, node: number]>((a, b) => a[0] - b[0])
    heap.push([0, 0])
    while (heap.size) {
      const [curDist, cur] = heap.pop()!
      if (dist[cur] < curDist) continue
      this._g[cur].forEach(([next, weight]) => {
        const cand = dist[cur] + weight
        if (cand < dist[next]) {
          dist[next] = cand
          heap.push([cand, next])
        }
      })
    }
    return [dist, true]
  }
}

if (require.main === module) {
  const [n, m] = [10, 3]
  const limits = [
    [1, 4, 2],
    [3, 6, 2],
    [10, 10, 1]
  ]
  const D = new DualShortestPath(n + 10, false)
  limits.forEach(([l, r, w]) => {
    D.addEdge(r, l - 1, w)
  })
  for (let i = 1; i <= n + 1; i++) {
    D.addEdge(i - 1, i, 0)
    D.addEdge(i, i - 1, 1)
  }
  const [dist, ok] = D.run()
  if (!ok) {
    console.log('No solution')
  }
  console.log(dist[n])
}

export { DualShortestPath }
