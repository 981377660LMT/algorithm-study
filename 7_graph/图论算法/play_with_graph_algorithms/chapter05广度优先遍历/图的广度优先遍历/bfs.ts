import { AdjMap } from '../../chapter02图的基本表示/图的基本表示/2_邻接表'

interface IBFS {
  adjMap: AdjMap
  bfs: (start?: number) => {
    pre: Map<number, number>
  }
}

/**
 * @description 主要用于寻找单源最短路径
 */
class BFS implements IBFS {
  public readonly adjMap: AdjMap
  private connect?: Map<number, number>

  private constructor(adjMap: AdjMap) {
    this.adjMap = adjMap
    this.bfs()
  }

  static async asyncBuild(fileName: string) {
    const adjMap = await AdjMap.asyncBuild(fileName)
    return new BFS(adjMap)
  }

  /**
   * @param start 从哪个顶点开始 不传则默认每个顶点
   * @description bfs入口函数
   * @description preMap记录每一步对应的前一个节点
   */
  bfs(start?: number): ReturnType<IBFS['bfs']> {
    const pre = new Map<number, number>()
    const visited = new Map<number, number>()

    for (let v = 0; v < this.adjMap.V; v++) {
      // Infinity 代表没有访问过
      visited.set(v, Infinity)
    }

    if (start !== undefined) {
      const queue = [start]
      this._bfs(start, pre, visited, queue)
    } else {
      for (let v = 0; v < this.adjMap.V; v++) {
        if (!this.isVisited(visited, v)) {
          const queue = [v]
          this._bfs(v, pre, visited, queue)
        }
      }
    }

    this.connect = visited
    return { pre }
  }

  /**
   * @description 从start到to的路径 利用pre Map 记录每个节点的pre
   */
  minPath(start: number, to: number): number[] {
    if (start === to) return []
    if (!this.isConnected(start, to)) return []

    const { pre } = this.bfs(to)
    const res: number[] = [start]
    let p = start
    while (pre.get(p) !== to) {
      p = pre.get(p)!
      res.push(p)
    }
    res.push(to)
    return res
  }

  /**
   * @description  是否在同一个连通分量
   */
  isConnected(v: number, w: number) {
    if (this.connect === undefined) {
      this.bfs()
    }
    return this.connect!.get(v) === this.connect!.get(w)
  }

  /**
   * @description 图的bfs遍历
   */
  private _bfs(
    root: number,
    pre: Map<number, number>,
    visited: Map<number, number>,
    queue: number[]
  ): void {
    while (queue.length) {
      // 表示cur属于root所在的连通分量
      const head = queue.shift()!
      visited.set(head, root)

      this.adjMap.adj(head).forEach(w => {
        if (!this.isVisited(visited, w)) {
          queue.push(w)
          visited.set(w, root)
          pre.set(w, head)
        }
      })
    }
  }

  private isVisited(visited: Map<number, number>, key: number) {
    return visited.get(key) !== Infinity
  }
}

const main = async () => {
  const bfs = await BFS.asyncBuild('../g.txt')
  console.log(bfs.adjMap)
  console.log(bfs.minPath(6, 4))
}

main()

export { BFS }
