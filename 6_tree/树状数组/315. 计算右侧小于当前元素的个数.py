from bisect import bisect_right
from random import randint
from typing import List, Sequence
from collections import defaultdict, deque
from sortedcontainers import SortedList


def countSmaller(nums: List[int]) -> List[int]:
    """sortedList求每个位置处的逆序对数量"""
    n = len(nums)
    res = [0] * n
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res[i] = smaller
        visited.add(nums[i])
    return res


def shiftAndInversions(nums: List[int]) -> List[int]:
    """求出每个轮转数组的逆序对数量"""
    sl = SortedList()
    inv = 0
    for num in nums[::-1]:
        inv += sl.bisect_left(num)
        sl.add(num)

    res = []
    for num in nums:
        res.append(inv)
        inv -= sl.bisect_left(num)
        sl.remove(num)
        inv += len(sl) - sl.bisect_right(num)
        sl.add(num)
    return res


def countInv(nums: Sequence[int]) -> int:
    """求数组逆序对数量之和"""
    n = len(nums)
    res = 0
    visited = SortedList()
    for i in range(n - 1, -1, -1):
        smaller = visited.bisect_left(nums[i])
        res += smaller
        visited.add(nums[i])
    return res


class WindowInv:
    """滑动窗口逆序对"""

    __slots__ = ("_inv", "_sl", "_queue")

    def __init__(self):
        self._inv = 0
        self._sl = SortedList()
        self._queue = deque()

    def append(self, num: int) -> None:
        self._inv += len(self._sl) - self._sl.bisect_right(num)
        self._sl.add(num)
        self._queue.append(num)

    def appendleft(self, num: int) -> None:
        self._inv += self._sl.bisect_left(num)
        self._sl.add(num)
        self._queue.appendleft(num)

    def pop(self) -> int:
        num = self._queue.pop()
        self._inv -= len(self._sl) - self._sl.bisect_right(num)
        self._sl.remove(num)
        return num

    def popleft(self) -> int:
        num = self._queue.popleft()
        self._inv -= self._sl.bisect_left(num)
        self._sl.remove(num)
        return num

    def query(self) -> int:
        return self._inv

    def __len__(self):
        return len(self._queue)

    def __repr__(self):
        return f"InvManager({self._queue}, {self._inv})"


def countSmallerBIT(nums: List[int]) -> List[int]:
    """求每个位置处的逆序对数量  注意值域很大时需要离散化"""
    n = len(nums)
    arr = sorted(nums)
    res = [0] * n
    bit = BIT1(n + 10)
    for i in range(len(nums) - 1, -1, -1):
        pos1 = bisect_right(arr, nums[i] - 1) + 1
        cur = bit.query(pos1)
        res[i] = cur
        pos2 = bisect_right(arr, nums[i]) + 1
        bit.add(pos2, 1)
    return res


class BIT1:
    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


if __name__ == "__main__":
    nums = [1, 3, 2, 3, 1]
    assert countSmaller(nums) == [0, 2, 1, 1, 0]
    assert countSmallerBIT(nums) == [0, 2, 1, 1, 0]
    assert countInv(nums) == 4

    # use bruteforce to test windowInv
    nums1, nums2 = deque(), WindowInv()
    for _ in range(10000):
        op = randint(0, 3)
        if op == 0:
            num = randint(0, 100)
            nums1.append(num)
            nums2.append(num)
        elif op == 1:
            num = randint(0, 100)
            nums1.appendleft(num)
            nums2.appendleft(num)
        elif op == 2:
            if nums1:
                nums1.pop()
                nums2.pop()
        else:
            if nums1:
                nums1.popleft()
                nums2.popleft()
        assert nums2.query() == countInv(nums1)
    print("test windowInv passed")
