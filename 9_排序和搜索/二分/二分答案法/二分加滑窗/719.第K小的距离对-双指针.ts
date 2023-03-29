/**
 * @link https://leetcode-cn.com/problems/find-k-th-smallest-pair-distance/solution/719-zhao-chu-di-k-xiao-de-ju-chi-dui-er-g1i76/
 * @param {number[]} nums 1e4
 * @param {number} k k>=1
 * @return {number}
 * @description 适应二分法
 * @description 给定一个整数数组，返回所有数对之间的第 k 个最小距离
 * 786. 第 K 个最小的素数分数.py
 */
const smallestDistancePair = function (nums: number[], k: number): number {
  nums.sort((a, b) => a - b)

  // !如果不大于距离k的点对数正好是k个，并不代表第K小的距离对的距离就是K，因为有可能比K小
  let left = 0
  let right = 4e15
  while (left <= right) {
    const mid = (left + right) >> 1
    if (countNGT(mid) < k) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }
  return left

  // 不超过mid的答案个数
  function countNGT(mid: number): number {
    let count = 0
    let left = 0
    for (let right = 0; right < nums.length; right++) {
      while (right < nums.length && nums[right] - nums[left] > mid) {
        left++
      }
      count += right - left
    }
    return count
  }
}

console.log(smallestDistancePair([62, 100, 4], 2))

export {}
