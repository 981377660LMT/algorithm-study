1751. 2008. 带权区间选择问题

`模板参考 1751`

```
1.按`结束时间`排序
2.每个区间为结尾初始化dp
3.bisectRight找到之前的最右区间，根据 存在/不存在pre 列转移方程

```

```Python
class Solution:
    def maxTaxiEarnings(self, n: int, rides: List[List[int]]) -> int:
        rides.sort(key=lambda x: x[1])
        dp = [e - s + t for s, e, t in rides]
        ends = [e for _, e, _ in rides]

        for i in range(1, len(rides)):
            pre = bisect_right(ends, rides[i][0]) - 1
            if pre >= 0:
                dp[i] = max(dp[i - 1], dp[pre] + rides[i][1] - rides[i][0] + rides[i][2])
            else:
                dp[i] = max(dp[i - 1], rides[i][1] - rides[i][0] + rides[i][2])
        return dp[-1]

```
