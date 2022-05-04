from collections import deque
from itertools import accumulate


n = int(input())
a, b, c = map(int, input().split())
nums = list(map(int, input().split()))  # 表示每名士兵的初始战斗力

# 士兵从 1 到 n 编号，要将他们拆分成若干个特别行动队调入战场
# 你总结出一支特别行动队的初始战斗力 x 将按如下公式修正为 x′
# x′=ax2+bx+c  (a<0)
# 作为部队统帅，你要为这支部队进行编队，使得所有特别行动队修正后的战斗力之和最大。
# 试求出这个最大和。

# https://www.acwing.com/solution/content/107979/
# 不管是啥 dp，先想暴力
# dp[i]=max(dp[i],dp[j]+a*(p[i]-p[j])**2+b*(p[i]-p[j])+c)

# 即 (dp[j]+a*p[j]*p[j]-b*p[j]) = (2*a*p[i])*p[j] + (dp[i]-a*p[i]*p[i]-b*p[i]-c)
# 斜率不变;截距越大,dp[i]越大
# 哪根直线(通过哪个点的直线)截距最大呢
# i变化时,斜率逐渐变小
# 发现不在上凸壳的点肯定取不到最大值
# 单调队列维护一个上凸壳 (凸包相邻两点连线的斜率具有单调性)
# 维护凸壳用Andrew定理
# 587. 安装栅栏 copy


def calSlope(j1: int, j2: int) -> float:
    """两点连线斜率"""

    def calX(j: int) -> int:
        """横坐标"""
        return preSum[j]

    def calY(j: int) -> int:
        """纵坐标"""
        return dp[j] + a * preSum[j] * preSum[j] - b * preSum[j]

    return (calY(j2) - calY(j1)) / (calX(j2) - calX(j1))


preSum = [0] + list(accumulate(nums))
dp = [-int(1e20)] * (n + 1)
dp[0] = 0

# 队列维护凸包与最大值
queue = deque([0])

for i in range(1, n + 1):
    # 1.维护队列最值(斜率单调性)
    while len(queue) >= 2 and calSlope(queue[0], queue[1]) >= 2 * a * preSum[i]:
        queue.popleft()

    maxJ = queue[0]
    dp[i] = max(
        dp[i], dp[maxJ] + a * (preSum[i] - preSum[maxJ]) ** 2 + b * (preSum[i] - preSum[maxJ]) + c
    )

    # 2.维护凸包
    while len(queue) >= 2 and calSlope(queue[-2], queue[-1]) <= calSlope(queue[-1], i):
        queue.pop()
    queue.append(i)


print(dp[-1])
