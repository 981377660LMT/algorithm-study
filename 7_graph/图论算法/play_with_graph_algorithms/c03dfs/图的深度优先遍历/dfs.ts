import { AdjMap } from '../../chapter02/图的基本表示/2_邻接表'

type Color = -1 | 0 | 1
interface MetaInfo {
  hasLoop: boolean
  isBiPartical: boolean
  colors: { verticalColors: Color[]; curColor: Exclude<Color, -1> }
}

interface IDFS {
  adjMap: AdjMap
  dfs: (start?: number) => {
    path: number[]
    connect: Map<number, number>
    metaInfo: MetaInfo
  }
}

/**
 * @description 需要从每个顶点开始dfs 因为可能有多个连通分量
 * @description 应用:求无向图的联通分量(Conneted Component)个数/判断是否联通/两点间是否可达/两点路径/寻找桥和割点
 */
class DFS implements IDFS {
  public readonly adjMap: AdjMap
  private _CCCount?: number
  private connect?: Map<number, number>
  private metaInfo?: MetaInfo

  private constructor(adjMap: AdjMap) {
    this.adjMap = adjMap
    this.dfs()
  }

  static async asyncBuild(fileName: string) {
    const adjMap = await AdjMap.asyncBuild(fileName)
    return new DFS(adjMap)
  }

  /**
   * @description 连通分量
   */
  get CCCount() {
    if (this._CCCount === undefined) {
      this.dfs()
    }
    return this._CCCount!
  }

  /**
   * visited Map 转数组
   */
  get connectDetail() {
    if (this.connect === undefined) {
      this.dfs()
    }

    const connectMap = new Map<number, number[]>()
    for (const [k, v] of this.connect!) {
      const tmpArr = connectMap.get(v) || []
      tmpArr.push(k)
      connectMap.set(v, tmpArr)
    }

    return [...connectMap.values()]
  }

  get hasLoop() {
    if (this.metaInfo === undefined) {
      this.dfs()
    }
    return this.metaInfo!.hasLoop
  }

  get isBiPartial() {
    if (this.metaInfo === undefined) {
      this.dfs()
    }
    return this.metaInfo!.isBiPartical
  }

  /**
   * @param start 从哪个顶点开始 不传则默认每个顶点
   * @description dfs入口函数
   * @description 求无向图的联通情况:用visitedMap存节点与起始点对应关系，表示属于不同的连通分量/没有访问过
   * @description 求单源路径问题:只dfs起始点
   * @description 环检测:path记录走过的路，visitedSet记录看过的点，走回了之前走过的非上一个节点的节点则return有环true
   * @description 二分图检测:遍历整个图间隔染色0/1，对于访问过的节点，颜色要与相邻不同
   */
  dfs(start?: number): ReturnType<IDFS['dfs']> {
    const path: number[] = []
    const visited = new Map<number, number>()
    const verticalColors = Array<Color>(this.adjMap.V).fill(-1)
    const metaInfo: MetaInfo = {
      colors: { verticalColors, curColor: 0 },
      hasLoop: false,
      isBiPartical: true,
    }

    for (let v = 0; v < this.adjMap.V; v++) {
      // Infinity 代表没有访问过
      visited.set(v, Infinity)
    }

    if (start !== undefined) {
      this._dfs(start, start, path, visited, metaInfo)
    } else {
      for (let v = 0; v < this.adjMap.V; v++) {
        if (!this.isVisited(visited, v)) {
          this._dfs(v, v, path, visited, metaInfo)
          // 求解无向图连通分量
          this._CCCount ? this._CCCount++ : (this._CCCount = 1)
        }
      }
    }

    this.connect = visited
    this.metaInfo = metaInfo

    return { path, metaInfo, connect: visited }
  }

  /**
   * @description  是否在同一个连通分量
   */
  isConnected(v: number, w: number) {
    if (this.connect === undefined) {
      this.dfs()
    }
    return this.connect?.get(v) === this.connect?.get(w)
  }

  /**
   * @description 从start到to的路径
   */
  path(start: number, to: number): number[] {
    if (!this.isConnected(start, to)) return []
    const { path } = this.dfs(start)
    const toIndex = path.indexOf(to)
    return path.slice(0, toIndex + 1)
  }

  /**
   * @description 图的先序dfs遍历
   */
  private _dfs(
    cur: number,
    root: number,
    path: number[],
    visited: Map<number, number>,
    metaInfo: MetaInfo
  ): void {
    // 表示cur属于root所在的连通分量
    visited.set(cur, root)
    path.push(cur)
    const { verticalColors, curColor } = metaInfo.colors
    verticalColors[cur] = curColor

    this.adjMap.adj(cur).forEach(w => {
      if (!this.isVisited(visited, w)) {
        metaInfo.colors.curColor = (1 - curColor) as 0 | 1
        this._dfs(w, root, path, visited, metaInfo)
      } else {
        // 二分图中，对于访问过的节点，颜色要与相邻不同
        if (verticalColors[cur] === verticalColors[w]) {
          metaInfo.isBiPartical = false
        }
        // 检测环，走回了之前走过的非上一个节点的节点
        if (cur !== path[path.length - 1]) {
          metaInfo.hasLoop = true
        }
      }
    })
  }

  private isVisited(visited: Map<number, number>, key: number) {
    return visited.get(key) !== Infinity
  }
}

const main = async () => {
  const dfs = await DFS.asyncBuild('../g4.txt')
  // console.log(dfs.connectDetail)
  // console.log(dfs.CCCount)
  // console.log(dfs.connectDetail)
  // console.log(dfs.isConnected(1, 5))
  // console.log(dfs.dfs(1))
  // console.log(dfs.path(1, 2))
  // console.log(dfs.hasLoop)
  console.log(dfs.isBiPartial)
}

main()

export {}
