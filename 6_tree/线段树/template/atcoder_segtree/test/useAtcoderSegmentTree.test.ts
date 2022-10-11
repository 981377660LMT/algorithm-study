import { useAtcoderSegmentTree, Operation, AtcoderSegmentTree } from '../AtcoderSegmentTree'

const INF = 2e15

describe('useAtcoderSegmentTree', () => {
  describe('MaxSegmentTree', () => {
    const maxOperation: Operation<number, number> = {
      dataUnit: () => -INF,
      lazyUnit: () => 0,
      mergeChildren: (a, b) => Math.max(a, b),
      updateDataByLazy: (data, lazy) => data + lazy,
      updateLazyByLazy: (lazy1, lazy2) => lazy1 + lazy2
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
})
