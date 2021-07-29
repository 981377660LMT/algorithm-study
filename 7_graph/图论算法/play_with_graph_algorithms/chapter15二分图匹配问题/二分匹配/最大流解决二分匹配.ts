import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
import { WeightedAdjList } from '../../chapter11无向带权图最小生成树/带权图/WeighedAdjList'
import { MaxFlow } from '../../chapter14有向有权图最大流算法/网络流/MaxFlow'

class BipartiteMaxGrow {
  private readonly dfs: DFS
  private readonly maxflow: MaxFlow
  constructor(dfs: DFS) {
    this.dfs = dfs
    this.maxflow = this.buildMaxflowFromDFS(dfs)
  }

  static async asyncBuild(fileName: string) {
    // 二分图是无权图
    const dfs = await DFS.asyncBuild('AdjMap', fileName, false)
    return new BipartiteMaxGrow(dfs)
  }

  /**
   *
   * @returns 二分图最大匹配数
   */
  get maxMatching(): number {
    return this.maxflow.getMaxFlow(this.dfs.adjMap.V, this.dfs.adjMap.V + 1)
  }

  get isPerfectMatching() {
    return this.maxMatching * 2 === this.dfs.adjMap.V
  }

  private buildMaxflowFromDFS(dfs: DFS) {
    if (dfs.isBiPartial === false) {
      throw new Error('不是二分图')
    }

    const netWork = new WeightedAdjList(
      dfs.adjMap.V + 2,
      dfs.adjMap.E,
      Array.from({ length: dfs.adjMap.V + 2 }, () => new Map()),
      true,
      Array(dfs.adjMap.V + 2).fill(0),
      Array(dfs.adjMap.V + 2).fill(0)
    )

    // 原点s是V 汇点t是V+1
    // 颜色为0的点与源点连接 颜色为1的点与汇点连接
    // 注意原图无向每条边只能算一次 要加上v<w限制
    const s = dfs.adjMap.V
    const t = dfs.adjMap.V + 1
    for (let v = 0; v < dfs.adjMap.V; v++) {
      if (dfs.colors[v] === 0) {
        netWork.addEdge(s, v, 1)
      } else if (dfs.colors[v] === 1) {
        netWork.addEdge(v, t, 1)
      }
      for (const w of dfs.adjMap.adj(v)) {
        if (v < w) {
          if (dfs.colors[v] === 0) {
            netWork.addEdge(v, w, 1)
          } else if (dfs.colors[v] === 1) {
            netWork.addEdge(w, v, 1)
          }
        }
      }
    }

    return new MaxFlow(netWork)
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g2.txt')
    const bmg = await BipartiteMaxGrow.asyncBuild(fileName)
    console.log(bmg.maxMatching)
    console.log(bmg.isPerfectMatching)
  }
  main()
}

export { BipartiteMaxGrow }
