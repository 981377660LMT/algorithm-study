from functools import lru_cache
from typing import List

# 1458. 两个子序列的最大点积
# 请你返回 nums1 和 nums2 中两个长度相同的 非空 子序列的最大点积。
# 数组的非空子序列是通过删除原数组中某些元素（可能一个也不删除）后剩余数字组成的序列，
# 但不能改变数字间相对顺序。
# 比方说，[2,3,5] 是 [1,2,3,4,5] 的一个子序列而 [1,5,3] 不是。

INF = int(1e18)


class Solution:
    def maxDotProduct(self, nums1: List[int], nums2: List[int]) -> int:
        """
        数组中每个元素只有选和不选两种状态
        """

        @lru_cache(None)
        def dfs(i: int, j: int) -> int:
            if i == n or j == m:
                return -INF

            # 选
            cand1 = nums1[i] * nums2[j]
            cand2 = dfs(i + 1, j + 1) + nums1[i] * nums2[j]
            # 不选
            cand3 = dfs(i + 1, j)
            cand4 = dfs(i, j + 1)
            cand5 = dfs(i + 1, j + 1)
            return max(cand1, cand2, cand3, cand4, cand5)

        n, m = len(nums1), len(nums2)
        res = dfs(0, 0)
        dfs.cache_clear()
        return res


print(Solution().maxDotProduct(nums1=[2, 1, -2, 5], nums2=[3, 0, -6]))
# 输出：18
# 解释：从 nums1 中得到子序列 [2,-2] ，从 nums2 中得到子序列 [3,-6] 。
# 它们的点积为 (2*3 + (-2)*(-6)) = 18 。
