# 295. 数据流的中位数

from typing import List, Optional, Tuple
from heapq import heappush, heappushpop


class MedianFinderAddOnly:
    __slots__ = "_small", "_large", "_smallSum", "_largeSum"

    def __init__(self):
        self._small = []  # 左边，大顶堆
        self._large = []  # 右边，小顶堆
        self._smallSum = 0
        self._largeSum = 0

    def add(self, num: int) -> None:
        if len(self._small) == len(self._large):
            leftMax = -heappushpop(self._small, -num)
            heappush(self._large, leftMax)
            self._smallSum += num - leftMax
            self._largeSum += leftMax
        elif len(self._small) < len(self._large):
            rightMin = heappushpop(self._large, num)
            heappush(self._small, -rightMin)
            self._smallSum += rightMin
            self._largeSum += num - rightMin

    def median(self) -> Tuple[int, Optional[int]]:
        if len(self._small) == len(self._large):
            return -self._small[0], self._large[0]
        elif len(self._small) < len(self._large):
            return self._large[0], None
        else:
            raise Exception("invalid")

    def distToMedian(self) -> int:
        median = self.median()[0]
        sum1 = len(self._small) * median - self._smallSum
        sum2 = self._largeSum - len(self._large) * median
        return sum1 + sum2


if __name__ == "__main__":
    # LCP 24. 数字游戏
    # https://leetcode.cn/problems/5TxKeK/description/
    # 主办方请小扣回答出一个长度为 N 的数组，
    # 第 i 个元素(0 <= i < N)表示将 0~i 号计数器 初始 所示数字操作成满足所有条件
    # nums[a]+1 == nums[a+1],(0 <= a < i) 的最小操作数。
    #
    # 首先将数据进行等效转换以方便计算：g[i] = nums[i] - i。
    # 转换完之后题目变为：给你一个数列g，求把g的所有元素变成同一个值的最小代价，那么显然这个值应该是g的中位数。
    # 使用双堆维护中位数以及中位数两边的数据累加和，便于快速累加差值。

    MOD = int(1e9 + 7)

    class Solution:
        def numsGame(self, nums: List[int]) -> List[int]:
            res = []
            medianFinder = MedianFinderAddOnly()
            for i, num in enumerate(nums):
                v = num - i
                medianFinder.add(v)
                res.append(medianFinder.distToMedian() % MOD)
            return res
