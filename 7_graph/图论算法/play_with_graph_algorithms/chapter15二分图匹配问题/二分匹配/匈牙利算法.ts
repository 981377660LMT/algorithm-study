import path from 'path'
import { DFS } from '../../chapter04深度优先遍历应用/c03dfs/图的深度优先遍历/dfs'

/**
 * @description 对左侧每一个尚未匹配的点，不断地寻找可以增广的交替路
 * @description 可以利用 BFS/DFS 寻找增广路径
 */
class Hungarian {
  private readonly dfs: DFS
  public readonly matching: number[]
  public maxMatching: number
  private constructor(dfs: DFS) {
    this.dfs = dfs
    this.matching = Array(dfs.adjMap.V).fill(-1)
    this.maxMatching = 0
    this.solve()
  }

  static async asyncBuild(fileName: string) {
    const dfs = await DFS.asyncBuild('AdjMap', fileName, false)
    if (!dfs.isBiPartial) throw new Error('不是二分图')

    return new Hungarian(dfs)
  }

  solve() {
    const colors = this.dfs.colors
    // dfs使用
    const visited = new Set<number>()

    // 颜色为0的点在左边 颜色为1的点在右边
    for (let v = 0; v < this.dfs.adjMap.V; v++) {
      // 从未匹配的左侧点开始寻找增广路径
      if (colors[v] === 0 && this.matching[v] === -1) {
        // this.bfsFind(v)
        this.dfsFind(v, visited)
      }
    }
  }

  /**
   * @param 从左边的v开始寻找增广路径
   * @description BFS 队列中只存储左边的点
   */
  private bfsFind(v: number) {
    const queue: number[] = []
    // 记录路径并充当visited的作用
    const pre: number[] = Array(this.dfs.adjMap.V).fill(-1)
    queue.push(v)
    pre[v] = v

    while (queue.length) {
      const cur = queue.shift()!

      for (const next of this.dfs.adjMap.adj(cur)) {
        if (pre[next] === -1) {
          // next已经匹配
          // this.matching[next]是左侧的点
          if (this.matching[next] !== -1) {
            pre[next] = cur
            pre[this.matching[next]] = next
            queue.push(this.matching[next])
          } else {
            // 找到了v到next的增广路径augPath
            this.maxMatching++
            pre[next] = cur
            const augPath = this.getAugPath(pre, v, next)
            // 更新增广路径上的匹配关系(偶数个点相邻两两重新配对)
            for (let i = 0; i < augPath.length; i += 2) {
              this.matching[augPath[i]] = this.matching[augPath[i + 1]]
              this.matching[augPath[i + 1]] = this.matching[augPath[i]]
            }
            break
          }
        }
      }
    }
  }

  private dfsFind(v: number, visited: Set<number>) {
    visited.add(v)
    for (const w of this.dfs.adjMap.adj(v)) {
      if (!visited.has(w)) {
        visited.add(w)
        if (this.matching[w] === -1) {
          // 找到增广路径
          this.matching[v] = w
          this.matching[w] = v
          this.maxMatching++
          return
        } else this.dfsFind(this.matching[w], visited)
      }
    }
  }

  private getAugPath(pre: number[], start: number, end: number) {
    const res: number[] = []
    let p = end

    while (pre[p] !== start) {
      res.push(p)
      p = pre[p]
    }

    res.push(start)
    return res.reverse()
  }
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g2.txt')
    const hg = await Hungarian.asyncBuild(fileName)
    console.log(hg.maxMatching)
  }
  main()
}

export { Hungarian }
