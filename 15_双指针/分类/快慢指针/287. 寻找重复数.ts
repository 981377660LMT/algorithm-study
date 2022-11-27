/**
 * @param {number[]} nums 给定一个包含 n + 1 个整数的数组 nums ，其数字都在 1 到 n 之间（包括 1 和 n），可知至少存在一个重复的整数。
 * @return {number}
 * 假设 nums 只有 一个重复的整数 ，找出 这个重复的数 。
 * 你设计的解决方案必须不修改数组 nums 且只用常量级 O(1) 的额外空间。
 * 这道题与那道链表题一样
 * 链表有next 数组的next可以看作是nums[i]
 * [1, 3, 4, 2, 2]就相当于1->3->2->4->2这个链表 在2处循环
 */
function findDuplicate (nums: number[]): number {
  // 先从0开始各走一步
  let fast = nums[nums[0]]
  let slow = nums[0]

  while (slow !== fast) {
    fast = nums[nums[fast]]
    slow = nums[slow]
  }

  slow = 0

  while (slow !== fast) {
    fast = nums[fast]
    slow = nums[slow]
  }

  return slow
}

console.log(findDuplicate([1, 3, 4, 2, 2]))
