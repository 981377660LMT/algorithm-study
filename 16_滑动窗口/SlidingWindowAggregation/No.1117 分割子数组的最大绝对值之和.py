# https://yukicoder.me/problems/no/1117
# 分割子数组的最大绝对值之和
# !将数组分成k段,每段长度不超过m,最大化`每段子数组和的绝对值`的和
# !dp[k][i]表示前i个数分成k段的最大得分

# n<=3000

from itertools import accumulate
from SlidingWindowAggregation import SlidingWindowAggregation

INF = int(1e18)


n, k, m = map(int, input().split())
nums = list(map(int, input().split()))


def max(x, y):
    if x > y:
        return x
    return y


preSum = [0] + list(accumulate(nums))
dp = [-INF] * (n + 1)
dp[0] = 0
for _ in range(k):
    ndp = [-INF] * (n + 1)
    s1 = SlidingWindowAggregation(lambda: -INF, max)
    s2 = SlidingWindowAggregation(lambda: -INF, max)
    for i in range(n + 1):
        ndp[i] = max(ndp[i], s1.query() - preSum[i])
        ndp[i] = max(ndp[i], s2.query() + preSum[i])
        s1.append(dp[i] + preSum[i])
        s2.append(dp[i] - preSum[i])
        if len(s1) > m:
            s1.popleft()
        if len(s2) > m:
            s2.popleft()
    dp = ndp


print(dp[n])
