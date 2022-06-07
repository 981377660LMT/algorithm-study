# N≤1e6,
# 有 n 个工厂，由高到低分布在一座山上，工厂 1 在山顶，工厂 n 在山脚。
# 第 i 个工厂目前有成品 pi 件，在第 i 个工厂位置建立仓库的费用是 ci.
# 对于没有建立仓库的工厂，其产品被运往其他的仓库，
# 产品只能往山下运（只能运往编号更大的工厂的仓库），一件产品运送一个单位距离的费用是 1.
# 假设建立的仓库容量都足够大。工厂 i 与 1 的距离是 xi，问总费用最小值。

# 假设在i建工厂 那么
# dp[i]=min(dp[j]+(disti-distj)*countk+...)+costi
# 对数量求前缀和得
# dp[i]=min(dp[j]+ dist[i]*(p1[i]-p1[j]) -(p2[i]-p2[j]))+costi
# 其中 p1是数量前缀和 p2 是dist*count的前缀和

# 移项得
# Fi = Fj + xi(p1i - p1j) - (p2i - p2j) +ci
# !Fj+p2j=xi*p1j+(Fi+p2i-ci-xi*p1i)
# !斜率大于0 要让截距Fi最小 需要单调队列维护下凸包

from collections import deque


def calSlope(j1: int, j2: int) -> float:
    """两点连线斜率"""

    def calX(j: int) -> int:
        """横坐标"""
        return preSum1[j]

    def calY(j: int) -> int:
        """纵坐标"""
        return dp[j] + preSum2[j]

    # if calX(j1) == calX(j2):
    #     return int(1e20)
    return (calY(j2) - calY(j1)) / (calX(j2) - calX(j1))


n = int(input())
factories = []
for _ in range(n):
    # 工厂 i 距离工厂 1 的距离 Xi（其中 X1=0）；
    # 工厂 i 目前已有成品数量 Pi；
    # 在工厂 i 建立仓库的费用 Ci。
    dist, remain, cost = map(int, input().split())
    factories.append((dist, remain, cost))

preSum1, preSum2 = [0], [0]
for dist, remain, _ in factories:
    preSum1.append(preSum1[-1] + remain)
    preSum2.append(preSum2[-1] + dist * remain)


dp = [int(1e20)] * (n + 1)
dp[0] = 0

# 队列维护凸包与最大值
queue = deque([0])

for i in range(1, n + 1):
    distI, _, costI = factories[i - 1]
    # 1.不是答案的直线出队
    while len(queue) >= 2 and calSlope(queue[0], queue[1]) <= distI:
        queue.popleft()
    # Fj + xi(p1i - p1j) - (p2i - p2j) +ci
    minJ = queue[0]
    dp[i] = min(
        dp[i],
        dp[minJ] + distI * (preSum1[i] - preSum1[minJ]) - (preSum2[i] - preSum2[minJ]) + costI,
    )

    # 2.维护下凸包
    while len(queue) >= 2 and calSlope(queue[-2], queue[-1]) >= calSlope(queue[-1], i):
        queue.pop()
    queue.append(i)


print(dp[-1])
