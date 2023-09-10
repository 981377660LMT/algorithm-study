# 求L到R中的数中出现子串s的次数
# 数位dp+KMP


MOD = 998244353


from functools import lru_cache
from typing import List, Optional


class KMP:
    """单模式串匹配"""

    @staticmethod
    def getNext(pattern: str) -> List[int]:
        next = [0] * len(pattern)
        j = 0
        for i in range(1, len(pattern)):
            while j and pattern[i] != pattern[j]:
                j = next[j - 1]
            if pattern[i] == pattern[j]:
                j += 1
            next[i] = j
        return next

    __slots__ = ("next", "_pattern")

    def __init__(self, pattern: str):
        self._pattern = pattern
        self.next = self.getNext(pattern)

    def match(self, s: str, start=0) -> List[int]:
        res = []
        pos = 0
        for i in range(start, len(s)):
            pos = self.move(pos, s[i])
            if self.isMatched(pos):
                res.append(i - len(self._pattern) + 1)
                pos = self.next[pos - 1]  # rollback
        return res

    def move(self, pos: int, char: str) -> int:
        assert 0 <= pos < len(self._pattern)
        while pos and char != self._pattern[pos]:
            pos = self.next[pos - 1]  # rollback
        if char == self._pattern[pos]:
            pos += 1
        return pos

    def isMatched(self, pos: int) -> bool:
        return pos == len(self._pattern)

    def period(self, i: Optional[int] = None) -> int:
        """
        求字符串 S 的前缀 s[:i+1] 的最短周期(0<=i<n)
        如果不存在周期, 返回0.
        """
        if i is None:
            i = len(self._pattern) - 1
        assert 0 <= i < len(self._pattern)
        res = (i + 1) - self.next[i]
        if res and (i + 1) > res and (i + 1) % res == 0:
            return res
        return 0


def cal(upper: str, need: str) -> int:
    """不超过upper的数中有多少个子串匹配need"""

    @lru_cache(None)
    def dfs(index: int, hasLeadingZero: bool, isLimit: bool, pos: int, count: int) -> int:
        """当前在第pos位,isLimit表示是否贴合上界,pos表示匹配了多少个need字符(KMP的状态)"""
        if index == n:
            return count
        res = 0
        up = upperNums[index] if isLimit else 9
        for cur in range(up + 1):
            if hasLeadingZero and cur == 0:
                res += dfs(index + 1, True, isLimit and cur == up, pos, count)
            else:
                select = str(cur)
                nextPos = kmp.move(pos, select)
                nextCount = count
                if kmp.isMatched(nextPos):
                    nextCount += 1
                    # !后退
                    nextPos = kmp.next[nextPos - 1]
                res += dfs(index + 1, False, (isLimit and cur == up), nextPos, nextCount)
        return res

    n = len(upper)
    kmp = KMP(need)
    upperNums = [int(x) for x in upper]
    res = dfs(0, True, True, 0, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    # 数位dp
    # L到R中的数中出现子串s的次数
    T = int(input())
    for _ in range(T):
        s, left, right = input().split()
        left = int(left)
        right = int(right)
        print((cal(str(right), s) - cal(str(left - 1), s)))
