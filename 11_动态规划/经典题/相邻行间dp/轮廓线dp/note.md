dfs 状态定义方法

# dfs(row,col,cState) 轮廓线 DP，一个格子一个格子地 dp (dp の遷移を細かくする)

**复杂度一般是 `ROW*COL*2^COL`**
之前 COL(+k?) 个格子的安排方式是 cState 用元组记录比较好(轮廓线 DP)

![1656041130719](image/note/1656041130719.png)
[1659. 最大化网格幸福感 2](1659.%20%E6%9C%80%E5%A4%A7%E5%8C%96%E7%BD%91%E6%A0%BC%E5%B9%B8%E7%A6%8F%E6%84%9F2.py)
[LCP 04. 覆盖](LCP%2004.%20%E8%A6%86%E7%9B%96.py)

ps:有的轮廓线 dp 的题中 `不存在连续两个元素相同`的状态 由此可以剪枝
`将状态从 2^COL 变为 fib(COL) 约等于1.62^COL 大幅度剪枝`

# dp(row, preState,curState) 一般的行间状态转移状压 dp，一行一行地 dp

**复杂度一般是 `ROW*4^COL` 或者 `ROW*3^COL`**
