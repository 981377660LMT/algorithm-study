// 给定一个整数数组 a，其中1 ≤ a[i] ≤ n （n为数组长度）,
// 其中有些元素出现两次而其他元素出现一次。

// 找到所有出现两次的元素。
// 遍历输入数组，给对应位置的数字取相反数，如果已经是负数，
// 说明前面已经出现过，直接放入输出数组。
function findDuplicates(nums: number[]): number[] {
  const res = []

  for (let i = 0; i < nums.length; i++) {
    const pos = Math.abs(nums[i]) - 1
    if (nums[pos] < 0) res.push(pos + 1)
    nums[pos] *= -1
  }

  return res
}

export {}
