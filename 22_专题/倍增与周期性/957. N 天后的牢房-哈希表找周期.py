from math import floor, log2
from operator import xor
from typing import List, Tuple

# cells.length == 8
# 1 <= N <= 10^9
# 如果一间牢房的两个相邻的房间都被占用或都是空的，那么该牢房就会被占用。 => 即 left^right^1
# 否则，它就会被空置。

# 答案：模拟每一天监狱的状态。
# 注意loop从第一天开始，最好是先算出第一天
# https://leetcode.com/problems/prison-cells-after-n-days/discuss/591304/Simple-Python-Solution

# n要与状态同时更新

# !哈希表记录周期(鸽巢原理)


class Solution:
    def prisonAfterNDays(self, cells: List[int], n: int) -> List[int]:
        def move(preState: List[int]):
            return [int(i > 0 and i < 7 and preState[i - 1] == preState[i + 1]) for i in range(8)]

        visited = dict()
        while n:
            visited[tuple(cells)] = n  # !保存当前状态
            cells = move(cells)  # !线性转移
            n -= 1

            if tuple(cells) in visited:  # !寻找周期加速
                period = visited[tuple(cells)] - n

                n %= period

        return cells

    def prisonAfterNDays2(self, cells: List[int], k: int) -> List[int]:
        """倍增dp

        cells.length == 8
        1 <= k <= 1e9
        """

        def move(preState: int) -> int:
            s1, s2 = preState >> 1, preState << 1
            nextState = s1 ^ s2 ^ 0b11111111  # 两个相邻的房间都被占用或都是空的，那么该牢房就会被占用
            nextState &= 0b01111110  # 行中的第一个和最后一个房间无法有两个相邻的房间
            return nextState

        n = len(cells)
        maxJ = floor(log2(k)) + 1
        dp = [[0] * (1 << n) for _ in range(maxJ + 1)]  # dp[j][i] 表示第i天后2^j天的状态

        for i in range(1 << n):
            dp[0][i] = move(i)

        for j in range(maxJ):
            for i in range(1 << n):
                dp[j + 1][i] = dp[j][dp[j][i]]

        res = int("".join(map(str, cells)), 2)
        for bit in range(maxJ + 1):
            if (k >> bit) & 1:
                res = dp[bit][res]

        return [int(res >> i & 1) for i in reversed(range(n))]


# print(Solution().prisonAfterNDays(cells=[0, 1, 0, 1, 1, 0, 0, 1], n=7))
print(Solution().prisonAfterNDays2(cells=[0, 1, 0, 1, 1, 0, 0, 1], k=7))
# 输出：[0,0,1,1,0,0,0,0]
# 解释：
# 下表概述了监狱每天的状况：
# Day 0: [0, 1, 0, 1, 1, 0, 0, 1]
# Day 1: [0, 1, 1, 0, 0, 0, 0, 0]
# Day 2: [0, 0, 0, 0, 1, 1, 1, 0]
# Day 3: [0, 1, 1, 0, 0, 1, 0, 0]
# Day 4: [0, 0, 0, 0, 0, 1, 0, 0]
# Day 5: [0, 1, 1, 1, 0, 1, 0, 0]
# Day 6: [0, 0, 1, 0, 1, 1, 0, 0]
# Day 7: [0, 0, 1, 1, 0, 0, 0, 0]
