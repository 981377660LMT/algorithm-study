/**
 * @param {number[]} arr  一个排序好的数组 arr
 * @param {number} k
 * @param {number} x
 * @return {number[]}  从数组中找到最靠近 x（两数之差最小）的 k 个数 =》必定连续
 * 直接把K个元素当成数组上面一个能够左右滑动的框. arr[mid]到arr[mid+k]
 * # 排除法，从左右两端依次排除数字，最后留下k个数
 * https://leetcode-cn.com/problems/find-k-closest-elements/solution/shuang-zhi-zhen-by-bullimito-46jj/
 */
const findClosestElements = function (arr: number[], k: number, x: number): number[] {
  const n = arr.length
  if (n === k) return arr
  let l = 0
  let r = n - 1
  let remove = n - k

  while (remove) {
    // 尽量左移 因为整数 a 比整数 b 更接近 x 需要满足 |a - x| == |b - x| 且 a < b
    if (x - arr[l] <= arr[r] - x) r--
    else l++
    remove--
  }

  return arr.slice(l, r + 1)
}

// 原本的数组是有序的，所以我们可以像如下步骤利用这一特点。
// 如果目标 x 小于等于有序数组的第一个元素，那么前 k 个元素就是答案。
// 类似的，如果目标 x 大于等于有序数组的最后一个元素，那么最后 k 个元素就是答案。
