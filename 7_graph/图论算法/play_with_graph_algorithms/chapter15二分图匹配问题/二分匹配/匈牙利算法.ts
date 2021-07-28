import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

class Hungarian {
  private constructor(private readonly dfs: DFS) {}

  static async asyncBuild(fileName: string) {
    // 拓扑排序只针对有向图
    const dfs = await DFS.asyncBuild('AdjMap', fileName, false)
    return new Hungarian(dfs)
  }

  solve() {
    let maxMathing = 0
    if (!this.dfs.isBiPartial) throw new Error('不是二分图')
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g2.txt')
    const hg = await Hungarian.asyncBuild(fileName)
    console.log(hg.solve())
  }
  main()
}

export { Hungarian }
