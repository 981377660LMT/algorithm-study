import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

class ToPoSort {
  private constructor(private readonly dfs: DFS) {}

  static async asyncBuild(fileName: string) {
    // 拓扑排序只针对有向图
    const dfs = await DFS.asyncBuild('AdjMap', fileName, true)
    return new ToPoSort(dfs)
  }

  /**
   * 拓扑排序常规解法
   */
  topoSort() {
    let hasCycle = false
    const res: number[] = []
    const queue: number[] = this.dfs.adjMap.outDegrees.filter(inDegree => inDegree === 0)

    while (queue.length) {
      const v = queue.shift()!
      res.push(v)
      this.dfs.adjMap.adj(v).forEach(w => {
        this.dfs.adjMap.outDegrees[w]--
        this.dfs.adjMap.outDegrees[w] === 0 && queue.push(w)
      })
    }

    // 有向无环图能够拓扑排序，否则有环则无解
    if (res.length < this.dfs.adjMap.V) {
      hasCycle = true
      res.splice(0)
    }

    return { res, hasCycle }
  }

  /**
   * @description 深度优先后序遍历的逆序就是拓扑排序结果
   * @description 缺点是不能有环
   */
  // anotherTopoSort() {
  //   if (this.dfs.hasLoop) throw new Error('不能有环')
  // }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../ug.txt')
    const topoSort = await ToPoSort.asyncBuild(fileName)
    console.log(topoSort.topoSort())
    // console.log(topoSort.anotherTopoSort())
  }
  main()
}

export { ToPoSort }
