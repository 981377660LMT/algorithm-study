import { randint } from '../randint'

class Solution {
  private nums: number[]
  /**
   *
   * @param nums
   * 给定一个可能含有重复元素的整数数组，
   * 要求随机输出给定的数字的索引。
   * 您可以假设给定的数字一定存在于数组中。
   */
  constructor(nums: number[]) {
    this.nums = nums
  }

  /**
   *
   * @param target
   * 随机输出给定的数字的索引
   * @summary
   * 蓄水池抽样
   * 第i个被选中 = 被选中*不被替换
   */
  pick(target: number): number {
    let count = 0
    let res = 0

    for (let i = 0; i < this.nums.length; i++) {
      if (this.nums[i] === target) {
        count++
        if (randint(1, count) === count) {
          res = i
        }
      }
    }

    return res
  }
}

const solution = new Solution([1, 2, 3, 3, 3])
console.log(solution.pick(1)) // 应该返回索引 2,3 或者 4。每个索引的返回概率应该相等。
console.log(solution.pick(3)) // 应该返回 0。因为只有nums[0]等于1。
export default 1
