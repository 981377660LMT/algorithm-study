import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
import { PriorityQueue } from './PriorityQueue'
import { UnionFind } from './UnionFind'
import { WeightedAdjList, WeightedEdge } from './WeighedAdjList'

class MST {
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
   * @description 逐边遍历，权值从小到大选出边，不构成环(使用并查集判断)则是MST的边
   * @description 时间复杂度O(ElogE)
   */
  getMSTFromKruskal(): WeightedEdge[] {
    if (this.cc.CCCount > 1) return []
    const res: WeightedEdge[] = []

    this.edgeList.forEach(edge => {
      if (!this.uf.isConnected(edge.v1, edge.v2)) {
        res.push(edge)
        this.uf.union(edge.v1, edge.v2)
      }
    })

    return res
  }

  /**
   * @description 逐点遍历(V-1次)，将当前最短的**横切边(一端看过一端没看过)**添加到 mst 中
   * @description 复杂度O(VE)
   * @description todo:优化:没有必要每次循环重新扫描整个图
   * 将visited的相邻点所连的边加入优先队列扩展切分  时间复杂度O(ElogE)
   * 从优先队列中取出的横切边(一个visited一个没visited)是最小生成树的边
   */
  getMSTFromPrim(): WeightedEdge[] {
    if (this.cc.CCCount > 1) return []
    const res: WeightedEdge[] = []
    const visited = new Set<number>([0])

    // for (let steps = 1; steps < this.adjList.V; steps++) {
    //   //  找最短横切边:一个断点看过一个端点没看过
    //   let minEdge = new WeightedEdge(Infinity, Infinity, Infinity)

    //   for (let v = 0; v < this.adjList.V; v++) {
    //     if (visited.has(v)) {
    //       this.adjList.adj(v).forEach(w => {
    //         if (!visited.has(w) && this.adjList.getWeight(v, w) < minEdge.weight) {
    //           minEdge = new WeightedEdge(v, w, this.adjList.getWeight(v, w))
    //         }
    //       })
    //     }
    //   }

    //   res.push(minEdge)
    //   visited.add(minEdge.v1).add(minEdge.v2)
    // }

    // 初始化queue
    const compareFunction = (a: WeightedEdge, b: WeightedEdge) => a.weight - b.weight
    const priorityQueue = new PriorityQueue<WeightedEdge>(compareFunction)
    for (const w of this.adjList.adj(0)) {
      priorityQueue.push(new WeightedEdge(0, w, this.adjList.getWeight(0, w)))
    }

    while (priorityQueue.length) {
      if (res.length === this.adjList.V - 1) break

      const edge = priorityQueue.shift()
      // 不是合法的横切边
      if (visited.has(edge.v1) && visited.has(edge.v2)) continue
      res.push(edge)
      const next = visited.has(edge.v1) ? edge.v2 : edge.v1
      visited.add(next)
      for (const newNext of this.adjList.adj(next)) {
        if (!visited.has(newNext)) {
          priorityQueue.push(new WeightedEdge(next, newNext, this.adjList.getWeight(next, newNext)))
        }
      }
    }

    return res
  }

  static async asyncBuild(fileName: string) {
    const weightedAdjList = await WeightedAdjList.asyncBuild(fileName)
    const cc = await DFS.asyncBuild('WeightedAdjList', fileName)
    return new MST(weightedAdjList, cc)
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g.txt')
    const mst = await MST.asyncBuild(fileName)
    // console.log(mst.adjList)
    // console.log(mst.edgeList)
    // console.log(mst.getMSTFromKruskal())
    console.log(mst.getMSTFromPrim())
  }
  main()
}

export { MST }
