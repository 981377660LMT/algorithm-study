/* eslint-disable no-inner-declarations */

// RangeToRangeGraph (区间图)
// !原图的连通分量/最短路在新图上仍然等价
// 线段树优化建图

const INF = 2e15

/**
 * 线段树优化区间建图.
 */
class RangeToRangeGraph {
  private readonly _n: number
  private readonly _edges: [u: number, v: number, w: number][] = []
  private _nodeCount: number

  constructor(n: number) {
    this._n = n
    this._nodeCount = n * 3
    for (let i = 2; i < n + n; i++) {
      this._edges.push([this._toUpperIdx(i >> 1), this._toUpperIdx(i), 0])
    }
    for (let i = 2; i < n + n; i++) {
      this._edges.push([this._toLowerIdx(i), this._toLowerIdx(i >> 1), 0])
    }
  }

  /**
   * 返回`新图的有向邻接表`和`新图的节点数`.
   */
  build(): [newGraph: [next: number, weight: number][][], vertex: number] {
    const vertex = this._nodeCount
    const adjList = Array(vertex).fill(0)
    for (let i = 0; i < adjList.length; i++) {
      adjList[i] = []
    }
    this._edges.forEach(([u, v, w]) => {
      adjList[u].push([v, w])
    })
    return [adjList, vertex]
  }

  /**
   * 添加有向边 from -> to, 权重为 weight.
   */
  add(from: number, to: number, weight: number): void {
    this._edges.push([from, to, weight])
  }

  /**
   * 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条有向边，权重为 weight.
   */
  addFromRange(fromStart: number, fromEnd: number, to: number, weight: number): void {
    let left = fromStart + this._n
    let right = fromEnd + this._n
    while (left < right) {
      if (left & 1) {
        this.add(this._toLowerIdx(left), to, weight)
        left++
      }
      if (right & 1) {
        right--
        this.add(this._toLowerIdx(right), to, weight)
      }
      left >>= 1
      right >>= 1
    }
  }

  /**
   * 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
   */
  addToRange(from: number, toStart: number, toEnd: number, weight: number): void {
    let left = toStart + this._n
    let right = toEnd + this._n
    while (left < right) {
      if (left & 1) {
        this.add(from, this._toUpperIdx(left), weight)
        left++
      }
      if (right & 1) {
        right--
        this.add(from, this._toUpperIdx(right), weight)
      }
      left >>= 1
      right >>= 1
    }
  }

  /**
   * 从区间 [fromStart, fromEnd) 中的每个点到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
   */
  addRangeToRange(
    fromStart: number,
    fromEnd: number,
    toStart: number,
    toEnd: number,
    weight: number
  ): void {
    const newNode = this._nodeCount
    this._nodeCount++
    this.addFromRange(fromStart, fromEnd, newNode, weight)
    this.addToRange(newNode, toStart, toEnd, 0)
  }

  private _toUpperIdx(i: number): number {
    return i >= this._n ? i - this._n : this._n + i
  }

  private _toLowerIdx(i: number): number {
    return i >= this._n ? i - this._n : this._n + this._n + i
  }
}

if (require.main === module) {
  // https://leetcode.cn/problems/zui-xiao-tiao-yue-ci-shu/
  // LCP 09. 最小跳跃次数 (MLE)
  function minJump(jump: number[]): number {
    const n = jump.length
    const G = new RangeToRangeGraph(n + 1)
    for (let i = 0; i < n; i++) {
      // i => [0,i)
      // i => Math.min(i + jump[i], n)
      G.addToRange(i, 0, i, 1)
      G.add(i, Math.min(i + jump[i], n), 1)
    }
    const [adjList] = G.build()
    const dist = bfs(adjList, 0)
    return dist[n]
  }

  function bfs(adjList: [next: number, weight: number][][], start: number) {
    const dist = Array(adjList.length).fill(INF)
    dist[start] = 0
    let stack = [start]
    while (stack.length) {
      const cur = stack.pop()!
      adjList[cur].forEach(([next, weight]) => {
        const newDist = dist[cur] + weight
        if (newDist < dist[next]) {
          dist[next] = newDist
          stack.push(next)
        }
      })
    }
    return dist
  }
}

export { RangeToRangeGraph }
