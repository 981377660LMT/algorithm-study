// https://leetcode.cn/problems/gcd-sort-of-an-array/
// 如果两个数gcd>1 那么这两个数可交换位置(连通)
// 如果能使用上述交换方式将 nums 按 非递减顺序 排列，返回 true
//  1 <= nums.length <= 3 * 104
//  2 <= nums[i] <= 105

// 输入：nums = [7,21,3]
// 输出：true
// 解释：可以执行下述操作完成对 [7,21,3] 的排序：
// - 交换 7 和 21 因为 gcd(7,21) = 7 。nums = [21,7,3]
// - 交换 21 和 3 因为 gcd(21,3) = 3 。nums = [3,7,21]

// 最终可交换位置=> 本身相等或他们具有相同的根
// 如果用两层循环来判断来合并任意两个数，此时必然会超时
// 因此考虑`将每个数和自己的所有因子进行合并`

import { getFactors } from '../../../19_数学/因数筛/prime'
import { UnionFindArray } from '../../UnionFind'

function gcdSort(nums: number[]): boolean {
  const n = Math.max(...nums, 0)
  const uf = new UnionFindArray(n + 1)
  for (let i = 0; i < nums.length; i++) {
    const cur = nums[i]
    const factors = getFactors(cur)
    for (let j = 0; j < factors.length; j++) {
      const f = factors[j]
      // !大于1的因子
      if (f > 1) {
        uf.union(cur, f)
      }
    }
  }

  const sorted = nums.slice().sort((a, b) => a - b)
  for (let i = 0; i < nums.length; i++) {
    if (!uf.isConnected(sorted[i], nums[i])) return false
  }

  return true
}

console.log(gcdSort([7, 21, 3]))

export {}
