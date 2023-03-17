/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-param-reassign */

interface Operation<DpItem> {
  /**
   * dp值的幺元.
   */
  e: (this: void, root: number) => DpItem

  /**
   * 合并两个dp值.
   */
  op: (this: void, childRes1: DpItem, childRes2: DpItem) => DpItem

  /**
   * 更新dp值.
   *
   * 当 direction 为 0 时，表示从子结点到父结点 (cur -> parent，fromRes表示子结点的dp值);
   * 当 direction 为 1 时，表示从父结点到子结点 (parent -> cur，fromRes表示父结点的dp值).
   */
  composition: (
    this: void,
    fromRes: DpItem,
    parent: number,
    cur: number,
    direction: 0 | 1
  ) => DpItem
}

class Rerooting<DpItem = number> {
  readonly adjList: number[][] = []
  private readonly _n: number
  private readonly _decrement: number

  constructor(n: number, decrement = 0) {
    this.adjList = Array.from({ length: n }, () => [])
    this._n = n
    this._decrement = decrement
  }

  /**
   * 添加一条无向边 u - v
   */
  addEdge(u: number, v: number): void {
    u -= this._decrement
    v -= this._decrement
    this.adjList[u].push(v)
    this.adjList[v].push(u)
  }

  reRooting(operation: Operation<DpItem>, root = 0): DpItem[] {
    const { e, op, composition } = operation
    root -= this._decrement

    const parents = new Int32Array(this._n).fill(-1)
    const order = [root]
    const stack = [root]
    while (stack.length) {
      const cur = stack.pop()!
      for (const next of this.adjList[cur]) {
        if (next !== parents[cur]) {
          parents[next] = cur
          order.push(next)
          stack.push(next)
        }
      }
    }

    const dp1 = Array.from({ length: this._n }, (_, i) => e(i))
    const dp2 = Array.from({ length: this._n }, (_, i) => e(i))
    for (let i = this._n - 1; i > -1; i--) {
      const cur = order[i]
      let res = e(cur)
      for (const next of this.adjList[cur]) {
        if (next !== parents[cur]) {
          dp2[next] = res
          res = op(res, composition(dp1[next], cur, next, 0))
        }
      }

      res = e(cur)
      for (let j = this.adjList[cur].length - 1; j > -1; j--) {
        const next = this.adjList[cur][j]
        if (next !== parents[cur]) {
          dp2[next] = op(res, dp2[next])
          res = op(res, composition(dp1[next], cur, next, 0))
        }
      }

      dp1[cur] = res
    }

    for (let i = 1; i < this._n; i++) {
      const newRoot = order[i]
      const parent = parents[newRoot]
      dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
      dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
    }

    return dp1
  }
}

// 6294. 最大价值和与最小价值和的差值
// 求每个点作为根节点时，到叶子节点的最大点权和(不包括自身)
if (require.main === module) {
  function maxOutput(n: number, edges: number[][], price: number[]): number {
    const R = new Rerooting(n)
    for (const [u, v] of edges) {
      R.addEdge(u, v)
    }

    const subSize = new Uint32Array(n).fill(1)
    // dfsForSubSize(0, -1)

    const dp = R.reRooting({
      e: () => 0,
      op: (childRes1, childRes2) => Math.max(childRes1, childRes2),
      composition: (fromRes, parent, cur, direction) => {
        if (direction === 0) return fromRes + price[cur] // cur => parent
        return fromRes + price[parent] // parent => cur
      }
    })

    return Math.max(...dp)

    function dfsForSubSize(cur: number, parent: number): void {
      R.adjList[cur].forEach(next => {
        if (next !== parent) {
          dfsForSubSize(next, cur)
          subSize[cur] += subSize[next]
        }
      })
    }
  }
}

export { Rerooting }
