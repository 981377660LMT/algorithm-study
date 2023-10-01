import { SparseTable } from '../../22_专题/RMQ问题/SparseTable'

export {}

const INF = 2e15
// class Solution:
//     def maxSubarrays(self, nums: List[int]) -> int:
//         and_ = reduce(lambda x, y: x & y, nums)
//         # 每个子数组and为and_，且尽可能多
function maxSubarrays(nums: number[]): number {
  const and = nums.reduce((x, y) => x & y)
  if (and !== 0) return 1

  const st = new SparseTable(
    nums,
    () => 2 ** 31 - 1,
    (x, y) => x & y
  )

  let res = 0
  let curAnd = 2 ** 31 - 1
  for (let i = 0; i < nums.length; ++i) {
    curAnd &= nums[i]
    if (curAnd === and && (i + 1 === nums.length || st.query(i + 1, nums.length) === and)) {
      res++
      curAnd = 2 ** 31 - 1
    }
  }

  return res
}

// nums = [5,7,1,3]
// console.log(maxSubarrays([5, 7, 1, 3]))
// [1,0,2,0,1,2]
// console.log(maxSubarrays([1, 0, 2, 0, 1, 2]))
// [22,21,29,22]
console.log(maxSubarrays([22, 21, 29, 22]))
