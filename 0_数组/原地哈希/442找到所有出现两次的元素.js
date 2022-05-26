// 给定一个整数数组 a，其中1 ≤ a[i] ≤ n （n为数组长度）, 其中有些元素出现两次而其他元素出现一次。

// 找到所有出现两次的元素。
// 遍历输入数组，给对应位置的数字取相反数，如果已经是负数，说明前面已经出现过，直接放入输出数组。
var findDuplicates = function (nums) {
  const result = []

  for (let i = 0; i < nums.length; i++) {
    const idx = Math.abs(nums[i])
    if (nums[idx - 1] < 0) result.push(idx)
    nums[idx - 1] *= -1
  }
  return result
}
