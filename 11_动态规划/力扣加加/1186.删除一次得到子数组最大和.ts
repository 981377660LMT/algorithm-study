/**
 * @param {number[]} arr
 * @return {number}
 * @description 求连续子数组最大和 可以删除一次元素;删除一个元素后，子数组 不能为空。
 */

// dp[i][0]：不删除元素，以i结尾的连续子数组最大和
// dp[i][1]：删除其中一个得到的最大值，有两种情况（1，删除i 2. 不删除i）
// dp[i+1][0]的计算有两种情况：
// 前面所有数之和大于0，dp[i+1][0]=dp[i][0]+arr[i]
// 前面所有数之和小于0，dp[i+1][0]=arr[i]
// dp[i+1][1]的计算有两种情况：
// 删除位置i的数：dp[i+1][1]=dp[i][0]
// 删除其他位置的数：dp[i+1][1]=dp[i][1]+arr[i]
// const maximumSum = (arr: number[]): number => {
//   const n = arr.length
//   if (n === 1) return arr[0]
//   let res = -Infinity
//   const dp = Array.from<number, [number, number]>({ length: arr.length + 1 }, () => [0, 0])
//   for (let i = 0; i < n; i++) {
//     dp[i + 1][0] = Math.max(dp[i][0] + arr[i], arr[i])
//     dp[i + 1][1] = Math.max(dp[i][1] + arr[i], dp[i][1])
//     res = Math.max(res, dp[i + 1][0], dp[i + 1][1])
//   }
//   return res
// }
const maximumSum = (arr: number[]): number => {
  const n = arr.length
  if (n === 1) return arr[0]

  return 1
}
console.log(maximumSum([1, -2, 0, 3]))
// 输出：4
// 解释：我们可以选出 [1, -2, 0, 3]，然后删掉 -2，这样得到 [1, 0, 3]，和最大。
export default 1
