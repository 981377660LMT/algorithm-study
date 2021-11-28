import { BIT } from './树状数组单点更新模板'

/**
 *
 * @param nums 1 <= nums.length <= 10**5
 * @param lower
 * @param upper
 * @description
 * 数组 A 有多少个连续的子数组，其元素只和在 [lower, upper]的范围内。
 * 即：前缀和之差不超过[lower,upper]
 */
const countRangeSum = (nums: number[], lower: number, upper: number): number => {}

console.log(countRangeSum([-2, 5, -1], -2, 2))
// 输出：3
// 解释：存在三个区间：[0,0]、[2,2] 和 [0,2] ，对应的区间和分别是：-2 、-1 、2 。

export {}
