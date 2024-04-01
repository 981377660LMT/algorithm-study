import { CompressedBinaryLiftWithSum } from '../../CompressedBinaryLiftWithSum/CompressedBinaryLiftWithSum'
import { ITreePath, TreePath } from '../TreePath'

describe('TreePath.ts', () => {
  let n: number
  let edges: [number, number][]
  let tree: number[][]
  let values: number[]
  let bl: CompressedBinaryLiftWithSum
  const createPath = (from: number, to: number): ITreePath => {
    return new TreePath(from, to, { depth: bl.depth, kthAncestorFn: bl.kthAncestor, lcaFn: bl.lca })
  }

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

  it('should support kthNodeOnPath', () => {
    const path1 = createPath(3, 6)
    expect(path1.kthNodeOnPath(0)).toBe(3)
    expect(path1.kthNodeOnPath(1)).toBe(1)
    expect(path1.kthNodeOnPath(2)).toBe(4)
    expect(path1.kthNodeOnPath(3)).toBe(6)
    expect(path1.kthNodeOnPath(4)).toBe(-1)
    const path2 = createPath(6, 3)
    expect(path2.kthNodeOnPath(0)).toBe(6)
    expect(path2.kthNodeOnPath(1)).toBe(4)
    expect(path2.kthNodeOnPath(2)).toBe(1)
    expect(path2.kthNodeOnPath(3)).toBe(3)
    expect(path2.kthNodeOnPath(4)).toBe(-1)
    const path3 = createPath(3, 3)
    expect(path3.kthNodeOnPath(0)).toBe(3)
    expect(path3.kthNodeOnPath(1)).toBe(-1)
    const path4 = createPath(5, 6)
    expect(path4.kthNodeOnPath(0)).toBe(5)
    expect(path4.kthNodeOnPath(1)).toBe(2)
    expect(path4.kthNodeOnPath(2)).toBe(0)
    expect(path4.kthNodeOnPath(3)).toBe(1)
    expect(path4.kthNodeOnPath(4)).toBe(4)
    expect(path4.kthNodeOnPath(5)).toBe(6)
    expect(path4.kthNodeOnPath(6)).toBe(-1)
  })

  it('should support onPath', () => {
    const path1 = createPath(3, 6)
    expect(path1.onPath(3)).toBeTruthy()
    expect(path1.onPath(1)).toBeTruthy()
    expect(path1.onPath(4)).toBeTruthy()
    expect(path1.onPath(6)).toBeTruthy()
    expect(path1.onPath(0)).toBeFalsy()
    expect(path1.onPath(2)).toBeFalsy()
    expect(path1.onPath(5)).toBeFalsy()

    const path2 = createPath(6, 3)
    expect(path2.onPath(6)).toBeTruthy()
    expect(path2.onPath(4)).toBeTruthy()
    expect(path2.onPath(1)).toBeTruthy()
    expect(path2.onPath(3)).toBeTruthy()
    expect(path2.onPath(0)).toBeFalsy()
    expect(path2.onPath(2)).toBeFalsy()
    expect(path2.onPath(5)).toBeFalsy()

    const path3 = createPath(3, 3)
    expect(path3.onPath(3)).toBeTruthy()

    const path4 = createPath(5, 6)
    expect(path4.onPath(5)).toBeTruthy()
    expect(path4.onPath(2)).toBeTruthy()
    expect(path4.onPath(0)).toBeTruthy()
    expect(path4.onPath(1)).toBeTruthy()
    expect(path4.onPath(4)).toBeTruthy()
    expect(path4.onPath(6)).toBeTruthy()
    expect(path4.onPath(3)).toBeFalsy()
  })

  it('should support hasIntersection', () => {
    expect(createPath(3, 5).hasIntersection(createPath(1, 6))).toBeTruthy()
    expect(createPath(0, 5).hasIntersection(createPath(1, 6))).toBeFalsy()
  })

  it('should support getIntersection', () => {
    const res1 = createPath(3, 5).getIntersection(createPath(1, 6))
    expect(res1).toEqual({ p1: 1, p2: 1 })
    const res2 = createPath(3, 6).getIntersection(createPath(1, 4))
    expect(res2).toEqual({ p1: 4, p2: 1 })
    const res3 = createPath(0, 5).getIntersection(createPath(1, 6))
    expect(res3).toBeUndefined()
  })

  it('should support countIntersection', () => {
    expect(createPath(3, 5).countIntersection(createPath(1, 6))).toBe(1)
    expect(createPath(0, 5).countIntersection(createPath(1, 6))).toBe(0)
    expect(createPath(3, 6).countIntersection(createPath(1, 4))).toBe(2)
    expect(createPath(3, 3).countIntersection(createPath(3, 3))).toBe(1)
    expect(createPath(5, 6).countIntersection(createPath(4, 2))).toBe(4)
  })
})
