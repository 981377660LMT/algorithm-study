"""使序列递增的最小操作次数"""

# n ≤ 1e5
# 时间复杂度O(n)
# !dp[index][preChange] 表示前index个数，最后一个数是否改变过的最小操作次数


from functools import lru_cache
from typing import List

INF = int(1e20)


class Solution:
    def minSwap(self, A: List[int], B: List[int]) -> int:
        """
        返回 使 nums1 和 nums2 严格递增 所需操作的最小次数 。
        用例保证可以实现操作。
        """
        n = len(A)
        dp = [0, 1]  # 不改变最后一个数、改变最后一个数
        for i in range(1, n):
            ndp = [INF, INF]
            for preSwap in range(2):
                preA, preB = (A[i - 1], B[i - 1]) if preSwap == 0 else (B[i - 1], A[i - 1])
                for curSwap in range(2):
                    curA, curB = (A[i], B[i]) if curSwap == 0 else (B[i], A[i])
                    if curA > preA and curB > preB:
                        ndp[curSwap] = min(ndp[curSwap], dp[preSwap] + curSwap)
            dp = ndp

        return min(dp)

    def minSwap2(self, nums1: List[int], nums2: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, isPreChanged: bool) -> int:
            if index == len(nums1):
                return 0

            preA, preB = -INF, -INF
            if index > 0:
                preA, preB = (
                    [nums1[index - 1], nums2[index - 1]]
                    if not isPreChanged
                    else [nums2[index - 1], nums1[index - 1]]
                )

            res = INF
            if nums1[index] > preA and nums2[index] > preB:
                res = min(res, dfs(index + 1, False))
            if nums1[index] > preB and nums2[index] > preA:
                res = min(res, dfs(index + 1, True) + 1)
            return res

        res = dfs(0, False)
        dfs.cache_clear()
        return res


# 返回 使 nums1 和 nums2 严格递增 所需操作的最小次数 。
# 用例保证可以实现操作。
print(Solution().minSwap(A=[1, 4], B=[3, 2]))
print(Solution().minSwap(A=[1, 3, 5, 4], B=[1, 2, 3, 7]))
print(Solution().minSwap2(nums1=[1, 3, 5, 4], nums2=[1, 2, 3, 7]))
