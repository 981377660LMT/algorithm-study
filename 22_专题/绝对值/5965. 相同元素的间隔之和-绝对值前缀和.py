from typing import List
from collections import defaultdict
from itertools import accumulate


# 1685. 母题-有序数组中差绝对值之和 前i+1项 后 n-i项
class Solution:
    def getDistances(self, arr: List[int]) -> List[int]:
        """排序+前缀和"""

        def getSumAbsoluteDifferences(nums: List[int]) -> List[int]:
            """到各个点绝对值距离之和"""
            n = len(nums)
            preSum = [0] + list(accumulate(nums))
            return [
                (num * (i + 1) - preSum[i + 1])
                + (preSum[n] - preSum[i] - (n - i) * num)
                for i, num in enumerate(nums)
            ]

        n = len(arr)
        indexMap = defaultdict(list)
        for i, num in enumerate(arr):
            indexMap[num].append(i)

        res = [0] * n
        for indexes in indexMap.values():
            dists = getSumAbsoluteDifferences(indexes)
            for i, v in zip(indexes, dists):
                res[i] = v
        return res


print(Solution().getDistances(arr=[2, 1, 3, 1, 2, 3, 3]))
# print(Solution().getDistances(arr=[10, 5, 10, 10]))
