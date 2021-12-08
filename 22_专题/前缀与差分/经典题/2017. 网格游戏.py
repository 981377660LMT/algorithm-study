from typing import List
from itertools import accumulate

# 第一个 机器人想要打击竞争对手，
# 使 第二个 机器人收集到的点数 最小化 。
# 与此相对，第二个 机器人想要 最大化 自己收集到的点数。
# 两个机器人都发挥出自己的 最佳水平 的前提下，返回 第二个 机器人收集到的 点数 。
# 每个机器人只会 向右 ((r, c) 到 (r, c + 1)) 或 向下 ((r, c) 到 (r + 1, c)) 。

# 让第一个机器人取最大值是不对的
class Solution:
    def gridGame(self, grid: List[List[int]]) -> int:
        s1, s2 = [0] + list(accumulate(grid[0])), [0] + list(accumulate(grid[1]))
        res = 2 ** 63 - 1
        # 枚举机器人1的转折点
        for mid in range(len(grid[0])):
            up = s1[-1] - s1[mid + 1]
            down = s2[mid]
            res = min(res, max(up, down))
        return res


print(Solution().gridGame(grid=[[2, 5, 4], [1, 5, 1]]))
