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

// 适用于森林的换根dp,从每个联通分量的根节点开始dp.
class RerootingForest<DpItem = number> {
  readonly adjList: number[][] = []
  private readonly _decrement: number

  constructor(n: number, decrement = 0) {
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    this.adjList = adjList
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

  /**
   * @param groupRoot 从每个联通分量的根节点开始dp.
   */
  reRooting(operation: Operation<DpItem>, groupRoot: number): Map<number, DpItem> {
    const { e, op, composition } = operation
    groupRoot -= this._decrement

    const parents = new Map<number, number>()
    const order = [groupRoot]
    const stack = [groupRoot]
    while (stack.length) {
      const cur = stack.pop()!
      for (const next of this.adjList[cur]) {
        if (next !== parents.get(cur)) {
          parents.set(next, cur)
          order.push(next)
          stack.push(next)
        }
      }
    }

    const dp1 = new Map<number, DpItem>()
    const dp2 = new Map<number, DpItem>()
    for (let i = 0; i < order.length; i++) {
      const v = order[i]
      dp1.set(v, e(v))
      dp2.set(v, e(v))
    }
    for (let i = order.length - 1; i > -1; i--) {
      const cur = order[i]
      let res = e(cur)
      for (const next of this.adjList[cur]) {
        if (next !== parents.get(cur)) {
          dp2.set(next, res)
          res = op(res, composition(dp1.get(next)!, cur, next, 0))
        }
      }

      res = e(cur)
      for (let j = this.adjList[cur].length - 1; j > -1; j--) {
        const next = this.adjList[cur][j]
        if (next !== parents.get(cur)) {
          dp2.set(next, op(res, dp2.get(next)!))
          res = op(res, composition(dp1.get(next)!, cur, next, 0))
        }
      }

      dp1.set(cur, res)
    }

    for (let i = 1; i < order.length; i++) {
      const newRoot = order[i]
      const parent = parents.get(newRoot)!
      dp2.set(newRoot, composition(op(dp2.get(newRoot)!, dp2.get(parent)!), parent, newRoot, 1))
      dp1.set(newRoot, op(dp1.get(newRoot)!, dp2.get(newRoot)!))
    }

    return dp1
  }
}

// 6294. 最大价值和与最小价值和的差值
// 求每个点作为根节点时，到叶子节点的最大点权和(不包括自身)
if (require.main === module) {
  function maxOutput(n: number, edges: number[][], price: number[]): number {
    const R = new RerootingForest(n)
    for (const [u, v] of edges) {
      R.addEdge(u, v)
    }

    const dp = R.reRooting(
      {
        e: () => 0,
        op: (childRes1, childRes2) => Math.max(childRes1, childRes2),
        composition: (fromRes, parent, cur, direction) => {
          if (direction === 0) return fromRes + price[cur] // cur => parent
          return fromRes + price[parent] // parent => cur
        }
      },
      0
    )

    return Math.max(...dp.values())
  }

  //   6
  // [[0,1],[1,2],[1,3],[3,4],[3,5]]
  // [9,8,7,6,10,5]
  console.log(
    maxOutput(
      6,
      [
        [0, 1],
        [1, 2],
        [1, 3],
        [3, 4],
        [3, 5]
      ],
      [9, 8, 7, 6, 10, 5]
    )
  )
}

export { RerootingForest }
