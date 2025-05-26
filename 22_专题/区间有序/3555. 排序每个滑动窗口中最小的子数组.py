from typing import Callable, List, TypeVar, Generic
from collections import deque

from sortedcontainers import SortedList

T = TypeVar("T")


class MonoQueue(Generic[T]):

    __slots__ = ("minQueue", "_minQueueCount", "_less", "_len")

    def __init__(self, less: Callable[[T, T], bool]) -> None:
        self.minQueue = deque()
        self._minQueueCount = deque()
        self._less = less
        self._len = 0

    def append(self, value: T) -> "MonoQueue[T]":
        count = 1
        while self.minQueue and self._less(value, self.minQueue[-1]):
            self.minQueue.pop()
            count += self._minQueueCount.pop()
        self.minQueue.append(value)
        self._minQueueCount.append(count)
        self._len += 1
        return self

    def popleft(self) -> None:
        if not self._len:
            raise IndexError("popleft from empty queue")
        self._minQueueCount[0] -= 1
        if self._minQueueCount[0] == 0:
            self.minQueue.popleft()
            self._minQueueCount.popleft()
        self._len -= 1

    def head(self) -> T:
        if not self._len:
            raise ValueError("monoQueue is empty")
        return self.minQueue[0]

    @property
    def min(self) -> "T":
        return self.head()

    def __len__(self) -> int:
        return self._len

    def __repr__(self) -> str:
        class Item:
            __slots__ = ("value", "count")

            def __init__(self, value: T, count: int):
                self.value = value
                self.count = count

            def __repr__(self) -> str:
                return f"[value: {self.value}, count: {self.count}]"

        res = []
        for i in range(len(self.minQueue)):
            res.append(Item(self.minQueue[i], self._minQueueCount[i]))
        return f"MonoQueue({res})"


def nearestLeftGreater(nums: List[int]) -> List[int]:
    """对每个下标i, 返回 i 左侧最近的严格大于 nums[i] 的下标.若不存在则为 -1."""
    n = len(nums)
    res = [-1] * n
    stack = []
    for i, v in enumerate(nums):
        while stack and nums[stack[-1]] <= v:
            stack.pop()
        res[i] = stack[-1] if stack else -1
        stack.append(i)
    return res


def nearestRightSmaller(nums: List[int]) -> List[int]:
    """对每个下标i, 返回 i 右侧最近的严格小于 nums[i] 的下标.若不存在则为 n."""
    n = len(nums)
    res = [n] * n
    stack = []
    for i, v in enumerate(nums):
        while stack and nums[stack[-1]] > v:
            j = stack.pop()
            res[j] = i
        stack.append(i)
    return res


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def minSubarraySort(self, nums: List[int], k: int) -> List[int]:
        n = len(nums)
        lefts, rights = nearestLeftGreater(nums), nearestRightSmaller(nums)
        pairs = []
        for i, (left, right) in enumerate(zip(lefts, rights)):
            if left != -1 and right != n:
                pairs.append((left, right))
            elif left != -1:
                pairs.append((left, i))
            elif right != n:
                pairs.append((i, right))
            else:
                pairs.append((-1, n))

        # !需要排序的区间的左右边界
        leftMin = SortedList(key=lambda x: x[0])
        rightMax = SortedList(key=lambda x: -x[1])

        def add(index: int) -> None:
            leftMin.add(pairs[index])
            rightMax.add(pairs[index])
            L, R = index - k + 1, index
            while leftMin and leftMin[0][0] < L:
                leftMin.pop(0)
            while rightMax and rightMax[0][1] > R:
                rightMax.pop(0)

        def remove(index: int) -> None:
            leftMin.discard(pairs[index])
            rightMax.discard(pairs[index])

        def query() -> int:
            print(list(leftMin), list(rightMax))
            if not leftMin or not rightMax:
                return 0
            return max2(rightMax[0][1] - leftMin[0][0] + 1, 0)

        res = []
        for right in range(n):
            add(right)
            if right >= k:
                remove(right - k)
            if right >= k - 1:
                res.append(query())
        return res


# nums = [1,3,2,4,5], k = 3
# print(Solution().minSubarraySort([1, 3, 2, 4, 5], 3))  # [2, 2, 0]
# [5,4,3,2,1]
print(Solution().minSubarraySort([5, 4, 3, 2, 1], 4))  # [4,4]
