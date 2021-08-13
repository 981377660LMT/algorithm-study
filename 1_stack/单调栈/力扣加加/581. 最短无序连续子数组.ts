/**
 * @param {number[]} nums
 * @return {number}
 * 给你一个整数数组 nums ，你需要找出一个 连续子数组 ，如果对这个子数组进行升序排序，那么整个数组都会变为升序排序。
 * 请你找出符合题意的 最短 子数组，并输出它的长度。
 * 即对于每个数 找到下一个比他小的数的位置
 * @summary 最简单的解法:排序数组，然后比较两边不同 O(nlogn)
 * 你可以设计一个时间复杂度为 O(n) 的解决方案吗？
 */
var findUnsortedSubarray = function (nums: number[]): number {
  return 0
}

console.log(findUnsortedSubarray([2, 6, 4, 8, 10, 9, 15]))
// 输出：5
// 解释：你只需要对 [6, 4, 8, 10, 9] 进行升序排序，那么整个表都会变为升序排序。
