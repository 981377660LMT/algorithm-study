import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../ug.txt')
    const dfs = await DFS.asyncBuild('AdjMap', fileName, true)
    console.log(dfs.adjMap)
    // console.log(dfs.CCCount)
    // console.log(dfs.connectDetail)
    // console.log(dfs.isConnected(1, 5))
    // console.log(dfs.dfs(1))
    // console.log(dfs.path(1, 2))
    console.log(dfs.hasLoop)
    console.log(dfs.adjMap.degree(1))
    // console.log(dfs.isBiPartial)
  }
  main()
}

export {}
