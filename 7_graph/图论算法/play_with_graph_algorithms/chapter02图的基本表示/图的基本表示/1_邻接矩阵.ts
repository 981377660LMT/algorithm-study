import { FileReader } from './FileReader'

class AdjMatrix {
  private constructor(
    public readonly V: number,
    public readonly E: number,
    public readonly adjMatrix: number[][]
  ) {}

  static async asyncBuild(fileName: string) {
    const fileReader = await FileReader.asyncBuild(fileName)
    const V = parseInt(fileReader.fileData[0][0])
    const E = parseInt(fileReader.fileData[0][1])
    const adjMatrix = Array.from({ length: V }, () => Array(V).fill(0))

    fileReader.fileData.slice(1).forEach(([v1_, v2_]) => {
      const v1 = parseInt(v1_)
      const v2 = parseInt(v2_)
      if (v1 === v2) throw new Error('检测到自环边')
      if (adjMatrix[v1][v2] === 1) throw new Error('检测到平行边')
      adjMatrix[v1][v2] = 1
      adjMatrix[v2][v1] = 1
    })

    return new AdjMatrix(V, E, adjMatrix)
  }

  hasEdge(v: number, w: number): boolean {
    this.validateVertex(v, w)
    return this.adjMatrix[v][w] === 1
  }

  /**
   *
   * @param v 与v相邻的顶点
   */
  adj(v: number): number[] {
    this.validateVertex(v)
    const res: number[] = []
    this.adjMatrix[v].forEach((w, i) => w === 1 && res.push(i))
    return res
  }

  /**
   *
   * @param v 求顶点的度
   */
  degree(v: number): number {
    return this.adj(v).length
  }

  private validateVertex(...vArr: number[]) {
    vArr.forEach(v => {
      if (v < 0 || v >= this.V) throw new Error(`不合法的顶点序号${v}`)
    })
  }
}

if (require.main === module) {
  async function main() {
    const adjMatrix = await AdjMatrix.asyncBuild('./g.txt')
    console.table(adjMatrix.adjMatrix)
    console.log(adjMatrix)
    console.log(adjMatrix.adj(-1))
  }
  main()
}
