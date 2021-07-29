import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'
/**
 * @description 求解有向图强连通分量
 */
class SCC {
  constructor(private readonly dfs: DFS) {}

  static async asyncBuild(fileName: string) {
    const dfs = await DFS.asyncBuild('AdjMap', fileName, true)
    return new SCC(dfs)
  }

  /**
   * @description 将所有强连通分量看做一个点，得到的有向图一定是DAG
   */
  kosaraju() {
    console.log(this.dfs.adjMap)
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../ug4.txt')
    const scc = await SCC.asyncBuild(fileName)
    console.log(scc.kosaraju())
  }
  main()
}

export { SCC }
