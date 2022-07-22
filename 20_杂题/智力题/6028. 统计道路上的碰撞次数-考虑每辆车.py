from itertools import accumulate
from typing import List, Tuple


MOD = int(1e9 + 7)

# 当两辆移动方向 相反 的车相撞时，碰撞次数加 2 。
# 当一辆移动的车和一辆静止的车相撞时，碰撞次数加 1 。
# !道路上的碰撞次数

# !ps:如果只要判断是否发生碰撞 只需
# !"RL" in s or "RS" in s or "SL" in s


class Solution:
    def countCollisions(self, s: str) -> int:
        """directions[i] 的值为 'L'、'R' 或 'S'"""
        # !考虑每个车撞不撞
        return len(s.lstrip("L").rstrip("R")) - s.count("S")

    def countCollisions2(self, directions: str) -> int:
        """directions[i] 的值为 'L'、'R' 或 'S'"""

        res = 0
        # !考虑每个车撞不撞 前缀和判断某个前缀/后缀是否存在相反的车或者墙
        rNum = [1 if char == "R" else 0 for char in directions]
        LNum = [1 if char == "L" else 0 for char in directions]
        SNum = [1 if char == "S" else 0 for char in directions]
        preSumR = [0] + list(accumulate(rNum))
        preSumL = [0] + list(accumulate(LNum))
        preSumS = [0] + list(accumulate(SNum))

        for i, char in enumerate(directions):
            if char == "R":
                count = preSumL[-1] - preSumL[i] + preSumS[-1] - preSumS[i]
                if count > 0:
                    res += 1
            elif char == "L":
                count = preSumR[i] + preSumS[i]
                if count > 0:
                    res += 1

        return res
