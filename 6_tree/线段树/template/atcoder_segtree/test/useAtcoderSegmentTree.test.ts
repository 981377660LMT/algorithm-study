import { useAtcoderSegmentTree, Operation, AtcoderSegmentTree } from '../AtcoderSegmentTree'

const INF = 2e15

describe('useAtcoderSegmentTree', () => {
  // !叠加更新 区间最大值查询
  describe('MaxSegmentTree', () => {
    const maxOperation: Operation<number, number> = {
      dataUnit: () => -INF,
      lazyUnit: () => 0,
      mergeChildren: (a, b) => Math.max(a, b),
      updateData: (data, lazy) => data + lazy,
      updateLazy: (lazy1, lazy2) => lazy1 + lazy2
    }

    let tree: AtcoderSegmentTree<number, number>
    beforeEach(() => {
      tree = useAtcoderSegmentTree([1, 2, 3, 4, 5], maxOperation)
    })

    it('should support query', () => {
      expect(tree.query(0, 5)).toBe(5)
      expect(tree.query(1, 3)).toBe(3)
      expect(tree.query(0, 0)).toBe(-INF)
    })

    it('should support update', () => {
      tree.update(2, 3, 10)
      expect(tree.query(0, 5)).toBe(13)
      tree.update(4, 5, 10)
      expect(tree.query(4, 5)).toBe(15)
      tree.update(3, 5, -10)
      expect(tree.query(0, 5)).toBe(13)
    })

    it('should support queryAll', () => {
      expect(tree.queryAll()).toBe(5)
    })

    it('should support maxRight', () => {
      expect(tree.maxRight(0, v => v < 3)).toBe(2)
      expect(tree.maxRight(0, v => v <= 4)).toBe(4)
    })

    it('should support minLeft', () => {
      expect(tree.minLeft(5, v => v > 3)).toBe(5)
    })

    it('should throw error when range is invalid', () => {
      try {
        tree.query(3, 2)
      } catch (error) {
        expect(error instanceof RangeError).toBeTruthy()
      }
    })
  })

  // !叠加更新 区间和查询
  describe('SumSegmentTree', () => {
    const sumOperation: Operation<[sum: number, length: number], number> = {
      dataUnit: () => [0, 0],
      lazyUnit: () => 0,
      mergeChildren(data1, data2) {
        return [data1[0] + data2[0], data1[1] + data2[1]]
      },
      updateData(parentLazy, childData) {
        return [childData[0] + parentLazy * childData[1], childData[1]]
      },
      updateLazy(parentLazy, childLazy) {
        return parentLazy + childLazy
      }
    }

    let tree: AtcoderSegmentTree<[sum: number, length: number], number>
    beforeEach(() => {
      tree = useAtcoderSegmentTree(
        Array.from({ length: 5 }, (_, i) => [i, 1]),
        sumOperation
      )
    })

    it('should support query', () => {
      expect(tree.query(0, 5)[0]).toBe(10)
      expect(tree.query(1, 3)[0]).toBe(3)
      expect(tree.query(0, 0)[0]).toBe(0)
    })

    it('should support update', () => {
      tree.update(2, 3, 10)
      expect(tree.query(0, 5)).toStrictEqual([20, 5])
      tree.update(3, 5, -10)
      expect(tree.query(0, 5)).toStrictEqual([0, 5])
    })

    it('should support queryAll', () => {
      expect(tree.queryAll()).toStrictEqual([10, 5])
    })
  })

  // 如题，已知一个数列，你需要进行下面三种操作：
  // 将某区间每一个数乘上 xx
  // 将某区间每一个数加上 xx
  // 求出某区间每一个数的和
  describe('https://www.luogu.com.cn/problem/P3373', () => {
    it('should pass test case', () => {
      const n = 5
      const MOD = 38
      const operation: Operation<[sum: number, length: number], [mul: number, add: number]> = {}
      const tree = useAtcoderSegmentTree(n, operation)
    })
  })
})
