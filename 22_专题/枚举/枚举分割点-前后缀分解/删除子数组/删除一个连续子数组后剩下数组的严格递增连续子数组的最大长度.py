# 字节笔试 字节2022.5.6日笔试


from collections import defaultdict
from typing import List


class BIT:
    """单点修改 维护区间最值"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def update(self, index: int, target: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] = max(self.tree[index], target)
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res = max(res, self.tree[index])
            index -= self._lowbit(index)
        return res


def maxLenAfterRemove(nums: List[int]) -> int:
    """给出一个数组，最多删除一个连续子数组，求剩下数组的严格递增连续子数组的最大长度。
    n<=1e6
    """
    n = len(nums)
    pos, suf = [1] * n, [1] * n
    for i in range(1, n):
        if nums[i] > nums[i - 1]:
            pos[i] = pos[i - 1] + 1
    for i in range(n - 2, -1, -1):
        if nums[i] < nums[i + 1]:
            suf[i] = suf[i + 1] + 1

    #########################################
    res = 0
    bit = BIT(int(1e6 + 5))
    for i in range(n):
        res = max(res, bit.query(nums[i] - 1) + suf[i])
        bit.update(nums[i], pos[i])

    return res


print(maxLenAfterRemove([5, 3, 4, 9, 2, 8, 6, 7, 1]))  # 4
