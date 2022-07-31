from functools import lru_cache
from typing import List

INF = int(1e20)
# n ≤ 1,000

# 时间复杂度O(n)


class Solution:
    def minSwap(self, nums1: List[int], nums2: List[int]) -> int:
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
print(Solution().minSwap(nums1=[1, 3, 5, 4], nums2=[1, 2, 3, 7]))
