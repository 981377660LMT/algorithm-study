from typing import List

# 请你返回 nums1 和 nums2 中两个长度相同的 非空 子序列的最大点积。
# 数组的非空子序列是通过删除原数组中某些元素（可能一个也不删除）后剩余数字组成的序列，但不能改变数字间相对顺序。
# 比方说，[2,3,5] 是 [1,2,3,4,5] 的一个子序列而 [1,5,3] 不是。

# 1035. 不相交的线.ts
class Solution:
    def maxDotProduct(self, nums1: List[int], nums2: List[int]) -> int:
        n, m = len(nums1), len(nums2)

        if max(nums1) < 0 and min(nums2) > 0:
            return max(nums1) * min(nums2)
        if max(nums2) < 0 and min(nums1) > 0:
            return max(nums2) * min(nums1)

        dp = [[0] * (m + 1) for _ in range(n + 1)]

        # 选还是不选
        for i in range(1, n + 1):
            for j in range(1, m + 1):
                dp[i][j] = max(
                    dp[i - 1][j - 1] + nums1[i - 1] * nums2[j - 1], dp[i][j - 1], dp[i - 1][j]
                )

        return dp[-1][-1]

    def maxDotProduct2(self, nums1: List[int], nums2: List[int]) -> int:
        """
        视作类01背包问题，数组中每个元素只有选和不选两种状态
        两个数组组合就有四种状态
        """
        from functools import lru_cache

        @lru_cache(len(nums1) * len(nums2))
        def product(i, j):
            if i >= len(nums1) or j >= len(nums2):
                return 0
            # 两个都不选
            op1 = product(i + 1, j + 1)
            # 两个都选
            op2 = product(i + 1, j + 1) + nums1[i] * nums2[j]
            # 只选第一个
            op3 = product(i + 1, j)
            # 只选第二个
            op4 = product(i, j + 1)
            return max([op1, op2, op3, op4])

        ans = product(0, 0)
        if ans > 0:
            return ans
        else:
            return max(min(nums1) * max(nums2), max(nums1) * min(nums2))


print(Solution().maxDotProduct(nums1=[2, 1, -2, 5], nums2=[3, 0, -6]))
# 输出：18
# 解释：从 nums1 中得到子序列 [2,-2] ，从 nums2 中得到子序列 [3,-6] 。
# 它们的点积为 (2*3 + (-2)*(-6)) = 18 。
