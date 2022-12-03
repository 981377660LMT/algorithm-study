from typing import List
from collections import defaultdict


# 给定一个整数数组 A，返回其中元素之和可被 K 整除的（连续、非空）子数组的数目。
# 1 <= A.length <= 30000


class Solution:
    def subarraysDivByK(self, nums: List[int], k: int) -> int:
        preSum = defaultdict(int, {0: 1})
        res, curSum = 0, 0
        for num in nums:
            curSum = (curSum + num) % k
            res += preSum[curSum]
            preSum[curSum] += 1
        return res


print(Solution().subarraysDivByK([4, 5, 0, -2, -3, 1], 5))
