from typing import List
from itertools import accumulate


class Solution:
    def getSumAbsoluteDifferences(self, nums: List[int]) -> List[int]:
        """
        求每个点到其他所有点的距离之和
        排序+前缀和
        """
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        res = []
        for i, num in enumerate(nums):
            leftDist = num * i - preSum[i]
            rightDist = preSum[n] - preSum[i] - num * (n - i)
            res.append(leftDist + rightDist)
        return res


print(Solution().getSumAbsoluteDifferences(nums=[2, 3, 5]))
print(Solution().getSumAbsoluteDifferences(nums=[2]))
# 输出：[4,3,5]
# 解释：假设数组下标从 0 开始，那么
# result[0] = |2-2| + |2-3| + |2-5| = 0 + 1 + 3 = 4，
# result[1] = |3-2| + |3-3| + |3-5| = 1 + 0 + 2 = 3，
# result[2] = |5-2| + |5-3| + |5-5| = 3 + 2 + 0 = 5。
