import { AdjMap } from '../chapter02图的基本表示/图的基本表示/2_邻接表'

/**
 * @description 就是回溯法 O(n!)
 * 使用记忆化搜索可以优化到O(n*2^n)
 */
class HamiltonLoop {
  public readonly adjMap: AdjMap

  private constructor(adjMap: AdjMap) {
    this.adjMap = adjMap
  }

  static async asyncBuild(fileName: string) {
    const adjMap = await AdjMap.asyncBuild(fileName)
    return new HamiltonLoop(adjMap)
  }

  dfs(start: number) {
    const visited = new Set<number>([start])
    const path: number[] = [start]
    const res: number[] = []

    this._dfs(start, start, visited, path, res)

    return res
  }

  private _dfs(cur: number, root: number, visited: Set<number>, path: number[], res: number[]) {
    for (const next of this.adjMap.adj(cur)) {
      if (!visited.has(next)) {
        visited.add(next)
        path.push(next)
        this._dfs(next, root, visited, path, res)
      } else if (next === root && this.allVisited(visited)) {
        res.push(...path)
        return
      }
    }

    // 回溯 注意这个位置
    visited.delete(cur)
    path.pop()
  }

  private allVisited(visited: Set<number>) {
    return visited.size === this.adjMap.V
  }
}

const main = async () => {
  const fb = await HamiltonLoop.asyncBuild('./g2.txt')
  console.log(fb.adjMap)
  console.log(fb.dfs(0))
}
main()

export { HamiltonLoop }
