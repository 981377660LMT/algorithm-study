// 如果两个数gcd>1 那么这两个数可交换位置(连通)
// 如果能使用上述交换方式将 nums 按 非递减顺序 排列，返回 true

import { useUnionFindArray } from '../../useUnionFind'

//  1 <= nums.length <= 3 * 104
//  2 <= nums[i] <= 105

// 输入：nums = [7,21,3]
// 输出：true
// 解释：可以执行下述操作完成对 [7,21,3] 的排序：
// - 交换 7 和 21 因为 gcd(7,21) = 7 。nums = [21,7,3]
// - 交换 21 和 3 因为 gcd(21,3) = 3 。nums = [3,7,21]

// 最终可交换位置=> 本身相等或他们具有相同的根
// 如果用两层循环来判断来合并任意两个数，此时必然会超时
// 因此考虑`将每个数和自己的所有(质)因子进行合并`

function gcdSort(nums: number[]): boolean {
  const n = Math.max(...nums)
  const uf = useUnionFindArray(n + 1)
  for (const num of nums) {
    for (let factor = 2; factor <= ~~Math.sqrt(num); factor++) {
      if (num % factor === 0) {
        uf.union(num, factor)
        uf.union(num, num / factor)
      }
    }
  }

  const sortedNums = nums.slice().sort((a, b) => a - b)
  for (let i = 0; i < nums.length; i++) {
    if (nums[i] === sortedNums[i]) continue
    if (uf.find(nums[i]) !== uf.find(sortedNums[i])) return false
  }

  return true
}

console.log(gcdSort([7, 21, 3]))

export {}
