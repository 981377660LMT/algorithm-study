/* eslint-disable prefer-destructuring */

/**
 * csr 是 `Compressed Sparse Row` 的缩写，是一种用于存储和处理稀疏矩阵的格式。
 * @link https://github.com/atcoder/ac-library/blob/master/atcoder/internal_csr.hpp
 * @example 可用来加速树的遍历.
 */
class InternalCsr {
  /** 每个点 `i` 对应范围 `[start[i],start[i+1])`. */
  private readonly _n: number
  private readonly _start: Uint32Array
  private readonly _elist: Uint32Array

  constructor(n: number, directedEdges: [from: number, to: number][]) {
    this._n = n
    this._start = new Uint32Array(n + 1)
    this._elist = new Uint32Array(directedEdges.length)

    for (let i = 0; i < directedEdges.length; i++) {
      this._start[directedEdges[i][0] + 1]++
    }

    for (let i = 1; i <= n; i++) {
      this._start[i] += this._start[i - 1]
    }

    const counter = this._start.slice()
    for (let i = 0; i < directedEdges.length; i++) {
      const e = directedEdges[i]
      this._elist[counter[e[0]]++] = e[1]
    }
  }

  /** 遍历 `cur` 的所有邻接点. */
  enumerateNeighbors(cur: number, f: (next: number) => void): void {
    for (let i = this._start[cur]; i < this._start[cur + 1]; i++) {
      f(this._elist[i])
    }
  }

  /** 遍历图中所有边. */
  enumerateEdges(f: (u: number, v: number) => void): void {
    for (let i = 0; i < this._n; i++) {
      for (let j = this._start[i]; j < this._start[i + 1]; j++) {
        f(i, this._elist[j])
      }
    }
  }
}

export { InternalCsr }

if (require.main === module) {
  const csr = new InternalCsr(4, [
    [0, 1],
    [0, 2],
    [1, 3],
    [2, 3]
  ])

  csr.enumerateNeighbors(0, console.log)
  csr.enumerateEdges(console.log)
}
