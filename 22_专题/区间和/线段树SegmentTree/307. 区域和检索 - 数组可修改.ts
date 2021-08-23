import { SegmentTree } from './区间和线段树'

class NumArray {
  private segmentTree: SegmentTree<number>

  /**
   *
   * @param nums 1 <= nums.length <= 3 * 10**4
   * 最多调用 3 * 10**4 次 update 和 sumRange 方法
   * @description 根据题目可知需要小于O(n)的复杂度:O(logn)线段树
   */
  constructor(nums: number[]) {
    this.segmentTree = new SegmentTree(nums, (a, b) => a + b)
  }

  update(index: number, val: number): void {
    this.segmentTree.update(index, val)
    console.log(this.segmentTree)
  }

  sumRange(left: number, right: number): number {
    return this.segmentTree.query(left, right)
  }
}

const numArray = new NumArray([1, 3, 5])
console.log(numArray.sumRange(0, 2)) // 9
numArray.update(1, 2)
console.log(numArray.sumRange(0, 2)) // 8
export {}
