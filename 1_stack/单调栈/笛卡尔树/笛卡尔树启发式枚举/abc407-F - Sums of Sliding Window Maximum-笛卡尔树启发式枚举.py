# https://atcoder.jp/contests/abc407/tasks/abc407_f
# abc407-F - Sums of Sliding Window Maximum-笛卡尔树启发式枚举
#
# !对所有长度为 k(k=1,2,..,n) 的滑动窗口求最大值之和


from typing import List
from itertools import accumulate


class DiffArray:
    """差分维护区间修改，区间查询."""

    __slots__ = ("_diff", "_dirty")

    def __init__(self, n: int) -> None:
        self._diff = [0] * (n + 1)
        self._dirty = False

    def add(self, start: int, end: int, delta: int) -> None:
        """区间 `[start,end)` 加上 `delta`."""
        if start < 0:
            start = 0
        if end >= len(self._diff):
            end = len(self._diff) - 1
        if start >= end:
            return
        self._dirty = True
        self._diff[start] += delta
        self._diff[end] -= delta

    def build(self) -> None:
        if self._dirty:
            self._diff = list(accumulate(self._diff))
            self._dirty = False

    def get(self, pos: int) -> int:
        """查询下标 `pos` 处的值."""
        self.build()
        return self._diff[pos]

    def getAll(self) -> List[int]:
        self.build()
        return self._diff[:-1]


def sumOfSlidingWindowMaximum(nums: List[int]) -> List[int]:
    # 计算每个 i 的有效区间 [L[i], R[i])，使得 A[i] 是该区间的最大值
    n = len(A)
    left = [0] * n
    right = [0] * n
    stack = []
    for i in range(n):
        while stack and (A[i] > A[stack[-1]] or (A[i] == A[stack[-1]] and i < stack[-1])):
            stack.pop()
        left[i] = stack[-1] + 1 if stack else 0
        stack.append(i)
    stack = []
    for i in range(n - 1, -1, -1):
        while stack and (A[i] > A[stack[-1]] or (A[i] == A[stack[-1]] and i < stack[-1])):
            stack.pop()
        right[i] = stack[-1] if stack else n
        stack.append(i)

    D = DiffArray(n + 1)
    for i, (v, l, r) in enumerate(zip(nums, left, right)):
        leftLen = i - l + 1
        rightLen = r - i
        if leftLen < rightLen:
            for a in range(l, i + 1):
                start, end = i + 1 - a, r + 1 - a
                D.add(start, end, v)
        else:
            for b in range(i + 1, r + 1):
                start, end = b - i, b - (l - 1)
                D.add(start, end, v)

    return [D.get(k) for k in range(1, n + 1)]


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    res = sumOfSlidingWindowMaximum(A)
    print(*res, sep="\n")
