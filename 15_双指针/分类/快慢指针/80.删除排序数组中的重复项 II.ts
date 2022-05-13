// 删除排序数组中的重复项，使得相同数字最多出现 k 次
/**
 * @param {number[]} nums
 * @return {number}
 * 你必须在 原地 修改输入数组 并在使用 O(1) 额外空间的条件下完成。
 */
const removeDuplicates = function (nums: number[]): number {
  const K = 2
  let slow = 0
  let fast = 0

  // fast将值丢给前面的slow指针 slow指针检查是否有K个 没有则允许赋值前进
  while (fast < nums.length) {
    if (nums[slow - K] !== nums[fast]) {
      // 注意这里是先赋值再前进
      nums[slow] = nums[fast]
      slow++
    }
    fast++
  }

  // 断开连接
  console.log(nums.slice(0, slow))

  return slow
}

console.log(removeDuplicates([1, 1, 1, 2, 2, 3]))
