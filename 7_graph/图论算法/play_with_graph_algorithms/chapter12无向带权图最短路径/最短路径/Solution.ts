import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
import { UnionFind } from '../../chapter11无向带权图最小生成树/带权图/UnionFind'
import {
  WeightedAdjList,
  WeightedEdge,
} from '../../chapter11无向带权图最小生成树/带权图/WeighedAdjList'

class Solution {
  public readonly adjList: WeightedAdjList
  public readonly edgeList: WeightedEdge[]
  private readonly cc: DFS
  private readonly uf: UnionFind

  private constructor(weightedAdjList: WeightedAdjList, cc: DFS) {
    this.adjList = weightedAdjList
    this.cc = cc

    const edges = weightedAdjList.adjList
      .flatMap((map, v1) => {
        const tmp = []
        for (const [v2, weight] of map.entries()) {
          v1 < v2 && tmp.push(new WeightedEdge(v1, v2, weight))
        }
        return tmp
      })
      .sort((e1, e2) => e1.weight - e2.weight)
    this.edgeList = edges

    this.uf = new UnionFind<number>()
    edges.forEach(edge => this.uf.add(edge.v1).add(edge.v2))
  }

  /**
   * @param start 原点
   * @description visited记录哪些已经确定最小值
   * @description
   * 1. 找最近的点v
   * 2. v加入visited
   * 3. 利用v更新相邻的没看过的点
   * @description 复杂度O(V^2)
   * @description pre数组求解路径，更新与dis数组更新同步
   * @description 不能处理负权边
   */
  dijkstra(start: number): {
    dis: number[]
    pre: number[]
  } {
    // dis[w]表示原点start到w的当前最短距离
    // Infinity哨兵
    const dis = Array<number>(this.adjList.V).fill(Infinity)
    dis[start] = 0
    const pre = Array<number>(this.adjList.V).fill(-1)
    pre[start] = start

    const visited = new Set<number>()

    while (true) {
      // 1.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
      let nearestDis = Infinity
      let nearestV = -1
      for (let v = 0; v < this.adjList.V; v++) {
        if (!visited.has(v) && dis[v] < nearestDis) {
          nearestDis = dis[v]
          nearestV = v
        }
      }

      // 都更新完了
      if (nearestV === -1) break

      // 2.加入visited
      visited.add(nearestV)

      // 3.利用nearestV点来更新其相邻节点w与原点s的距离
      for (const next of this.adjList.adj(nearestV)) {
        if (!visited.has(next)) {
          if (dis[nearestV] + this.adjList.getWeight(nearestV, next) < dis[next]) {
            dis[next] = dis[nearestV] + this.adjList.getWeight(nearestV, next)
            pre[next] = nearestV
          }
        }
      }
    }

    return { dis, pre }
  }

  /**
   * @description
   * ```js
    (1) 初始dis[s] = 0, 其余dis为正无穷
    (2) 对所有边进行一次**松弛操作**(换条道更短)，则求出了到所有点，经过的变数最多为1的最短路径
    (3) 对所有边再进行一次松弛操作，则求出了到所有点，经过的边数最多为2的最短路径
    (4) 对所有边进行V - 1(没有负权环时需要的次数的上界，与松弛顺序有关)次松弛操作，则求出了到所有点，经过的变数最多为V - 1的最短路径
    (5) 最后再进行一次松弛操作，如果有更新最短距离dis，则肯定有负权环，没有意义
    ```
    @description 可以直接应用于有向图
    @description O(V*E)
    @description 松弛操作时改变pre数组
   */
  bellmanFord(start: number) {
    const dis = Array<number>(this.adjList.V).fill(Infinity)
    dis[start] = 0
    const pre = Array<number>(this.adjList.V).fill(-1)
    pre[start] = start
    let hasNegativeCycle = false

    // 松弛操作
    const relax = (a: number, b: number) => {
      if (dis[a] + this.adjList.getWeight(a, b) < dis[b]) {
        dis[b] = dis[a] + this.adjList.getWeight(a, b)
      }
    }

    // V-1轮松弛操作
    for (let steps = 1; steps < this.adjList.V; steps++) {
      for (let v = 0; v < this.adjList.V; v++) {
        this.adjList.adj(v).forEach(w => relax(v, w))
      }
    }

    // 判断负权环
    for (let v = 0; v < this.adjList.V; v++) {
      for (const w of this.adjList.adj(v)) {
        if (dis[v] + this.adjList.getWeight(v, w) < dis[w]) {
          hasNegativeCycle = true
        }
      }
    }

    return { dis, hasNegativeCycle }
  }

  /**
   *@description 三重循环
   */
  floyd() {}

  static async asyncBuild(fileName: string) {
    const weightedAdjList = await WeightedAdjList.asyncBuild(fileName)
    const cc = await DFS.asyncBuild('WeightedAdjList', fileName)
    return new Solution(weightedAdjList, cc)
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g.txt')
    const solution = await Solution.asyncBuild(fileName)
    console.log(solution.adjList)
    console.log(solution.edgeList)
    console.log(solution.dijkstra(0))
    console.log(solution.bellmanFord(0))
  }
  main()
}

export { Solution }
