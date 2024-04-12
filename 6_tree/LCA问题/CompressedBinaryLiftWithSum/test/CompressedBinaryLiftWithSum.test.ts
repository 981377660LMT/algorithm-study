import { CompressedBinaryLiftWithSum } from '../CompressedBinaryLiftWithSum'

describe('CompressedBinaryLiftWithSum.ts', () => {
  let n: number
  let edges: [number, number][]
  let tree: number[][]
  let values: number[]
  let bl: CompressedBinaryLiftWithSum

  beforeEach(() => {
    //          0
    //        /   \
    //       1     2
    //      / \     \
    //     3   4     5
    //         /
    //        6
    n = 7
    edges = [
      [0, 1],
      [0, 2],
      [1, 3],
      [1, 4],
      [2, 5],
      [4, 6]
    ]
    tree = Array(n)
    for (let i = 0; i < n; i++) tree[i] = []
    edges.forEach(([u, v]) => {
      tree[u].push(v)
      tree[v].push(u)
    })
    values = [1, 1, 2, 3, 4, 5, 6]
    bl = new CompressedBinaryLiftWithSum(tree, i => values[i], {
      e: () => 0,
      op: (a, b) => a + b
    })
  })

  it('should support firstTrueWithSum', () => {
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 0, true)).toEqual({ node: 6, sum: 0 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 6, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 10, true)).toEqual({ node: 1, sum: 10 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 11, true)).toEqual({ node: 0, sum: 11 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 15, true)).toEqual({ node: -1, sum: 11 })

    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 0, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 6, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 10, false)).toEqual({ node: 4, sum: 10 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 11, false)).toEqual({ node: 1, sum: 11 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 12, false)).toEqual({ node: 0, sum: 12 })
    expect(bl.firstTrueWithSum(6, (i, sum) => sum >= 15, false)).toEqual({ node: -1, sum: 12 })

    expect(bl.firstTrueWithSum(6, i => bl.depth[i] <= 1, true)).toEqual({ node: 1, sum: 10 })
    expect(bl.firstTrueWithSum(6, i => bl.depth[i] <= 1, false)).toEqual({ node: 1, sum: 11 })
  })

  it('should support lastTrueWithSum', () => {
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= -1, true)).toEqual({ node: -1, sum: 0 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 0, true)).toEqual({ node: 6, sum: 0 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 5, true)).toEqual({ node: 6, sum: 0 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 6, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 7, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 10, true)).toEqual({ node: 1, sum: 10 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 11, true)).toEqual({ node: 0, sum: 11 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 12, true)).toEqual({ node: 0, sum: 11 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 13, true)).toEqual({ node: 0, sum: 11 })

    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= -1, false)).toEqual({ node: -1, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 0, false)).toEqual({ node: -1, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 5, false)).toEqual({ node: -1, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 6, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 7, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 10, false)).toEqual({ node: 4, sum: 10 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 11, false)).toEqual({ node: 1, sum: 11 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 12, false)).toEqual({ node: 0, sum: 12 })
    expect(bl.lastTrueWithSum(6, (i, sum) => sum <= 13, false)).toEqual({ node: 0, sum: 12 })

    expect(bl.lastTrueWithSum(6, i => bl.depth[i] >= 2, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.lastTrueWithSum(6, i => bl.depth[i] >= 2, false)).toEqual({ node: 4, sum: 10 })
  })

  it('should support upToDepthWithSum', () => {
    expect(bl.upToDepthWithSum(6, 1, true)).toEqual({ node: 1, sum: 10 })
    expect(bl.upToDepthWithSum(6, 1, false)).toEqual({ node: 1, sum: 11 })
    expect(bl.upToDepthWithSum(6, 2, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.upToDepthWithSum(6, 2, false)).toEqual({ node: 4, sum: 10 })
    expect(bl.upToDepthWithSum(6, 3, true)).toEqual({ node: 6, sum: 0 })
    expect(bl.upToDepthWithSum(6, 3, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.upToDepthWithSum(6, 4, true)).toEqual({ node: -1, sum: 0 })
    expect(bl.upToDepthWithSum(6, 4, false)).toEqual({ node: -1, sum: 0 })
  })

  it('should support kthAncestorWithSum', () => {
    expect(bl.kthAncestorWithSum(6, 0, true)).toEqual({ node: 6, sum: 0 })
    expect(bl.kthAncestorWithSum(6, 0, false)).toEqual({ node: 6, sum: 6 })
    expect(bl.kthAncestorWithSum(6, 1, true)).toEqual({ node: 4, sum: 6 })
    expect(bl.kthAncestorWithSum(6, 1, false)).toEqual({ node: 4, sum: 10 })
    expect(bl.kthAncestorWithSum(6, 2, true)).toEqual({ node: 1, sum: 10 })
    expect(bl.kthAncestorWithSum(6, 2, false)).toEqual({ node: 1, sum: 11 })
    expect(bl.kthAncestorWithSum(6, 3, true)).toEqual({ node: 0, sum: 11 })
    expect(bl.kthAncestorWithSum(6, 3, false)).toEqual({ node: 0, sum: 12 })
    expect(bl.kthAncestorWithSum(6, 4, true)).toEqual({ node: -1, sum: 0 })
    expect(bl.kthAncestorWithSum(6, 4, false)).toEqual({ node: -1, sum: 0 })
  })

  it('should support lcaWithSum', () => {
    const weightSum = (u: number, v: number, isEdge: boolean) => {
      if (bl.depth[u] < bl.depth[v]) [u, v] = [v, u]
      let sum = 0
      while (bl.depth[u] > bl.depth[v]) {
        sum += values[u]
        u = bl.parent[u]
      }
      while (u !== v) {
        sum += values[u] + values[v]
        u = bl.parent[u]
        v = bl.parent[v]
      }
      return isEdge ? sum : sum + values[u]
    }

    for (let i = 0; i < n; i++) {
      for (let j = 0; j < n; j++) {
        const lca = bl.lca(i, j)
        expect(bl.lcaWithSum(i, j, true)).toEqual({
          node: lca,
          sum: weightSum(i, j, true)
        })
        expect(bl.lcaWithSum(i, j, false)).toEqual({
          node: lca,
          sum: weightSum(i, j, false)
        })
      }
    }
  })

  it('should support inSubtree', () => {
    expect(bl.inSubtree(6, 6)).toBeTruthy()
    expect(bl.inSubtree(6, 4)).toBeTruthy()
    expect(bl.inSubtree(6, 1)).toBeTruthy()
    expect(bl.inSubtree(6, 0)).toBeTruthy()
    expect(bl.inSubtree(6, 2)).toBeFalsy()
    expect(bl.inSubtree(6, 3)).toBeFalsy()
    expect(bl.inSubtree(6, 5)).toBeFalsy()
  })
})
