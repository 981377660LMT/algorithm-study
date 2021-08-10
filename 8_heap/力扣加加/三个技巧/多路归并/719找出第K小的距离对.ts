/**
 * @link https://leetcode-cn.com/problems/find-k-th-smallest-pair-distance/solution/719-zhao-chu-di-k-xiao-de-ju-chi-dui-er-g1i76/
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @description 数据10000大小不能双循环
 * @description 适应二分法
 * @description 求第 k 小的数比较容易想到的就是堆和二分法。二分的原因在于求第 k 小，本质就是求不大于其本身的有 k - 1 个的那个数。而这个问题很多时候满足单调性，因此就可使用二分来解决。
 */
const smallestDistancePair = function (nums: number[], k: number): number {
  nums.sort((a, b) => a - b)
  /**
   *
   * @param n 距离不大于n的点对数
   * @description 复杂度O(n) 是一个滑动窗口
   */
  const countNGT = (n: number): number => {
    let count = 0
    let l = 0
    let r = 0
    while (r < nums.length) {
      while (nums[r] - nums[l] > n) {
        l++
      }
      count += r - l
      r++
    }

    return count
  }

  // 距离最大最小值之差
  let l = 0
  let r = nums[nums.length - 1] - nums[0]

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    // 如果不大于距离k的点对数正好是k个，并不代表第K小的距离对的距离就是K，因为有可能比K小
    countNGT(mid) >= k && (r = mid - 1)
    countNGT(mid) < k && (l = mid + 1)
  }

  return l
}

console.log(smallestDistancePair([62, 100, 4], 2))

export {}
