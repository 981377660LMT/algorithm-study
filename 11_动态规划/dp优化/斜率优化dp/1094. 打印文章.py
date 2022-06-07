# 给出 N 个单词，每个单词有个非负权值 Ci，现要将它们分成连续的若干段，
# 每段的代价为此段单词的权值和的平方，还要加一个常数 M，即 (∑Ci)2+M。
# 现在想求出一种最优方案，使得总费用之和最小。
# 0≤N≤5×105,


# dp[i]=min(dp[j]+(p[i]-p[j])**2+M)
# 即 dp[j]+p[j]**2=2*p[i]*p[j]+dp[i]-p[i]*p[i]-M

# 斜率为正，dp[i]越小，截距越小
# 单调队列维护一个下凸包

from collections import deque
from itertools import accumulate


def main():
    def calSlope(j1: int, j2: int) -> float:
        """两点连线斜率"""

        def calX(j: int) -> int:
            """横坐标"""
            return preSum[j]

        def calY(j: int) -> int:
            """纵坐标"""
            return dp[j] + preSum[j] * preSum[j]

        return (calY(j2) - calY(j1)) / (calX(j2) - calX(j1))

    n, M = map(int, input().split())
    nums = list(map(int, input().split()))
    preSum = [0] + list(accumulate(nums))
    dp = [int(1e20)] * (n + 1)
    dp[0] = 0
    queue = deque([0])
    for i in range(1, n + 1):
        # 1.不是答案的直线出队(斜率单调性)
        while len(queue) >= 2 and calSlope(queue[0], queue[1]) <= 2 * preSum[i]:
            queue.popleft()

        maxJ = queue[0]
        dp[i] = min(dp[i], dp[maxJ] + (preSum[i] - preSum[maxJ]) ** 2 + M)

        # 2.维护凸包
        while len(queue) >= 2 and calSlope(queue[-2], queue[-1]) >= calSlope(queue[-1], i):
            queue.pop()
        queue.append(i)

    print(dp[-1])


while True:
    try:
        main()
    except EOFError:
        break
