from functools import lru_cache
from typing import List, Literal


class Solution:
    def maximumsSplicedArray(self, nums1: List[int], nums2: List[int]) -> int:
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


print(Solution().maximumsSplicedArray(nums1=[60, 60, 60], nums2=[10, 90, 10]))
print(
    Solution().maximumsSplicedArray(nums1=[10, 20, 50, 15, 30, 10], nums2=[40, 20, 10, 100, 10, 10])
)


class Solution2:
    def maximumsSplicedArray(self, nums1: List[int], nums2: List[int]) -> int:

        diff = [nums2[j] - nums1[j] for j in range(len(nums1))]
        diff2 = [nums1[j] - nums2[j] for j in range(len(nums1))]

        dp = 0
        cur = 0
        for num in diff:
            cur = max(num, cur + num)
            dp = max(dp, cur)

        dp2 = 0
        cur = 0
        for num in diff2:
            cur = max(num, cur + num)
            dp2 = max(dp2, cur)

        return max(sum(nums1) + dp, sum(nums2) + dp2)


# class Solution:
#     def maximumsSplicedArray(self, nums1: List[int], nums2: List[int]) -> int:
#         arr = [x - y for x, y in zip(nums1, nums2)]
#         p1, p2 = inf, -inf
#         mx = -inf
#         mn = inf
#         for x in accumulate(arr):
#             p1 = min(p1, x)
#             p2 = max(p2, x)
#             mx = max(x - p1, mx)
#             mn = min(x - p2, mn)
#         return max(sum(nums2) + mx, sum(nums1) - mn)
