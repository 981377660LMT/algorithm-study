/**
 * @link https://leetcode-cn.com/problems/find-k-th-smallest-pair-distance/solution/719-zhao-chu-di-k-xiao-de-ju-chi-dui-er-g1i76/
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @description 数据10000大小不能双循环
 * @description 适应二分法
 * @description 求第 k 小的数比较容易想到的就是堆和二分法。二分的原因在于求第 k 小，本质就是求不大于其本身的有 k - 1 个的那个数。而这个问题很多时候满足单调性，因此就可使用二分来解决。
 * 786. 第 K 个最小的素数分数.py
 */
const smallestDistancePair = function (nums: number[], k: number): number {
  nums.sort((a, b) => a - b)

  // 距离最大最小值之差
  let left = 0
  let right = nums[nums.length - 1] - nums[0]

  while (left <= right) {
    const mid = (left + right) >> 1
    // 如果不大于距离k的点对数正好是k个，并不代表第K小的距离对的距离就是K，因为有可能比K小
    if (countNGT(mid) >= k) {
      right = mid - 1
    } else {
      left = mid + 1
    }
  }

  return left

  /**
   *
   * @param n 距离不大于n的点对数
   * @description 复杂度O(n) 是一个滑动窗口
   */
  function countNGT(n: number): number {
    let count = 0
    let left = 0

    for (let right = 0; right < nums.length; right++) {
      while (nums[right] - nums[left] > n) {
        left++
      }

      count += right - left
    }

    return count
  }
}

console.log(smallestDistancePair([62, 100, 4], 2))

export {}
