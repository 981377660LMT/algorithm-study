# O(S^(4/3))解法判断是否是repeating string

# 还是要z函数最快
from typing import List, Tuple

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
    return s in (s + s)[1:-1]


# 3. 扩展kmp求字符串的最小周期


def getZ(string: str) -> List[int]:
    """z算法求字符串公共前后缀的长度

    z[0]=0
    z[i]是s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
    """

    n = len(string)
    Z = [0] * n
    left, right = 0, 0
    for i in range(1, n):
        Z[i] = max(min(Z[i - left], right - i + 1), 0)
        while i + Z[i] < n and string[Z[i]] == string[i + Z[i]]:
            left, right = i, i + Z[i]
            Z[i] += 1
    return Z


def getMinPeriod(s: str) -> int:
    """求字符串的最小周期
    当区间[l+d,r]的哈希值与[l,r-d]的哈希值相等时，那么该区间[l,r]是以 d 为循环节的**
    """
    z = getZ(s)
    for i in range(1, len(s)):
        if len(s) % i == 0 and z[i] == len(s) - i:
            return i
    return -1
