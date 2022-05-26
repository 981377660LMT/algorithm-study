# O(S^(4/3))解法判断是否是repeating string

# 还是要z函数最快
from typing import Tuple


class Solution:
    def findRepetend(self, s: str) -> Tuple[bool, str]:
        """检测循环节"""
        n = len(s)
        for len_ in range(1, n // 2 + 1):
            # 因子大约为n^(1/3)个
            if n % len_ == 0 and s[:len_] * (n // len_) == s:
                return True, s[:len]
        return False, ''


def another(s: str) -> bool:
    period = (s + s).find(s, 1, -1)
    return period != -1
