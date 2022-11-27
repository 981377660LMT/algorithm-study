// 总的连续子数组个数等于：以索引为 0 结尾的子数组个数 + 以索引为 1 结尾的子数组个数 + ... + 以索引为 n - 1 结尾的子数组个数

// 求一个数组相邻差为 1 连续子数组(索引差 1 的同时，值也差 1)的总个数 dp
const countSubArray = (arr: number[]): number => {
  let res = 0
  let dp = 0
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] - arr[i - 1] === 1) dp++
    else dp = 0
    res += dp
  }
  return res
}

// 求上升子序列个数 dp
const countSubArray2 = (arr: number[]): number => {
  let res = 0
  let dp = 0
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] - arr[i - 1] >= 1) dp++
    else dp = 0
    res += dp
  }
  return res
}

// （连续)子数组的全部元素都不大于 k的子数组个数 dp
const atMostK = (k: number, nums: number[]): number => {
  let res = 0
  let dp = 0
  for (let i = 0; i < nums.length; i++) {
    if (nums[i] <= k) dp++
    else dp = 0
    res += dp
  }
  return res
}

// 子数组最大值刚好是 k 的子数组的个数
const exactK = (k: number, nums: number[]) => atMostK(k, nums) - atMostK(k - 1, nums)

// 子数组最大值刚好是 介于[k1,k2] 的子数组的个数
// 其中 k1 < k2
const betweenK = (k1: number, k2: number, nums: number[]) =>
  // eslint-disable-next-line implicit-arrow-linebreak
  atMostK(k2, nums) - atMostK(k1 - 1, nums)

if (require.main === module) {
  console.log(countSubArray([1, 2, 3]))
  console.log(countSubArray2([1, 2, 3]))
  console.log(atMostK(4, [1, 2, 3]))
  console.log(exactK(3, [1, 2, 3]))
  console.log(betweenK(2, 3, [1, 2, 3]))
}

export { atMostK, exactK, betweenK }
