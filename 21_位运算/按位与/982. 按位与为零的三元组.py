from collections import defaultdict
from typing import List

# 1 <= A.length <= 1000
# nums[i]<=2^16
# 找出索引为 (i, j, k) 的三元组
# A[i] & A[j] & A[k] == 0，其中 & 表示按位与（AND）操作符。
class Solution:
    def countTriplets(self, A: List[int]) -> int:
        """时间复杂度 n*max(nums[i])"""
        memo = defaultdict(int)
        for n1 in A:
            for n2 in A:
                memo[n1 & n2] += 1  # 结果必然小于2^16

        res = 0
        for num in A:
            for key, val in memo.items():
                if num & key == 0:
                    res += val
        return res

