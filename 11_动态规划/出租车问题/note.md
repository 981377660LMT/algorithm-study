1751. 2008. 带权区间选择问题 **线性 dp+二分查找优化**

**每个区间选还是不选(to jump or not to jump)**

1. 按`结束时间`排序,dp[i]表示前 i 个区间的最优解
2. 遍历，如果不选当前区间，dp[i] = dp[i-1]
3. 如果选当前区间，bisectRight 找到之前的最右区间 j，dp[i] = dp[j] + w[i]，其中 j 是最后一个不相交的区间，w[i]是当前区间的权重

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
