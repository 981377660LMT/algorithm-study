# 3327. 判断 DFS 字符串是否是回文串
# https://leetcode.cn/problems/check-if-dfs-strings-are-palindromes/description/

from order import dfsPostOrder

from random import randint
from typing import List, Optional, Sequence


class Solution:
    def findAnswer(self, parent: List[int], s: str) -> List[bool]:
        n = len(s)
        tree = [[] for _ in range(n)]
        for i, p in enumerate(parent):
            if p != -1:
                tree[p].append(i)

        lid, rid = dfsPostOrder(tree)
        ords = [ord(ch) - ord("a") for ch in s]
        data = [0] * n
        for i in range(n):
            data[rid[i] - 1] = ords[i]

        R = RollingHash()
        table1, table2 = R.build(data), R.build(data[::-1])

        def isPalindrome(start: int, end: int) -> bool:
            return R.query(table1, start, end) == R.query(table2, n - end, n - start)

        return [isPalindrome(lid[i], rid[i]) for i in range(n)]


MOD61 = (1 << 61) - 1


class RollingHash:
    """Rolling hash for sequence.

    Example:
    >>> R = RollingHash()
    >>> arr = [1, 2, 3, 2, 1]
    >>> n = len(arr)
    >>> table1, table2 = R.build(arr), R.build(arr[::-1])
    >>> def isPalindrome(start: int, end: int) -> bool:
             return R.query(table1, start, end) == R.query(table2, n - end, n - start)
    >>> isPalindrome(0, 5)
    True
    >>> isPalindrome(0, 4)
    False
    """

    __slots__ = ("_base", "_power")

    def __init__(self, base: Optional[int] = None):
        if base is None:
            base = randint(1, MOD61 - 1)
        self._base = base
        self._power = [1]

    def build(self, seq: Sequence[int]) -> List[int]:
        """Build hash table for sequence."""
        n = len(seq)
        table = [0] * (n + 1)
        for i in range(n):
            table[i + 1] = (table[i] * self._base + seq[i]) % MOD61
        return table

    def eval(self, seq: Sequence[int]) -> int:
        res = 0
        for x in seq:
            res = (res * self._base + x) % MOD61
        return res

    def query(self, table: List[int], start: int, end: int) -> int:
        if start < 0:
            start = 0
        if end > len(table):
            end = len(table)
        if start >= end:
            return 0
        self._expand(end - start)
        return (table[end] - table[start] * self._power[end - start]) % MOD61

    def combine(self, h1: int, h2: int, h2len: int) -> int:
        """Combine two hash values."""
        self._expand(h2len)
        return (h1 * self._power[h2len] + h2) % MOD61

    def addChar(self, h: int, x: int) -> int:
        return (h * self._base + x) % MOD61

    def lcp(
        self, table1: List[int], start1: int, end1: int, table2: List[int], start2: int, end2: int
    ) -> int:
        """Longest common prefix of s1[start1:end1] and s2[start2:end2]."""
        n = min(end1 - start1, end2 - start2)
        low, high = 0, n + 1
        while high - low > 1:
            mid = (low + high) >> 1
            if self.query(table1, start1, start1 + mid) == self.query(table2, start2, start2 + mid):
                low = mid
            else:
                high = mid
        return low

    def _expand(self, size: int) -> None:
        if len(self._power) < size + 1:
            preSize = len(self._power)
            for i in range(preSize - 1, size):
                self._power.append((self._power[i] * self._base) % MOD61)
