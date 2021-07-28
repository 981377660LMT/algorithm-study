import { FileReader } from './FileReader'
interface Graph<T> {
  V: number
  E: number
  hasEdge: (v: number, w: number) => boolean
  adj: (v: number) => number[]
  degree: (v: number) => number
  cloneAdj: () => T
  removeEdge: (v: number, w: number) => boolean
}

class AdjMap implements Graph<AdjMap> {
  protected constructor(
    public readonly V: number,
    public readonly E: number,
    public readonly adjMap: Map<number, Set<number>>
  ) {}

  static async asyncBuild(fileName: string) {
    const fileReader = await FileReader.asyncBuild(fileName)
    const V = parseInt(fileReader.fileData[0][0])
    const E = parseInt(fileReader.fileData[0][1])
    const adjMap = new Map<number, Set<number>>()

    fileReader.fileData.slice(1).forEach(([v1_, v2_]) => {
      const v1 = parseInt(v1_)
      const v2 = parseInt(v2_)
      if (v1 === v2) throw new Error('检测到自环边')
      if (adjMap.get(v1)?.has(v2) || adjMap.get(v2)?.has(v1)) throw new Error('检测到平行边')
      adjMap.set(v1, adjMap.get(v1)?.add(v2) || new Set([v2]))
      adjMap.set(v2, adjMap.get(v2)?.add(v1) || new Set([v1]))
    })

    return new AdjMap(V, E, adjMap)
  }

  hasEdge(v: number, w: number): boolean {
    this.validateVertex(v, w)
    return !!this.adjMap.get(v)?.has(w)
  }

  /**
   *
   * @param v 与v相邻的顶点
   */
  adj(v: number): number[] {
    this.validateVertex(v)
    return [...(this.adjMap.get(v) || [])]
  }

  /**
   *
   * @param v 求顶点的度
   */
  degree(v: number): number {
    return this.adj(v).length
  }

  cloneAdj() {
    return new AdjMap(this.V, this.E, new Map(this.adjMap))
  }

  /**
   * @description 删除成功返回true
   */
  removeEdge(v: number, w: number) {
    this.validateVertex(v, w)
    return !!(this.adjMap.get(v)?.delete(w) && this.adjMap.get(w)?.delete(v))
  }

  protected validateVertex(...vArr: number[]) {
    vArr.forEach(v => {
      if (v < 0 || v >= this.V) throw new Error(`不合法的顶点序号${v}`)
    })
  }
}

if (require.main === module) {
  async function main() {
    const adjMap = await AdjMap.asyncBuild('./g.txt')
    console.table(adjMap.adjMap)
    console.log(adjMap)
  }
  main()
}

export { AdjMap }
