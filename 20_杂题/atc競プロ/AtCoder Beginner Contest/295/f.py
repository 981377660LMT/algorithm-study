import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# T 個のテストケースについて、数字のみからなる文字列
# S と正整数
# L,R が与えられるので、以下の問題を解いてください。

# 正整数
# x に対して
# f(x)= (
# x を ( 先頭に
# 0 を含まないように ) 書き下した文字列の連続部分列のうち
# S と合致するものの個数 ) と定義します。

# 例えば
# S= 22 であるとき、
# f(122)=1,f(123)=0,f(226)=1,f(222)=2 となります。

# このとき、
# k=L
# ∑
# R
# ​
#  f(k) を求めてください。

from functools import lru_cache
from typing import List


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
                pos = 0
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


def cal(upper: str, evil: str) -> int:
    """字典序小于等于upper且不含evil的字符串个数"""

    @lru_cache(None)
    def dfs(index: int, hasLeadingZero: bool, isLimit: bool, pos: int, count: int) -> int:
        """当前在第pos位,isLimit表示是否贴合上界,pos表示匹配了多少个evil字符(KMP的状态)"""
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
                    # 后退
                    nextPos = kmp.next[nextPos - 1]  # TODO
                res += dfs(index + 1, False, (isLimit and cur == up), nextPos, nextCount)
        return res

    n = len(upper)
    kmp = KMP(evil)
    upperNums = [int(x) for x in upper]
    res = dfs(0, True, True, 0, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    # 数位dp
    # L到R中有多少个数满足
    T = int(input())
    for _ in range(T):
        s, L, R = input().split()
        L = int(L)
        R = int(R)
        print((cal(str(R), s) - cal(str(L - 1), s)))
