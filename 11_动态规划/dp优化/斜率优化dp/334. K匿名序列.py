# 给出一个长度为 n 的非严格递增整数序列，每次操作可以将其中的一个数减少一，
# 问最少多少次操作后能够使得序列中的任何一个数在序列中都至少有 k−1 个数与之相同。


# 就相当于是把一个序列分成若干组，每组都有至少k−1个数字，
# 花费就是这组的数字和sum，
# 再减去最小值min乘以这个组的cnt，
# 也就是sum−(min∗cnt)。
# 分成的组一定是连续的一段，那么就可以设计转移方程了
# dp(i)=min(dp(j)+[p(i)-p(j)]−nums(j)∗(i−j)) 其中i−j≥k。  斜率为i

# 求得是最小值，所以维护一个下凸壳。
# 由于可能除数为0，所以我们用叉乘来维护
# 注意转移时 i−j≥k 因此更新答案时需要判断
from collections import deque
from itertools import accumulate

INF = int(1e20)


def main():
    def calSlope(j1: int, j2: int) -> float:
        """两点连线斜率"""

        def calX(j: int) -> int:
            """横坐标"""
            return nums[j]

        def calY(j: int) -> int:
            """纵坐标"""
            return dp[j] - preSum[j] + j * nums[j]

        if calX(j1) == calX(j2):
            return INF
        return (calY(j2) - calY(j1)) / (calX(j2) - calX(j1))

    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    preSum = [0] + list(accumulate(nums))
    dp = [INF] * (n + 1)
    dp[0] = 0
    queue = deque([0])
    for i in range(1, n + 1):
        # 1.不是答案的直线出队(斜率单调性)
        while len(queue) >= 2 and calSlope(queue[0], queue[1]) <= i:
            queue.popleft()

        maxJ = queue[0]
        if i - maxJ >= k:
            dp[i] = min(dp[i], dp[maxJ] + (preSum[i] - preSum[maxJ]) - (i - maxJ) * nums[maxJ])

        # 2.维护凸包
        while len(queue) >= 2 and calSlope(queue[-2], queue[-1]) >= calSlope(queue[-1], i):
            queue.pop()
        queue.append(i)

    print(dp[-1])


T = int(input())
for _ in range(T):
    main()
