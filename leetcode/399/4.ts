export {}

const INF = 2e9 // !超过int32使用2e15

// 给你一个整数数组 nums 和一个二维数组 queries，其中 queries[i] = [posi, xi]。

// 对于每个查询 i，首先将 nums[posi] 设置为 xi，然后计算查询 i 的答案，该答案为 nums 中 不包含相邻元素 的子序列的 最大 和。

// 返回所有查询的答案之和。

// 由于最终答案可能非常大，返回其对 109 + 7 取余 的结果。

// 子序列 是指从另一个数组中删除一些或不删除元素而不改变剩余元素顺序得到的数组。

// impl Solution {
//   pub fn maximum_sum_subsequence(mut nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
//       unsafe { Solution::maximum_sum_subsequence_unsafe(nums, queries) }
//   }

//   #[target_feature(enable = "avx2")]
//   pub unsafe fn maximum_sum_subsequence_unsafe(nums: Vec<i32>, queries: Vec<Vec<i32>>) -> i32 {
//       let mut res: i64 = 0;
//       let mut nums64 = nums.iter().map(|&x| x as i64).collect::<Vec<i64>>();

//       for query in queries {
//           let pos = query[0] as usize;
//           let x = query[1];
//           nums64[pos] = x as i64;
//           let mut dp0: i64 = 0;
//           let mut dp1: i64 = nums64[0];
//           for v in nums64.iter().skip(1) {
//               let new_dp0 = dp0.max(dp1);
//               let new_dp1 = dp0 + v;
//               dp0 = new_dp0;
//               dp1 = new_dp1;
//           }
//           res += dp0.max(dp1);
//       }

//       (res % 1_000_000_007) as i32
//   }
// }

function maximumSumSubsequence(nums: number[], queries: number[][]): number {
  const n = nums.length
  nums.reverse()
  for (let i = 0; i < queries.length; i++) {
    queries[i][0] = n - 1 - queries[i][0]
  }
  let res = 0
  const history: { dp0: number; dp1: number }[] = []

  {
    let dp0 = 0
    let dp1 = nums[0]
    history.push({ dp0: 0, dp1: nums[0] })
    for (let j = 1; j < nums.length; j++) {
      const ndp0 = Math.max(dp0, dp1)
      const ndp1 = dp0 + nums[j]
      dp0 = ndp0
      dp1 = ndp1
      history.push({ dp0, dp1 })
    }
  }

  for (let i = 0; i < queries.length; i++) {
    const { 0: pos, 1: x } = queries[i]
    nums[pos] = x
    const startAt = Math.max(0, pos - 1)
    let dp0: number
    let dp1: number
    if (startAt === 0) {
      dp0 = 0
      dp1 = nums[0]
      history[0] = { dp0: 0, dp1: nums[0] }
    } else {
      dp0 = history[startAt].dp0
      dp1 = history[startAt].dp1
    }

    for (let j = startAt + 1; j < nums.length; j++) {
      const ndp0 = Math.max(dp0, dp1)
      const ndp1 = dp0 + nums[j]
      dp0 = ndp0
      dp1 = ndp1
      history[j].dp0 = dp0
      history[j].dp1 = dp1
    }
    res += Math.max(dp0, dp1)
  }

  return res % 1_000_000_007
}

// [4,0,-1,-2,3,1,-1]
// [[3,1],[0,-2],[1,-1],[0,-2],[5,4],[6,-3],[6,-2],[2,-1]]
console.log(
  maximumSumSubsequence(
    [4, 0, -1, -2, 3, 1, -1],
    [
      [3, 1],
      [0, -2],
      [1, -1],
      [0, -2],
      [5, 4],
      [6, -3],
      [6, -2],
      [2, -1]
    ]
  )
)
