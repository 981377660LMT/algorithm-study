import path from 'path'
import { FileReader } from '../../chapter02图的基本表示/图的基本表示/FileReader'

interface Graph<T> {
  V: number
  E: number
  hasEdge: (v: number, w: number) => boolean
  adj: (v: number) => number[]
  degree: (v: number) => number
  cloneAdj: () => T
  removeEdge: (v: number, w: number) => boolean
}

/**
 * @description 带权图
 */
class WeightedAdjList implements Graph<WeightedAdjList> {
  constructor(
    public readonly V: number,
    public readonly E: number,
    public readonly adjList: Map<number, number>[]
  ) {}

  static async asyncBuild(fileName: string): Promise<WeightedAdjList> {
    const fileReader = await FileReader.asyncBuild(fileName)
    const V = parseInt(fileReader.fileData[0][0])
    const E = parseInt(fileReader.fileData[0][1])
    // 第i位存储对应的点j与权值
    const adjList: Map<number, number>[] = Array.from({ length: V }, () => new Map())

    fileReader.fileData.slice(1).forEach(([v1_, v2_, weight_]) => {
      const v1 = parseInt(v1_)
      const v2 = parseInt(v2_)
      const weight = parseInt(weight_)
      if (v1 === v2) throw new Error('检测到自环边')
      if (adjList[v1].has(v2) || adjList[v2].has(v1)) throw new Error('检测到平行边')

      adjList[v1].set(v2, weight)
      adjList[v2].set(v1, weight)
    })

    return new WeightedAdjList(V, E, adjList)
  }

  hasEdge(v: number, w: number): boolean {
    this.validateVertex(v, w)
    return this.adjList[v].has(w)
  }

  /**
   *
   * @param v 获取v相邻的边
   * @returns
   */
  adj(v: number): number[] {
    this.validateVertex(v)
    return [...this.adjList[v].keys()]
  }

  degree(v: number): number {
    return this.adj(v).length
  }

  cloneAdj(): WeightedAdjList {
    return new WeightedAdjList(
      this.V,
      this.E,
      this.adjList.map(map => new Map(map))
    )
  }

  removeEdge(v: number, w: number): boolean {
    this.validateVertex(v, w)
    return this.adjList[v].delete(w) && this.adjList[w].delete(v)
  }

  getWeight(v: number, w: number): number {
    this.validateVertex(v, w)
    if (this.adjList[v].has(w)) {
      return this.adjList[v].get(w)!
    }
    throw new Error(`不存在${v}到${w}的边`)
  }

  protected validateVertex(...vArr: number[]): void {
    vArr.forEach(v => {
      if (v < 0 || v >= this.V) throw new Error(`不合法的顶点序号${v}`)
    })
  }
}

class WeightedEdge {
  constructor(public v1: number, public v2: number, public weight: number) {}
}

if (require.main === module) {
  const main = async () => {
    const fileName = path.join(__dirname, '../g.txt')
    const wam = await WeightedAdjList.asyncBuild(fileName)
    console.log(wam.adjList)
    console.log(wam.getWeight(3, 1))
  }
  main()
}

export { WeightedAdjList, WeightedEdge }
