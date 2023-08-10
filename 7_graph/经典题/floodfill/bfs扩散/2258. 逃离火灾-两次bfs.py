# 2258. 逃离火灾
# https://leetcode.cn/problems/escape-the-spreading-fire/description/

# 0 表示草地。
# 1 表示着火的格子。
# 2 表示一座墙，你跟火都不能通过这个格子。
# 2 <= m, n <= 300
# 4 <= m * n <= 2e4

# 一开始你在最左上角的格子 (0, 0) ，你想要到达最右下角的安全屋格子 (m - 1, n - 1) 。
# 每一分钟，你可以移动到 相邻 的草地格子。
# 每次你移动 之后 ，着火的格子会扩散到所有不是墙的 相邻 格子。
# 请你返回你在初始位置可以停留的 最多 分钟数，且停留完这段时间后你还能安全到达安全屋。
# 如果无法实现，请你返回 -1 。如果不管你在初始位置停留多久，你 总是 能到达安全屋，请你返回 1e9 。

from typing import Deque, List, Set, Tuple
from collections import deque


DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def maximumMinutes(self, grid: List[List[int]]) -> int:
        ...


print(
    Solution().maximumMinutes(
        grid=[
            [0, 2, 0, 0, 0, 0, 0],
            [0, 0, 0, 2, 2, 1, 0],
            [0, 2, 0, 0, 1, 2, 0],
            [0, 0, 2, 2, 2, 0, 2],
            [0, 0, 0, 0, 0, 0, 0],
        ]
    )
)
