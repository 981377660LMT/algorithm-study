import { BinaryIndexedTree } from './BIT'

/**
 * 动态区间和求解
 */
class NumArray {
  private nums: number[]
  private bit: BinaryIndexedTree

  constructor(nums: number[]) {
    this.nums = nums
    this.bit = new BinaryIndexedTree(nums)
  }

  update(index: number, val: number): void {
    this.bit.add(index + 1, val - this.nums[index])
    this.nums[index] = val
  }

  sumRange(left: number, right: number): number {
    return this.bit.query(right + 1) - this.bit.query(left)
  }
}

const numArray = new NumArray([1, 3, 5])

console.log(numArray.sumRange(0, 2))
numArray.update(1, 2)
console.log(numArray.sumRange(0, 2))
