// 每个物品的重量为一个范围 [lefti,righti] 可随意调节
// 每个物品价值为vi
// 现在选任意个物品 总重量为w 求最大价值
// 如果不可能完成 返回 -1
// n<=500
// 1<=L<=R<=1e4
// !dp[i][weight] 表示前index个 重量为 weight 时的最大价值
// !不选: dp[i][weight] = dp[i-1][weight]
// !选: dp[i][weight] = max(dp[i-1][weight-R],..., dp[i-1][weight-L]) +vi
// !这里需要RMQ 求区间最大值 由于需要动态更新 所以用线段树或者单调队列维护最值

import { useInput } from '../../../20_杂题/atc競プロ/ts入力'
import { MaxSegmentTree } from '../../../6_tree/线段树/template/线段树区间叠加最大值模板'
import { MaxSegmentTree2 } from '../../../6_tree/线段树/template/线段树区间叠加最大值模板2-atc'

const { input } = useInput()
const [v, n] = input().split(' ').map(Number)
const goods: [left: number, right: number, score: number][] = []
for (let i = 0; i < n; i++) {
  const [left, right, score] = input().split(' ').map(Number)
  goods.push([left, right, score])
}

let dp = Array(v + 1).fill(-Infinity)
dp[0] = 0
for (let i = goods[0][0]; i <= goods[0][1]; i++) {
  // eslint-disable-next-line prefer-destructuring
  dp[i] = goods[0][2]
}

for (let i = 1; i < n; i++) {
  const ndp = dp.slice()
  const [lower, upper, score] = goods[i]
  // const tree = new MaxSegmentTree(ndp)
  const tree = new MaxSegmentTree2(ndp)
  for (let j = lower; j <= v; j++) {
    const [left, right] = [Math.max(0, j - upper), j - lower]
    // const preMax = tree.query(left + 1, right + 1)
    const preMax = tree.query(left, right + 1)
    ndp[j] = Math.max(preMax + score, ndp[j])
  }
  dp = ndp
}

console.log(dp[v] < 0 ? -1 : dp[v])
