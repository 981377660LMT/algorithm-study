from itertools import accumulate
from typing import List, Tuple


MOD = int(1e9 + 7)


class Solution:
    def countCollisions(self, s: str) -> int:
        """directions[i] 的值为 'L'、'R' 或 'S'"""
        # 考虑每个车撞不撞
        return len(s.lstrip('L').rstrip('R')) - s.count('S')
