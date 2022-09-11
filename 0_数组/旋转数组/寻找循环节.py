# O(S^(4/3))解法判断是否是repeating string

# 还是要z函数最快
from typing import Tuple

# 1.
class Solution:
    def findRepetend(self, s: str) -> Tuple[bool, str]:
        """是否存在循环节"""
        n = len(s)
        for len_ in range(1, n // 2 + 1):
            # 因子大约为n^(1/3)个
            if n % len_ == 0 and s[:len_] * (n // len_) == s:
                return True, s[:len]
        return False, ""


# 2.
def another(s: str) -> bool:
    """给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。

    O(n) 可以换kmp
    """
    period = (s + s).find(s, 1, -1)
    return period != -1


# 3. 扩展kmp求字符串的最小周期
