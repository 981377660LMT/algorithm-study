# 使用方式类似于AC自动机:
# KMP(pattern)：构造函数, pattern为模式串.
# Match(longer,start): 返回模式串在s中出现的所有位置.
# Move(pos, char): 从当前状态pos沿着char移动到下一个状态, 如果不存在则移动到fail指针指向的状态.
# IsMatched(pos): 判断当前状态pos是否为匹配状态.
# Period(i): 求字符串 longer 的前缀 longer[:i+1] 的最短周期(0<=i<n). 如果不存在周期, 返回0.

# https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html

from typing import Generic, List, Optional, Sequence, TypeVar

T = TypeVar("T", int, str)


def getNext(needle: Sequence[T]) -> List[int]:
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


def indexOfAll(
    longer: Sequence[T], shorter: Sequence[T], start=0, nexts: Optional[List[int]] = None
) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return []
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter) if nexts is None else nexts
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


def indexOf(
    longer: Sequence[T], shorter: Sequence[T], start=0, nexts: Optional[List[int]] = None
) -> int:
    """kmp O(n+m)求搜索串 `longer` 中第一个匹配 `shorter` 的位置."""
    if not shorter:
        return 0
    if len(longer) < len(shorter):
        return -1
    next = getNext(shorter) if nexts is None else nexts
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            return i - len(shorter) + 1
    return -1


class KMP(Generic[T]):
    """单模式串匹配"""

    @staticmethod
    def getNext(pattern: Sequence[T]) -> List[int]:
        next = [0] * len(pattern)
        j = 0
        for i in range(1, len(pattern)):
            while j > 0 and pattern[i] != pattern[j]:
                j = next[j - 1]
            if pattern[i] == pattern[j]:
                j += 1
            next[i] = j
        return next

    __slots__ = ("next", "_pattern")

    def __init__(self, pattern: Sequence[T]):
        self._pattern = pattern
        self.next = self.getNext(pattern)

    def findAll(self, longer: Sequence[T], start=0) -> List[int]:
        """findAll/indexOfAll.
        `o(n+m)`求搜索串 longer 中所有匹配 pattern 的位置.
        """
        if len(longer) < len(self._pattern):
            return []
        res = []
        pos = 0
        for i in range(start, len(longer)):
            pos = self.move(pos, longer[i])
            if self.accept(pos):
                res.append(i - len(self._pattern) + 1)
                pos = self.next[pos - 1]  # rollback
        return res

    def find(self, longer: Sequence[T], start=0) -> int:
        """find/indexOf.
        `o(n+m)`求搜索串 longer 中第一个匹配 pattern 的位置.
        """
        pos = 0
        for i in range(start, len(longer)):
            pos = self.move(pos, longer[i])
            if self.accept(pos):
                return i - len(self._pattern) + 1
        return -1

    def move(self, pos: int, input_: T) -> int:
        assert 0 <= pos < len(self._pattern)
        while pos > 0 and input_ != self._pattern[pos]:
            pos = self.next[pos - 1]  # fail
        if input_ == self._pattern[pos]:
            pos += 1
        return pos

    def accept(self, pos: int) -> bool:
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
    # https://www.luogu.com.cn/problem/P4824
    # 在longer中不断删除shorter，求剩下的字符串(消消乐).
    def P4824():
        import sys

        input = sys.stdin.readline
        longer = input().strip()
        shorter = input().strip()
        kmp = KMP(shorter)

        pos = 0
        stack = []
        posRecord = [0] * len(longer)
        for i, c in enumerate(longer):
            pos = kmp.move(pos, c)
            posRecord[i] = pos
            stack.append(i)

            if kmp.accept(pos):
                for _ in range(len(shorter)):
                    stack.pop()
                pos = posRecord[stack[-1]] if stack else 0

        res = []
        for v in stack:
            res.append(longer[v])
        print("".join(res))

    P4824()

    # 2855. 使数组成为递增数组的最少右移次数
    # https://leetcode.cn/problems/minimum-right-shifts-to-sort-the-array/description/
    class Solution:
        def minimumRightShifts(self, nums: List[int]) -> int:
            index = indexOf(nums[::-1] + nums[::-1], sorted(nums)[::-1])
            return index if index != -1 else -1

    # 3036. 匹配模式数组的子数组数目 II
    # https://leetcode.cn/problems/number-of-subarrays-that-match-a-pattern-ii/description/
    class Solution2:
        def countMatchingSubarrays(self, nums: List[int], pattern: List[int]) -> int:
            arr = [1 if a > b else -1 if a < b else 0 for a, b in zip(nums, nums[1:])]
            return len(indexOfAll(arr, pattern))

    print(Solution().minimumRightShifts([3, 4, 5, 1, 2]))
    # [1,3,5]
    print(Solution().minimumRightShifts([1, 3, 5]))

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
    assert kmp.accept(nextPos)

    next = KMP.getNext([1, 2, 3, 1, 2, 3, 1, 2, 3, 4])
    assert next == [0, 0, 0, 1, 2, 3, 4, 5, 6, 0]
    next = KMP.getNext([1, 2, 3, 1, 2, 3, 1, 2, 3, 1])
    assert next == [0, 0, 0, 1, 2, 3, 4, 5, 6, 7]
    kmp = KMP([1, 2, 3, 1, 2, 3, 1, 2, 3, 1])
    assert kmp.findAll([1, 2, 3, 1, 2, 3, 1, 2, 3, 1]) == [0]
