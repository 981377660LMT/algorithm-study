# 字节笔试 字节2022.5.6日笔试
# 删除一个连续子数组后剩下数组的严格递增连续子数组的最大长度

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
        """将[index,size]区间的最大值更新为target"""
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index] = max(self.tree[index], target)
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        """查询[1,index]的最大值"""
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res = max(res, self.tree[index])
            index -= self._lowbit(index)
        return res


def maxLenAfterRemove(nums: List[int]) -> int:
    """给出一个数组，最多删除一个连续子数组，求剩下数组的严格递增连续子数组的最大长度。
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
    for i, num in enumerate(nums):
        res = max(res, bit.query(num - 1) + suf[i])
        bit.update(num, pos[i])  # 更新[nums[i],size]的区间最值

    return res


print(maxLenAfterRemove([5, 3, 4, 9, 2, 8, 6, 7, 1]))  # 4

if __name__ == '__main__':
    testBit = BIT(100)
    testBit.update(1, 2)
    testBit.update(5, 100)
    print(testBit.query(3))
    print(testBit.query(5))
    print(testBit.query(6))
    print(testBit.query(7))
    print(testBit.query(8))
    print(testBit.query(9))

