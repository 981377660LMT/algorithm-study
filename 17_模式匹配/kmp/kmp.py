# 使用方式类似于AC自动机:
# KMP(pattern)：构造函数, pattern为模式串.
# Match(longer,start): 返回模式串在s中出现的所有位置.
# Move(pos, char): 从当前状态pos沿着char移动到下一个状态, 如果不存在则移动到fail指针指向的状态.
# IsMatched(pos): 判断当前状态pos是否为匹配状态.
# Period(i): 求字符串 longer 的前缀 longer[:i+1] 的最短周期(0<=i<n). 如果不存在周期, 返回0.

# https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html

from typing import List, Optional


def getNext(needle: str) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组
    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    """
    next = [0] * len(needle)
    j = 0
    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]
        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1
        next[i] = j
    return next


def indexOfAll(longer, shorter, start=0) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return [0]
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter)
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            res.append(i - len(shorter) + 1)
            hitJ = next[hitJ - 1]
    return res


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

    def findAll(self, longer: str, start=0) -> List[int]:
        """findAll/indexOfAll.
        `o(n+m)`求搜索串 longer 中所有匹配 pattern 的位置.
        """
        res = []
        pos = 0
        for i in range(start, len(longer)):
            pos = self.move(pos, longer[i])
            if self.isMatched(pos):
                res.append(i - len(self._pattern) + 1)
                pos = self.next[pos - 1]  # rollback
        return res

    def find(self, longer: str, start=0) -> int:
        """find/indexOf.
        `o(n+m)`求搜索串 longer 中第一个匹配 pattern 的位置.
        """
        pos = 0
        for i in range(start, len(longer)):
            pos = self.move(pos, longer[i])
            if self.isMatched(pos):
                return i - len(self._pattern) + 1
        return -1

    def move(self, pos: int, input_: str) -> int:
        assert 0 <= pos < len(self._pattern)
        while pos and input_ != self._pattern[pos]:
            pos = self.next[pos - 1]  # rollback
        if input_ == self._pattern[pos]:
            pos += 1
        return pos

    def isMatched(self, pos: int) -> bool:
        return pos == len(self._pattern)

    def period(self, i: Optional[int] = None) -> int:
        """
        求字符串 longer 的前缀 longer[:i+1] 的最短周期(0<=i<n)
        如果不存在周期, 返回0.
        """
        if i is None:
            i = len(self._pattern) - 1
        assert 0 <= i < len(self._pattern)
        res = (i + 1) - self.next[i]
        if res and (i + 1) > res and (i + 1) % res == 0:
            return res
        return 0


if __name__ == "__main__":
    next = getNext("aabaabaabaab")  # 模式串的next数组
    assert next == [0, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

    allPos = indexOfAll("aabaabaabaab", "aab")
    assert allPos == [0, 3, 6, 9]

    kmp = KMP("aab")
    assert kmp.findAll("aabaabaabaab") == [0, 3, 6, 9]
    assert kmp.findAll("aabaabaabaab", 1) == [3, 6, 9]

    pos = 0
    nextPos = kmp.move(pos, "a")
    assert nextPos == 1
    nextPos = kmp.move(nextPos, "a")
    assert nextPos == 2
    nextPos = kmp.move(nextPos, "b")
    assert kmp.isMatched(nextPos)
