# 一个区间的木块要粉刷,`每一个木块最多只能粉刷一次`,
# 可以不粉刷,对于一个粉刷匠i而言,他必须粉刷第Si木板(不同工匠的 Si 不同),
# !而且只能粉刷一段连续的木块,并且要求粉刷长度`不超过`Li,
# 而且每刷一个木板,就可以得到Pi的报酬,现在要求报酬最大化,
# 1≤木板数≤1.6e4,
# 1≤人数≤1e2,


# !暗示状态为木板索引，当前人数索引 1.6e6


# 当前工匠，当前木板结尾
# !dp[i][j]=max(dp[i][j-1],dp[i-1][j]) 这个工匠不刷/这个格子不刷
# !dp[i][j]=max(dp[i−1][k]+pi*(j−k)) 其中j−Li≤k≤Si−1
# 将方程中与k有关的部分放到单调队列中(滑窗最值=>单调队列优化)
# 维护一个f[i-1][k]-p[i]*k单减,k单增的队列，即可完成转移
# https://www.acwing.com/solution/content/33067/

from typing import Any
from collections import deque


class MaxQueue(deque):
    @property
    def max(self) -> int:
        if not self:
            raise ValueError('maxQueue is empty')
        return self[0][0]

    def append(self, value: int, *metaInfo: Any) -> None:
        count = 1
        while self and self[-1][0] < value:
            count += self.pop()[-1]
        super().append([value, *metaInfo, count])

    def popleft(self) -> None:
        if not self:
            raise IndexError('popleft from empty queue')
        self[0][-1] -= 1
        if self[0][-1] == 0:
            super().popleft()


N, M = map(int, input().split())
people = []
for _ in range(M):
    length, money, include = map(int, input().split())
    people.append((include, length, money))
people.sort(key=lambda x: x[0])

dp = [0] * (N + 1)

# 换一下下标
for i in range(M):
    include, length, money = people[i]
    ndp, queue = [0] * (N + 1), MaxQueue()  # !注意queue作用于当前层的dp 用于优化jk的两层循环

    for j in range(N + 1):  # !注意是从0到n
        # 不刷格子||不看工匠
        ndp[j] = dp[j]
        if j - 1 >= 0:
            ndp[j] = max(ndp[j], ndp[j - 1])

        # 第i个工匠粉刷第k+1块到第j块木板 该工匠粉刷总数不能超过Li，且必须粉刷Si
        while queue and (queue[0][1] < j - length):
            queue.popleft()
        if queue and j >= include:
            ndp[j] = max(ndp[j], queue.max + money * j)
        if j < include:  # !每次选择的k需要在include前
            queue.append(dp[j] - money * j, j)

    dp = ndp


print(dp[-1])
