from random import randint
from typing import List, Optional, Sequence


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


if __name__ == "__main__":

    class Solution:
        # 3327. 判断 DFS 字符串是否是回文串
        # https://leetcode.cn/problems/check-if-dfs-strings-are-palindromes/
        def findAnswer(self, parent: List[int], s: str) -> List[bool]:
            n = len(s)
            adjList = [[] for _ in range(n)]
            for i in range(1, n):
                adjList[parent[i]].append(i)

            H = RollingHash()
            hash1, hash2 = [0] * n, [0] * n
            size = [0] * n

            def dfs(cur: int, pre: int) -> None:
                curH1, curH2 = 0, ord(s[cur])
                curSize = 1
                for next in adjList[cur]:
                    if next == pre:
                        continue
                    dfs(next, cur)
                    curH1 = H.combine(curH1, hash1[next], size[next])
                    curSize += size[next]

                for next in reversed(adjList[cur]):
                    if next == pre:
                        continue
                    curH2 = H.combine(curH2, hash2[next], size[next])

                curH1 = H.addChar(curH1, ord(s[cur]))
                hash1[cur] = curH1
                hash2[cur] = curH2
                size[cur] = curSize

            dfs(0, -1)
            return [hash1[i] == hash2[i] for i in range(n)]
