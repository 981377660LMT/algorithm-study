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
  constructor(
    public readonly V: number,
    public readonly E: number,
    public readonly adjMap: Map<number, Set<number>>,
    public readonly directed: boolean,
    public readonly outDegrees: number[],
    public readonly inDegrees: number[]
  ) {}

  static async asyncBuild(fileName: string, directed: boolean = false) {
    const fileReader = await FileReader.asyncBuild(fileName)
    const V = parseInt(fileReader.fileData[0][0])
    const E = parseInt(fileReader.fileData[0][1])
    const adjMap = new Map<number, Set<number>>()
    const outDegrees = Array<number>(V).fill(0)
    const inDegrees = Array<number>(V).fill(0)

    fileReader.fileData.slice(1).forEach(([v1_, v2_]) => {
      const v1 = parseInt(v1_)
      const v2 = parseInt(v2_)
      if (v1 === v2) throw new Error('检测到自环边')
      if (adjMap.get(v1)?.has(v2) || adjMap.get(v2)?.has(v1)) throw new Error('检测到平行边')

      adjMap.set(v1, adjMap.get(v1)?.add(v2) || new Set([v2]))

      if (!directed) {
        adjMap.set(v2, adjMap.get(v2)?.add(v1) || new Set([v1]))
      }

      if (directed) {
        outDegrees[v1]++
        inDegrees[v2]++
      }
    })

    return new AdjMap(V, E, adjMap, directed, outDegrees, inDegrees)
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
   * @description 有向图包括出度和入度
   */
  degree(v: number): number {
    if (this.directed) {
      return this.outDegrees[v] - this.inDegrees[v]
    } else {
      return this.adj(v).length
    }
  }

  cloneAdj() {
    return new AdjMap(
      this.V,
      this.E,
      new Map(this.adjMap),
      this.directed,
      this.outDegrees,
      this.inDegrees
    )
  }

  /**
   * @description 删除成功返回true
   */
  removeEdge(v: number, w: number) {
    this.validateVertex(v, w)
    if (this.directed) {
      this.outDegrees[v]--
      this.inDegrees[w]--
    }
    return !!(this.adjMap.get(v)?.delete(w) && (this.directed || this.adjMap.get(w)?.delete(v)))
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
