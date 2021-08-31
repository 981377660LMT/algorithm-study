/**
 * @param {number[]} nums
 * @return {number}
 * 给定一个未排序的整数数组 nums ，找出数字连续的最长序列(不要求序列元素在原数组中连续）的长度
 * 请你设计并实现时间复杂度为 O(n) 的算法解决此问题。
 */
const longestConsecutive = function (nums: number[]): number {
  if (!nums.length) return 0
  let res = 0
  const set = new Set(nums)
  for (let num of set) {
    // 不是左端点 则跳过
    if (set.has(num - 1)) continue
    let tmp = 1
    while (set.has(num + 1)) {
      tmp++
      num++
    }
    res = Math.max(tmp, res)
  }

  return res
}

console.log(longestConsecutive([100, 4, 200, 1, 3, 2]))

// 解释：最长数字连续序列是 [1, 2, 3, 4]。它的长度为 4。

// 第一思路:计数排序后计算连续1的个数最大值 但是数字太大
// 第二思路:需要对每个左端点查询操作 O(1)的查询使用集合
