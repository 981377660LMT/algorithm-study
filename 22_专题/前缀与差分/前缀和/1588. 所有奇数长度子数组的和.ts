/**
 * @param {number[]} arr
 * @return {number}
 * 前缀和 O(n^2) 可以计算贡献优化到 O(n)
 */
const sumOddLengthSubarrays = function (arr: number[]): number {
  const preSum = [0]
  for (let i = 1; i <= arr.length; i++) {
    preSum.push(preSum[i - 1] + arr[i - 1])
  }

  let res = 0
  for (let len = 1; len <= arr.length; len += 2) {
    for (let start = 0; start + len <= arr.length; start++) {
      const end = start + len
      res += preSum[end] - preSum[start]
    }
  }

  return res
}
console.log(sumOddLengthSubarrays([1, 4, 2, 5, 3]))

// 统计任意值 arr[i] 在奇数子数组的出现次数
// 对于原数组的任意位置 i 而言，其左边共有 i 个数，右边共有 n - i - 1 个数
// const sumOddLengthSubarrays2 = function (arr: number[]): number {
//   const pre = [0]
//   for (let i = 1; i <= arr.length; i++) {
//     pre[i] = pre[i - 1] + arr[i - 1]
//   }

//   let res = 0
//   for (let len = 1; len <= arr.length; len += 2) {
//     for (let l = 0; l + len <= arr.length; l++) {
//       const r = l + len
//       res += pre[r] - pre[l]
//     }
//   }

//   return res
// }
