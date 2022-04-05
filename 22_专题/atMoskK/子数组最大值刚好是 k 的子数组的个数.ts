// 子数组最大值不超过k的子数组的个数
const atMostK = (k: number, nums: number[]): number => {
  let res = 0
  let pre = 0
  for (let i = 0; i < nums.length; i++) {
    if (nums[i] <= k) pre++
    else pre = 0
    res += pre
  }
  return res
}

// 子数组最大值刚好是 k 的子数组的个数
const exactK = (k: number, nums: number[]) => atMostK(k, nums) - atMostK(k - 1, nums)

console.log(exactK(3, [1, 2, 3]))
export {}
