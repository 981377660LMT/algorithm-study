from typing import List


class Finder:
    __slots__ = ("_n", "_exist", "_prev", "_next")

    def __init__(self, n: int):
        self._n = n
        self._exist = [True] * n
        self._prev = [i - 1 for i in range(n)]
        self._next = [i + 1 for i in range(n)]

    def has(self, i: int) -> bool:
        return 0 <= i < self._n and self._exist[i]

    def erase(self, i: int) -> bool:
        if not self.has(i):
            return False
        l, r = self._prev[i], self._next[i]
        if l >= 0:
            self._next[l] = r
        if r < self._n:
            self._prev[r] = l
        self._exist[i] = False
        return True

    def prev(self, i: int) -> int:
        """
        返回`严格`小于i的最大元素,如果不存在,返回-1.
        !调用时需保证i存在.
        """
        return self._prev[i]

    def next(self, i: int) -> int:
        """
        返回`严格`大于i的最小元素.如果不存在,返回n.
        !调用时需保证i存在.
        """
        return self._next[i]


if __name__ == "__main__":

    class Solution:
        # 3639. 变为活跃状态的最小时间
        # https://leetcode.cn/problems/minimum-time-to-activate-string/description/
        def minTime(self, s: str, order: List[int], k: int) -> int:
            n = len(s)
            count = n * (n + 1) // 2
            if count < k:
                return -1

            finder = Finder(n)
            for t in range(n - 1, -1, -1):
                i = order[t]
                l, r = finder.prev(i), finder.next(i)
                count -= (i - l) * (r - i)
                if count < k:
                    return t
                finder.erase(i)

            raise RuntimeError("Should not reach here.")
