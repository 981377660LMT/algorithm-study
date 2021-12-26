from typing import List
from collections import defaultdict
from itertools import accumulate


# 1685. 母题-有序数组中差绝对值之和 前i+1项 后 n-i项
class Solution:
    def getDistances(self, arr: List[int]) -> List[int]:
        def getSumAbsoluteDifferences(nums: List[int]) -> List[int]:
            n = len(nums)
            preSum = [0] + list(accumulate(nums))
            return [
                (num * (i + 1) - preSum[i + 1]) + (preSum[n] - preSum[i] - (n - i) * num)
                for i, num in enumerate(nums)
            ]

        n = len(arr)
        indexes = defaultdict(list)
        res = [0] * n
        for i, num in enumerate(arr):
            indexes[num].append(i)

        for lis in indexes.values():
            diffSums = getSumAbsoluteDifferences(lis)
            for index, diff in zip(lis, diffSums):
                res[index] = diff
        return res


print(Solution().getDistances(arr=[2, 1, 3, 1, 2, 3, 3]))
# print(Solution().getDistances(arr=[10, 5, 10, 10]))
