from typing import List
from collections import defaultdict
from sortedcontainers import SortedList

# 1 <= nums.length <= 105
# -104 <= nums[i] <= 104


class Solution:
    def countSmaller(self, nums: List[int]) -> List[int]:
        """求逆序对数量"""
        OFFSET = int(1e4) + 10
        res = []
        bit = BIT(3 * OFFSET)
        for i in range(len(nums) - 1, -1, -1):
            cur = bit.query(0, nums[i] - 1 + OFFSET)
            res.append(cur)
            bit.add(nums[i] + OFFSET, nums[i] + OFFSET, 1)
        return list(reversed(res))


# 1 <= nums.length <= 105
# -104 <= nums[i] <= 104
class Solution2:
    def countSmaller(self, nums: List[int]):
        res = []
        visited = SortedList()
        for num in reversed(nums):
            index = visited.bisect_left(num)
            res.append(index)
            visited.add(num)

        return res[::-1]


class BIT:
    def __init__(self, n: int):
        self.size = n
        self._tree1 = defaultdict(int)
        self._tree2 = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')

        rawIndex = index
        while index <= self.size:
            self._tree1[index] += delta
            self._tree2[index] += (rawIndex - 1) * delta
            index += self._lowbit(index)

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1[index] - self._tree2[index]
            index -= self._lowbit(index)
        return res
