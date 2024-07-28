from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 正整数 数组 nums。

# 小红和小明正在玩游戏。在游戏中，小红可以从 nums 中选择所有个位数 或 所有两位数，剩余的数字归小明所有。如果小红所选数字之和 严格大于 小明的数字之和，则小红获胜。


# 如果小红能赢得这场游戏，返回 true；否则，返回 false
class Solution:
    def canAliceWin(self, nums: List[int]) -> bool:
        return sum(x if x < 10 else -x for x in nums) != 0
