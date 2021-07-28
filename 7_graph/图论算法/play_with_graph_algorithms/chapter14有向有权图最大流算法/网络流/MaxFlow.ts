import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

class MaxFlow {
  private constructor(private readonly dfs: DFS) {}

  static async asyncBuild(fileName: string) {
    // 有向有权图
    const dfs = await DFS.asyncBuild('WeightedAdjList', fileName, true)
    return new MaxFlow(dfs)
  }

  solve() {
    console.log(this.dfs.adjMap)
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../network.txt')
    const hg = await MaxFlow.asyncBuild(fileName)
    console.log(hg.solve())
  }
  main()
}

export { MaxFlow }
