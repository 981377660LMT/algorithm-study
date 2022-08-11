# n<=2e5
# 从i跳到j的花费为 (hj-hi)^2 + C
# 求从0跳到n-1的花费最小的跳法


# 斜率优化dp O(n)
# !f[i]=min(f[i],f[j]+(hi-hj)^2)+C
# j左边i右边之后
# 化为直线 dp[j] + h[j]^2 = 2h[i]*h[j] + dp[i] - h[i]^2 - C

# 斜率为正，dp[i]越小，截距越小
# 单调队列维护一个下凸包

from collections import deque


INF = int(1e18)


def calSlope(j1: int, j2: int) -> float:
    """两点连线斜率"""

    def calX(j: int) -> int:
        """横坐标"""
        return heights[j]

    def calY(j: int) -> int:
        """纵坐标"""
        return dp[j] + heights[j] * heights[j]

    if calX(j1) == calX(j2):
        return INF

    return (calY(j2) - calY(j1)) / (calX(j2) - calX(j1))


n, c = map(int, input().split())
heights = list(map(int, input().split()))
dp = [INF] * n
dp[0] = 0

queue = deque([0])
for i in range(n):
    # 1.不是答案的直线出队(下凸包斜率不减)
    while len(queue) >= 2 and calSlope(queue[0], queue[1]) <= 2 * heights[i]:
        queue.popleft()

    candJ = queue[0]
    dp[i] = min(
        dp[i], dp[candJ] + (heights[i] - heights[candJ]) * (heights[i] - heights[candJ]) + c
    )

    # 2.维护下凸包(下凸包斜率不减)
    while len(queue) >= 2 and calSlope(queue[-2], queue[-1]) >= calSlope(queue[-1], i):
        queue.pop()
    queue.append(i)

print(dp[-1])
