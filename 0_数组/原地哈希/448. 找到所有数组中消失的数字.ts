// 给你一个含 n 个整数的数组 nums ，其中 nums[i] 在区间 [1, n] 内。
// 请你找出所有在 [1, n] 范围内但没有出现在 nums 中的数字，并以数组的形式返回结果。

// 将所有正数作为数组下标，置对应数组值为负值。
// 那么，仍为正数的位置即为（未出现过）消失的数字。

// 原始数组：[4,3,2,7,8,2,3,1]
// 重置后为：[-4,-3,-2,-7,8,2,-3,-1]
// 结论：[8,2] 分别对应的index为[5,6]（消失的数字）
/**
 * @param {number[]} nums
 * @return {number[]}
 */
function findDisappearedNumbers(nums: number[]): number[] {
  const res: number[] = []

  for (let i = 0; i < nums.length; i++) {
    const mapped = Math.abs(nums[i]) - 1
    nums[mapped] = Math.abs(nums[mapped]) * -1 // 出现过则置为负数
  }

  for (let i = 0; i < nums.length; i++) {
    if (nums[i] > 0) res.push(i + 1)
  }

  return res
}
