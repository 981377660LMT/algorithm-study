import { AdjMap } from '../chapter02图的基本表示/图的基本表示/2_邻接表'

/**
 * @link https://nicodechal.github.io/2019/10/14/bridge-finding-algorithm/
 * @description
 * 1.首先对图进行一次 DFS 搜索，并对节点按照遍历到的顺序进行编号存放在 dfn 数组 dfn[i] 是第 i 个节点的遍历顺序。
 * 2.然后，对于每个节点，计算其在不经过其父节点时能够达到的编号最小值,存放在 low 数组中 ,low[i] 是第 i 个节点不走父节点能到的最小的 dfn 值
 * 3.最后，对于每条边 ab (a 节点先被遍历即b是a的nexxt)，如果有 dfn[a] < low[b] 则 ab 边为桥
 * @summary 使用两个数组记录每个节点的遍历编号 dfn 和不经过父节点的可达最小遍历编号 low 再通过比较边的两端点的 order 和 low 值了解该边是否为桥。
 * Tarjan算法从图的任意顶点进行DFS都可以得出割点集和割边集
 */
class FindBridge {
  public readonly adjMap: AdjMap
  private steps: number

  private constructor(adjMap: AdjMap) {
    this.adjMap = adjMap
    this.steps = 0
  }

  static async asyncBuild(fileName: string) {
    const adjMap = await AdjMap.asyncBuild(fileName)
    return new FindBridge(adjMap)
  }

  findBridge() {
    const visited = new Set<number>()
    const order = new Map<number, number>()
    const lower = new Map<number, number>()
    const bridge: [number, number][] = []

    for (let v = 0; v < this.adjMap.V; v++) {
      !visited.has(v) && this.dfs(v, v, visited, order, lower, bridge)
    }
    console.log(order, lower)
    return bridge
  }

  private dfs(
    cur: number,
    parent: number,
    visited: Set<number>,
    order: Map<number, number>,
    lower: Map<number, number>,
    bridge: [number, number][],
    steps: number = this.steps,
    path: number[] = []
  ) {
    visited.add(cur)
    // 初始值
    order.set(cur, steps)
    lower.set(cur, steps)
    path.push(cur)
    this.steps++
    for (const next of this.adjMap.adj(cur)) {
      if (!visited.has(next)) {
        this.dfs(next, cur, visited, order, lower, bridge, this.steps, path)
        // 倒退回溯 到这里 lower[next]已经算好了
        // 当前节点的 lower 值是子节点 lower 值的最小值
        lower.set(cur, Math.min(lower.get(cur)!, lower.get(next)!))
        // next的low大于cur的order说明next在前一组，是桥)
        if (lower.get(next)! > order.get(cur)!) {
          bridge.push([cur, next])
        }
      } else if (next !== parent) {
        lower.set(cur, Math.min(lower.get(cur)!, lower.get(next)!))
      }
    }
  }
}

const main = async () => {
  const fb = await FindBridge.asyncBuild('./g.txt')
  console.log(fb.findBridge())
}
main()

export { FindBridge }
