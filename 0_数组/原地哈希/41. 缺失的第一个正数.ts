/**
 * @param {number[]} nums
 * @return {number}
 * 给你一个未排序的整数数组 nums ，请你找出其中没有出现的最小的正整数。
 * 请你实现时间复杂度为 O(n) 并且只使用常数级别额外空间的解决方案。
 * @summary 原地哈希
 * We will position every positive integer in the array at its corresponding index
 * ex) 1 at index 0, 2 at index 1, 3 at index 2
 * Therefore, we can find the first missing positive integer by scanning through the array.
 * 总的来说思路就是把3丢到2号位，把1丢到0号位,...
 * 遍历一次数组把大于等于1的和小于数组大小的值放到原数组对应位置，然后再遍历一次数组查当前下标是否和值对应
 */
const firstMissingPositive = function (nums: number[]): number {
  const n = nums.length
  const swap = (i: number, j: number) => {
    ;[nums[i], nums[j]] = [nums[j], nums[i]]
  }
  for (let i = 0; i < n; i++) {
    if (nums[i] === i + 1) continue

    while (nums[i] >= 1 && nums[i] <= n && nums[i] !== nums[nums[i] - 1]) {
      swap(i, nums[i] - 1)
    }
  }

  for (let i = 0; i < n; i++) {
    if (nums[i] !== i + 1) return i + 1
  }

  return nums.length + 1
}

// console.log(firstMissingPositive([3, 4, -1, 1]))
// console.log(firstMissingPositive([1, 2, 0]))
console.log(firstMissingPositive([-1, 4, 2, 1, 9, 10]))
// 与[1,2,3,...,nums.length]逐个比较
