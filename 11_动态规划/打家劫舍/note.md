- 这类问题可以用两个变量滚动 dp
  `dp0,dp1=0,score[0] => 当前不偷/偷 的最大得分, 结果为max(dp0,dp1)`
  `dp0,dp1=1,pow(2,count,MOD)-1 => 当前不偷/偷 的方案数，结果为dp0+dp1`
