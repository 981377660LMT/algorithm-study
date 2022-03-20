from itertools import accumulate
from typing import List, Tuple


MOD = int(1e9 + 7)


class Solution:
    def countCollisions(self, directions: str) -> int:
        """directions[i] 的值为 'L'、'R' 或 'S'"""

        res = 0
        # 考虑每个车撞不撞
        rNum = [1 if char == 'R' else 0 for i, char in enumerate(directions)]
        LNum = [1 if char == 'L' else 0 for i, char in enumerate(directions)]
        SNum = [1 if char == 'S' else 0 for i, char in enumerate(directions)]
        preSumR = [0] + list(accumulate(rNum))
        preSumL = [0] + list(accumulate(LNum))
        preSumS = [0] + list(accumulate(SNum))

        for i, char in enumerate(directions):
            if char == 'R':
                count = preSumL[-1] - preSumL[i] + preSumS[-1] - preSumS[i]
                if count > 0:
                    res += 1
            elif char == 'L':
                count = preSumR[i] + preSumS[i]
                if count > 0:
                    res += 1

        return res

