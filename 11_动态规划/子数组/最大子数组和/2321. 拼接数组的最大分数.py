from functools import lru_cache
from typing import List, Literal
from kadane import kanade

# 可以交换一段子数组[left,right]
# 问能取到sum(num1)和sum(num2)的最大子数组和


class Solution:
    def maximumsSplicedArray(self, nums1: List[int], nums2: List[int]) -> int:
        """求出两个数组的diff 那么就是要求diff的最大子数组和"""

        def cal(A: List[int], B: List[int]) -> int:
            diff = [b - a for a, b in zip(A, B)]
            return sum(A) + kanade(diff, getMax=True)

        return max(cal(nums1, nums2), cal(nums2, nums1))

    def maximumsSplicedArray2(self, nums1: List[int], nums2: List[int]) -> int:
        """记忆化dfs+state
        
        state 0/1/2 分别表示 没开始换另一个数组/已经开始换了/已经换完了
        """

        def cal(A: List[int], B: List[int]) -> int:
            @lru_cache(None)
            def dfs(index: int, state: Literal[0, 1, 2]) -> int:
                if index == len(A):
                    return 0

                res = -int(1e20)
                if state == 0:  # 没开始换
                    res = max(res, A[index] + dfs(index + 1, 0))
                    res = max(res, B[index] + dfs(index + 1, 1))
                elif state == 1:  # 开始换的途中
                    res = max(res, B[index] + dfs(index + 1, 1))
                    res = max(res, A[index] + dfs(index + 1, 2))
                elif state == 2:  # 已经换了
                    res = max(res, A[index] + dfs(index + 1, 2))
                return res

            res = dfs(0, 0)
            dfs.cache_clear()
            return res

        return max(cal(nums1, nums2), cal(nums2, nums1))
