import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
import { WeightedAdjList } from '../../chapter11无向带权图最小生成树/带权图/WeighedAdjList'

class MaxFlow {
  public readonly netWork: WeightedAdjList
  public readonly residualGraph: WeightedAdjList
  private clonedResidualGraph: WeightedAdjList
  constructor(netWork: WeightedAdjList) {
    this.netWork = netWork
    this.residualGraph = this.buildResidualGraphFromNetwork(netWork)
    this.clonedResidualGraph = this.cloneResidualGraph()
  }

  static async asyncBuild(fileName: string) {
    // 有向带权图
    const dfs = await DFS.asyncBuild('WeightedAdjList', fileName, true)
    const network = dfs.adjMap as WeightedAdjList
    return new MaxFlow(network)
  }

  /**
   *
   * @param network 输入的原始图
   * @returns 残差图
   */
  private buildResidualGraphFromNetwork(network: WeightedAdjList) {
    // 初始化残量图
    const residualGraph = new WeightedAdjList(
      network.V,
      network.E,
      Array.from({ length: network.V }, () => new Map()),
      true,
      Array(network.V).fill(0),
      Array(network.V).fill(0)
    )

    for (let v = 0; v < network.V; v++) {
      for (const w of network.adj(v)) {
        const volumn = network.getWeight(v, w)
        // 正向流量
        residualGraph.addEdge(v, w, volumn)
        // 反向流量
        residualGraph.addEdge(w, v, 0)
      }
    }

    return residualGraph
  }
  /**
   * @param {number} s 源点
   * @param {number} t 汇点
   * @description EK算法
   * 1.构建残量图
   * 2.在残量图中使用bfs寻找增广路径
   * 3.计算增广路径上的最小值
   * 4.根据增广路径更新残量图
   * @description 残量图的正向加反向等于容量
   */
  getMaxFlow(s: number, t: number) {
    if (s === t || this.netWork.V < 2) throw new Error('不合法的输入')
    this.clonedResidualGraph = this.cloneResidualGraph()

    let maxFlow = 0

    while (true) {
      const augPath = this.getArgumentingPath(s, t)
      if (augPath.length === 0) break

      //  计算增广路径上的最小值
      let min = Infinity
      for (let i = 1; i < augPath.length; i++) {
        const v = augPath[i - 1]
        const w = augPath[i]
        min = Math.min(min, this.clonedResidualGraph.getWeight(v, w))
      }
      // console.log(min)
      maxFlow += min

      // 根据增广路径更新残量图
      for (let i = 1; i < augPath.length; i++) {
        const v = augPath[i - 1]
        const w = augPath[i]
        this.clonedResidualGraph.setWeight(v, w, this.clonedResidualGraph.getWeight(v, w) - min)
        this.clonedResidualGraph.setWeight(w, v, this.clonedResidualGraph.getWeight(w, v) + min)
      }
    }

    return maxFlow
  }

  private cloneResidualGraph() {
    return new WeightedAdjList(
      this.residualGraph.V,
      this.residualGraph.E,
      this.residualGraph.adjList.map(map => new Map(map)),
      this.residualGraph.directed,
      this.residualGraph.outDegrees.slice(),
      this.residualGraph.inDegrees.slice()
    )
  }

  /**
   * @description bfs寻找增广路径
   */
  private getArgumentingPath(s: number, t: number): number[] {
    const res: number[] = []
    const queue: number[] = [s]
    const visited = new Set<number>([s])
    const pre: number[] = Array(this.netWork.V).fill(-1)

    while (queue.length) {
      const v = queue.shift()!
      for (const w of this.clonedResidualGraph.adj(v)) {
        if (!visited.has(w) && this.clonedResidualGraph.getWeight(v, w) > 0) {
          visited.add(w)
          queue.push(w)
          pre[w] = v
        }
      }
    }

    // 没找到增广路径
    if (pre[t] === -1) return []

    let p = t
    while (p !== s) {
      res.push(p)
      p = pre[p]
    }
    res.push(s)

    return res.reverse()
  }
}

if (require.main === module) {
  const main = async () => {
    // const fileName = path.join(__dirname, '../network2.txt')
    const fileName = path.join(__dirname, '../baseball.txt')
    const mf = await MaxFlow.asyncBuild(fileName)
    console.log(mf.netWork)
    console.log(mf.residualGraph)
    console.log(mf.getMaxFlow(0, 10))
    console.log(mf.residualGraph)
  }
  main()
}

export { MaxFlow }
