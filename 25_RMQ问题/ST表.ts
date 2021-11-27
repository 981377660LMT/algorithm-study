// 令f[i][j]  表示区间[i,i+(1<<j)-1]  的最大值。
// dp第二维的大小根据数据范围决定，不小于log(一维长度)
// 给定数组长度
// https://zhuanlan.zhihu.com/p/105439034

class ST {
  private N: number
  private dp: number[][]
  private input: number[]

  constructor(input: number[]) {
    this.N = input.length
    this.input = input
    this.dp = Array.from<unknown, number[]>({ length: this.N }, () => Array(20).fill(0))

    for (let i = 1; i <= this.N; i++) {
      this.dp[i][0] = 0
    }

    for (let j = 1; j <= 20; j++) {
      for (let i = 1; i + (1 << j) - 1 <= this.N; i++) {
        this.dp[i][j] = Math.max(this.dp[i][j - 1], this.dp[i + (1 << (j - 1))][j - 1])
      }
    }
  }

  // 查询时，我们需要找到两个 [l,r] 的子区间，它们的并集恰是 [l,r]
  query(left: number, right: number): number {}
}
