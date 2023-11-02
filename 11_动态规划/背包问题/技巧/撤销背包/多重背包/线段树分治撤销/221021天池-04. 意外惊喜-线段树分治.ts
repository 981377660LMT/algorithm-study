// 阿里云天池专场
// 有n个单调不减数组，求从这n个数组中选择k个数的最大和，（只能拿数组的一个前缀）
// 221021天池-04. 意外惊喜 https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/
// 比 2218. 从栈中取出K个硬币的最大面值和 https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
// !多了每个礼物包中的礼物价值 非严格递增 的条件，去除了背包总数的限制
// n<=3000 k<=3000 O(nklogn)

// 解：
// !因为数组是单调的，所以能得到一个结论就是我们最多只会拿不全一个数组
// 因为当有两个数组都拿一部分时，那么肯定可以通过减少某个数组拿的数字去拿另一个数组来增大这个和
// 换句话说就是，只要选了一个序列的第一个数，那么最优方案一定是尽量把这个序列的所有数都选上。
// 即：给定 n 个物品，对 i=1,2,⋯,n，你需要求出：去掉第 i 个物品后，对其他物品做背包的结果(可撤销01背包)。
// !所以我们就可以通过枚举哪个数组不全选，采用分治删点(线段树分治/可撤销背包)的方式dp

import { mutateWithoutOne } from '../../../../../../10_分治法/线段树分治/mutateWithOutOne'

function brilliantSurprise(present: number[][], limit: number): number {
  const n = present.length
  const groupSum = present.map(arr => arr.reduce((a, b) => a + b, 0))

  let res = 0
  mutateWithoutOne(Array(limit + 1).fill(0), 0, n, {
    copy(dp) {
      return dp.slice()
    },
    mutate(dp, index) {
      const m = present[index].length
      // 01背包，整个组全选还是全不选
      for (let i = limit; i >= m; i--) {
        dp[i] = Math.max(dp[i], dp[i - m] + groupSum[index])
      }
    },
    query(dp, index) {
      // 可以不全选，即选择一个前缀
      res = Math.max(res, dp[limit])
      const group = present[index]
      let curSum = 0
      for (let i = 0; i < Math.min(group.length, limit); i++) {
        curSum += group[i]
        res = Math.max(res, dp[limit - (i + 1)] + curSum)
      }
    }
  })

  return res
}

export {}

if (require.main === module) {
  // present = [[1,2],[2,3],[3,4]], limit = 3
  console.log(
    brilliantSurprise(
      [
        [1, 2],
        [2, 3],
        [3, 4]
      ],
      3
    )
  )
}
