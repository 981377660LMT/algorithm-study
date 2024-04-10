import { CompressedBinaryLift } from '../../CompressedBinaryLift/CompressedBinaryLift'
import { splitPath, splitPathByJump } from '../splitPath'
import { ITreePath } from '../TreePath'

describe('splitPath.ts', () => {
  let n: number
  let edges: [number, number][]
  let tree: number[][]
  let bl: CompressedBinaryLift

  const createPath = (from: number, to: number): ITreePath => {
    const dep = {
      splitFn: (
        separator: number
      ): { path1: ITreePath | undefined; path2: ITreePath | undefined } => {
        const { from1, to1, from2, to2 } =
          Math.random() > 0.5 // 1
            ? splitPathByJump(from, to, separator, bl.jump.bind(bl))
            : splitPath(from, to, separator, {
                depth: bl.depth,
                kthAncestorFn: bl.kthAncestor.bind(bl),
                lcaFn: bl.lca.bind(bl)
              })
        let path1: ITreePath | undefined = undefined
        let path2: ITreePath | undefined = undefined
        if (from1 !== undefined && to1 !== undefined) {
          path1 = new TreePathAdapter(from1, to1, dep)
        }
        if (from2 !== undefined && to2 !== undefined) {
          path2 = new TreePathAdapter(from2, to2, dep)
        }
        return { path1, path2 }
      }
    }

    return new TreePathAdapter(from, to, dep)
  }

  class TreePathAdapter implements ITreePath {
    readonly from: number
    readonly to: number
    private readonly _splitFn: (separator: number) => {
      path1: ITreePath | undefined
      path2: ITreePath | undefined
    }

    constructor(
      from: number,
      to: number,
      dependencies: {
        splitFn: ITreePath['split']
      }
    ) {
      this.from = from
      this.to = to
      this._splitFn = dependencies.splitFn.bind(dependencies)
    }

    split(separator: number): { path1: ITreePath | undefined; path2: ITreePath | undefined } {
      return this._splitFn(separator)
    }

    kthNodeOnPath(k: number): number {
      throw new Error('Method not implemented.')
    }

    onPath(node: number): boolean {
      throw new Error('Method not implemented.')
    }

    hasIntersection(other: ITreePath): boolean {
      throw new Error('Method not implemented.')
    }

    getIntersection(other: ITreePath): { p1: number; p2: number } | undefined {
      throw new Error('Method not implemented.')
    }

    countIntersection(other: ITreePath): number {
      throw new Error('Method not implemented.')
    }

    get lca(): number {
      throw new Error('Method not implemented.')
    }

    get length(): number {
      throw new Error('Method not implemented.')
    }
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
    bl = new CompressedBinaryLift(tree, 0)
  })

  // it('debug', () => {
  //   const p3_3 = createPath(3, 3)
  //   console.log(p3_3.split(3))
  // })

  // 6种情况:
  // down和to在一条链上，此时separator为:down/to/非down非to
  // down和to不在一条链上，此时separator为:down/to/非down非to
  it('should support split', () => {
    const p3_6 = createPath(3, 6)
    const { path1: pathx_x_3_6_3, path2: path1_6 } = p3_6.split(3)
    expect(pathx_x_3_6_3).toBeUndefined()
    expect(path1_6?.from).toBe(1)
    expect(path1_6?.to).toBe(6)
    const { path1: path3_3, path2: path4_6 } = p3_6.split(1)
    expect(path3_3?.from).toBe(3)
    expect(path3_3?.to).toBe(3)
    expect(path4_6?.from).toBe(4)
    expect(path4_6?.to).toBe(6)
    const { path1: path3_1, path2: path6_6 } = p3_6.split(4)
    expect(path3_1?.from).toBe(3)
    expect(path3_1?.to).toBe(1)
    expect(path6_6?.from).toBe(6)
    expect(path6_6?.to).toBe(6)
    const { path1: path3_4, path2: pathx_x_3_6_6 } = p3_6.split(6)
    expect(path3_4?.from).toBe(3)
    expect(path3_4?.to).toBe(4)
    expect(pathx_x_3_6_6).toBeUndefined()

    const p5_3 = createPath(5, 3)
    const { path1: path5_3_5, path2: path2_3 } = p5_3.split(5)
    expect(path5_3_5).toBeUndefined()
    expect(path2_3?.from).toBe(2)
    expect(path2_3?.to).toBe(3)
    const { path1: path5_5, path2: path0_3 } = p5_3.split(2)
    expect(path5_5?.from).toBe(5)
    expect(path5_5?.to).toBe(5)
    expect(path0_3?.from).toBe(0)
    expect(path0_3?.to).toBe(3)
    const { path1: path5_2, path2: path1_3 } = p5_3.split(0)
    expect(path5_2?.from).toBe(5)
    expect(path5_2?.to).toBe(2)
    expect(path1_3?.from).toBe(1)
    expect(path1_3?.to).toBe(3)
    const { path1: path5_0, path2: path_5_3_3_3 } = p5_3.split(1)
    expect(path5_0?.from).toBe(5)
    expect(path5_0?.to).toBe(0)
    expect(path_5_3_3_3?.from).toBe(3)
    expect(path_5_3_3_3?.to).toBe(3)
    const { path1: path5_1, path2: path_5_3_x_x } = p5_3.split(3)
    expect(path5_1?.from).toBe(5)
    expect(path5_1?.to).toBe(1)
    expect(path_5_3_x_x).toBeUndefined()

    //! path命名不够用了，命名空间限制
    const check = (
      path: ITreePath,
      separator: number,
      from1: number | undefined,
      to1: number | undefined,
      from2: number | undefined,
      to2: number | undefined
    ): boolean => {
      const { path1, path2 } = path.split(separator)
      const autual: (number | undefined)[] = [path1?.from, path1?.to, path2?.from, path2?.to]
      const expected: (number | undefined)[] = [from1, to1, from2, to2]
      return autual.every((v, i) => v === expected[i])
    }

    const p3_3 = createPath(3, 3)
    expect(check(p3_3, 3, undefined, undefined, undefined, undefined)).toBeTruthy()

    const p6_0 = createPath(6, 0)
    expect(check(p6_0, 6, undefined, undefined, 4, 0)).toBeTruthy()
    expect(check(p6_0, 4, 6, 6, 1, 0)).toBeTruthy()
    expect(check(p6_0, 1, 6, 4, 0, 0)).toBeTruthy()
    expect(check(p6_0, 0, 6, 1, undefined, undefined)).toBeTruthy()

    const p1_6 = createPath(1, 6)
    expect(check(p1_6, 1, undefined, undefined, 4, 6)).toBeTruthy()
    expect(check(p1_6, 4, 1, 1, 6, 6)).toBeTruthy()
    expect(check(p1_6, 6, 1, 4, undefined, undefined)).toBeTruthy()
  })
})
