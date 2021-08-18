/**
 * @param {number[]} nums   在一个长度为 n 的数组 nums 里的所有数字都在 0～n-1 的范围内
 * @return {number}
 * 时间复杂度O(n)，空间复杂度O(1)
 * @summary 让 位置i 的地方放元素i
 * 请找出数组中任意一个重复的数字
 */
const findRepeatNumber = function (nums: number[]): number {
  function* helper(): Generator<number> {
    // 提前把len保存到变量中，这样可以避免每轮循环都计算一遍len
    const n = nums.length
    const swap = (i: number, j: number) => {
      ;[nums[i], nums[j]] = [nums[j], nums[i]]
    }

    for (let i = 0; i < n; i++) {
      if (nums[i] === i) continue
      if (nums[i] === nums[nums[i]]) yield nums[i]
      swap(i, nums[i])
    }
  }

  return helper().next().value
}

console.log(findRepeatNumber([2, 3, 1, 0, 2, 5, 3]))
