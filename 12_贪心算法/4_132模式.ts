// 给你一个整数数组 nums ，数组中共有 n 个整数。
// 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，并同时满足：
// i < j < k 和 nums[i] < nums[k] < nums[j] 。

// 定一移二 O(n^2)
// 在j移动时如果j对应的值小于i对应的值则将j赋给i
const find132Pattern = (nums: number[]): boolean => {
  if (nums.length <= 2) return false
  let i = 0
  for (let j = 1; j < nums.length; j++) {
    for (let k = j + 1; k < nums.length; k++) {
      if (nums[i] < nums[k] && nums[k] < nums[j]) return true
    }

    if (nums[j] < nums[i]) {
      i = j
    }
  }

  return false
}

console.log(find132Pattern([1, 2, 3, 4]))
console.log(find132Pattern([3, 1, 4, 2]))

export {}
