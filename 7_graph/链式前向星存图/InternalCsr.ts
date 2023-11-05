/* eslint-disable prefer-destructuring */

// csr存图

// namespace internal {
//   template <class E> struct csr {
//       std::vector<int> start;
//       std::vector<E> elist;
//       explicit csr(int n, const std::vector<std::pair<int, E>>& edges)
//           : start(n + 1), elist(edges.size()) {
//           for (auto e : edges) {
//               start[e.first + 1]++;
//           }
//           for (int i = 1; i <= n; i++) {
//               start[i] += start[i - 1];
//           }
//           auto counter = start;
//           for (auto e : edges) {
//               elist[counter[e.first]++] = e.second;
//           }
//       }
//   };
// }

/**
 * @link https://github.com/atcoder/ac-library/blob/master/atcoder/internal_csr.hpp
 * @example 可用来加速树的遍历.
 */
class InternalCsr {
  /** 每个点 `i` 对应范围 `[start[i],start[i+1])`. */
  private readonly _start: Uint32Array
  private readonly _elist: Uint32Array

  constructor(n: number, edges: [number, number][]) {
    this._start = new Uint32Array(n + 1)
    this._elist = new Uint32Array(edges.length)

    for (let i = 0; i < edges.length; i++) {
      this._start[edges[i][0] + 1]++
    }

    for (let i = 1; i <= n; i++) {
      this._start[i] += this._start[i - 1]
    }

    for (let i = 0; i < edges.length; i++) {
      const e = edges[i]
      this._elist[this._start[e[0]]++] = e[1]
    }

    console.log(this._start, this._elist)
  }

  enumerateNeighbor(cur: number, f: (next: number) => void): void {
    for (let i = this._start[cur]; i < this._start[cur + 1]; i++) {
      f(this._elist[i])
    }
  }

  enumerateEdge(f: (u: number, v: number) => void): void {}
}

export { InternalCsr }

if (require.main === module) {
  const csr = new InternalCsr(4, [
    [0, 1],
    [0, 2],
    [1, 3],
    [2, 3]
  ])

  // for (int i = g.start[v]; i < g.start[v + 1]; i++) {
  //   auto e = g.elist[i];

  csr.enumerateNeighbor(0, console.log)
}
