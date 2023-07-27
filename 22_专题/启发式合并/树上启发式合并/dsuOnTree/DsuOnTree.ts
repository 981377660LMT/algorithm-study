/* eslint-disable no-inner-declarations */
/* eslint-disable prefer-destructuring */

/**
 * `O(nlogn)` 静态查询每个子树内的信息, 空间复杂度优于启发式合并.
 * @see {@link https://blog.csdn.net/qq_43472263/article/details/104150940}
 */
class DsuOnTree {
  private readonly _tree: number[][]
  private readonly _root: number
  private readonly _subSize: Uint32Array
  private readonly _euler: Uint32Array
  private readonly _down: Uint32Array
  private readonly _up: Uint32Array
  private _order = 0

  constructor(n: number, tree: number[][], root = 0) {
    this._tree = tree
    this._root = root
    this._subSize = new Uint32Array(n)
    this._euler = new Uint32Array(n)
    this._down = new Uint32Array(n)
    this._up = new Uint32Array(n)

    this._dfs1(root, -1)
    this._dfs2(root, -1)
  }

  /**
   * @param add 添加root处的贡献.
   * @param remove 移除root处的贡献.
   * @param query 查询root的子树的贡献并更新答案.
   * @param reset 退出轻儿子时的回调函数.
   */
  run(
    add: (root: number) => void,
    remove: (root: number) => void,
    query: (root: number) => void,
    reset?: () => void
  ) {
    const dsu = (cur: number, pre: number, keep: boolean): void => {
      const nexts = this._tree[cur]
      for (let i = 1; i < nexts.length; i++) {
        const next = nexts[i]
        if (next !== pre) {
          dsu(next, cur, false)
        }
      }

      if (this._subSize[cur] !== 1) {
        dsu(nexts[0], cur, true)
      }

      if (this._subSize[cur] !== 1) {
        for (let i = this._up[nexts[0]]; i < this._up[cur]; i++) {
          add(this._euler[i])
        }
      }

      add(cur)
      query(cur)
      if (!keep) {
        for (let i = this._down[cur]; i < this._up[cur]; i++) {
          remove(this._euler[i])
        }
        reset && reset()
      }
    }

    dsu(this._root, -1, false)
  }

  private _dfs1(cur: number, pre: number): number {
    this._subSize[cur] = 1
    const nexts = this._tree[cur]
    if (nexts.length >= 2 && nexts[0] === pre) {
      const tmp = nexts[0]
      nexts[0] = nexts[1]
      nexts[1] = tmp
    }
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      this._subSize[cur] += this._dfs1(next, cur)
      if (this._subSize[next] > this._subSize[nexts[0]]) {
        const tmp = nexts[0]
        nexts[0] = next
        nexts[i] = tmp
      }
    }
    return this._subSize[cur]
  }

  private _dfs2(cur: number, pre: number): void {
    this._euler[this._order] = cur
    this._down[cur] = this._order++
    const nexts = this._tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next === pre) continue
      this._dfs2(next, cur)
    }
    this._up[cur] = this._order
  }
}

export { DsuOnTree }

if (require.main === module) {
  // 2003. 每棵子树内缺失的最小基因值
  // https://leetcode.cn/problems/smallest-missing-genetic-value-in-each-subtree/
  function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
    const n = parents.length
    const adjList: number[][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    for (let i = 1; i < n; i++) adjList[parents[i]].push(i)

    const dsu = new DsuOnTree(n, adjList)

    const res: number[] = Array(n).fill(1)
    const counter = new Map<number, number>()
    let mex = 1
    dsu.run(add, remove, query)
    return res

    function add(root: number): void {
      const num = nums[root]
      counter.set(num, (counter.get(num) || 0) + 1)
      while (counter.get(mex)) mex++
    }
    function remove(root: number): void {
      const num = nums[root]
      const count = counter.get(num)!
      counter.set(num, count - 1)
      if (count === 1) mex = Math.min(mex, num)
    }
    function query(root: number): void {
      res[root] = mex
    }
  }
}
