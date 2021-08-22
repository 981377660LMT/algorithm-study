// 总的连续子数组个数等于：以索引为 0 结尾的子数组个数 + 以索引为 1 结尾的子数组个数 + ... + 以索引为 n - 1 结尾的子数组个数

// 一维前缀和
const getPrefixArray = (arr: number[]): number[] => {
  const arrCopy = arr.slice()
  arrCopy.reduce((pre, _, index, arr) => (arr[index] += pre))
  return arrCopy
}

// 求一个数组相邻差为 1 连续子数组(索引差 1 的同时，值也差 1)的总个数
const countSubArray = (arr: number[]): number => {
  let res = 0
  let pre = 0
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] - arr[i - 1] === 1) pre++
    else pre = 0
    res += pre
  }
  return res
}

// 求上升子序列个数
const countSubArray2 = (arr: number[]): number => {
  let res = 0
  let pre = 0
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] - arr[i - 1] >= 1) pre++
    else pre = 0
    res += pre
  }
  return res
}

// （连续)子数组的全部元素都不大于 k的子数组个数
// atMostK 就是灵魂方法，一定要掌握，不明白建议多看几遍。
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

// 子数组最大值刚好是 介于[k1,k2] 的子数组的个数
// 其中 k1 < k2
const betweenK = (k1: number, k2: number, nums: number[]) =>
  atMostK(k2, nums) - atMostK(k1 - 1, nums)

if (require.main === module) {
  console.log(getPrefixArray([1, 2, 3, 4]))
  console.log(countSubArray([1, 2, 3]))
  console.log(countSubArray2([1, 2, 3]))
  console.log(atMostK(4, [1, 2, 3]))
  console.log(exactK(3, [1, 2, 3]))
  console.log(betweenK(2, 3, [1, 2, 3]))
}

export { getPrefixArray, atMostK, exactK, betweenK }

// 补充
/////////////////////////////////////////////////////////
// 二维前缀和
const getPrefixArray2 = function (mat: number[][], k: number): void {
  const m = mat.length
  const n = mat[0].length
  const pre = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))

  // 求前缀和
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      pre[i][j] =
        (i > 0 ? pre[i - 1][j] : 0) +
        (j > 0 ? pre[i][j - 1] : 0) -
        (i > 0 && j > 0 ? pre[i - 1][j - 1] : 0) +
        mat[i][j]
    }
  }

  console.table(pre)
}
