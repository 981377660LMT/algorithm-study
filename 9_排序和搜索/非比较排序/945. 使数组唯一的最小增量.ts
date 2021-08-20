/**
 * @param {number[]} nums
 * @return {number}
 * 每次 move 操作将会选择任意 A[i]，并将其递增 1。
   返回使 A 中的每个值都是唯一的最少操作次数。
   贪心 O(nlogn)
 */
// const minIncrementForUnique = function (nums: number[]): number {
//   nums.sort((a, b) => a - b)
//   let res = 0
//   for (let i = 1; i < nums.length; i++) {
//     if (nums[i - 1] >= nums[i]) {
//       const shouldPlus = nums[i - 1] - nums[i] + 1
//       nums[i] += shouldPlus
//       res += shouldPlus
//     }
//   }
//   return res
// }
// 先计数再遍历
// 两个 1 重复了，需要有一个增加到 2，这样 2 的数量变成了三个。在三个 2 中，又有两个需要增加到 3，然后又出现了两个 3…… 以此类推，
const minIncrementForUnique = function (nums: number[]): number {
  let res = 0
  const counter: number[] = []
  for (let i = 0; i < nums.length; i++) {
    const element = nums[i]
    if (counter[element]) counter[element] += 1
    else counter[element] = 1
  }

  for (let i = 0; i < counter.length; i++) {
    const freq = counter[i]
    if (freq >= 2) {
      counter[i] = 1
      if (counter[i + 1]) counter[i + 1] += freq - 1
      else counter[i + 1] = freq - 1
      res += freq - 1
    }
  }

  return res
}

console.log(minIncrementForUnique([3, 2, 1, 2, 1, 7]))

export default 1
