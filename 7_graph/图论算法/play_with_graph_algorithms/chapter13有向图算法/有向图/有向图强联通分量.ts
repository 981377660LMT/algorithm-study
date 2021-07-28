import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

// todo kosaraju算法
class SCC {
  private constructor(private readonly dfs: DFS) {}

  static async asyncBuild(fileName: string) {
    // 拓扑排序只针对有向图
    const dfs = await DFS.asyncBuild('AdjMap', fileName, true)
    return new SCC(dfs)
  }

  /**
   * @description 将所有强连通分量看做一个点，得到的有向图一定是DAG
   */
  cc() {}
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../ug4.txt')
    const cc = await SCC.asyncBuild(fileName)

    // console.log(topoSort.anotherTopoSort())
  }
  main()
}

export { SCC }
