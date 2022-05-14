// 计数法的解
// 1 <= nums.length <= 20000
// 1 <= k <= 99
// 1 <= nums[i] <= 100

// 请你返回数对 (i, j) 的数目，满足 i < j 且 |nums[i] - nums[j]| == k 。
// 小值域直接在值域搜答案
function countKDifference(nums: number[], k: number): number {
  const counter = Array<number>(101).fill(0)
  nums.forEach(num => counter[num]++)

  let res = 0
  for (let value = 0; value + k < counter.length; value++) {
    res += counter[value] * counter[value + k]
  }

  return res
}

console.log(countKDifference([3, 2, 1, 5, 4], 2))
