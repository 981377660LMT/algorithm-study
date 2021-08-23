import { SegmentTree } from './线段树'

class NumArray {
  private segmentTree: SegmentTree<number>

  constructor(nums: number[]) {
    this.segmentTree = new SegmentTree(nums, (a, b) => a + b)
  }

  sumRange(left: number, right: number): number {
    return this.segmentTree.query(left, right)
  }
}

const numArray = new NumArray([1, 3, 5])
console.log(numArray.sumRange(0, 2)) // 9

export {}
