from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumCandies(self, candies: List[int], k: int) -> int:
        def check(mid):
            """k个孩子 每个小孩都mid"""
            childCount = 0
            for candy in candies:
                childCount += candy // mid
            return childCount >= k

        # max_=len(candies)
        left, right = 1, int(1e14)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maximumCandies(candies=[5, 8, 6], k=3))
print(Solution().maximumCandies(candies=[2, 5], k=11))
print(
    Solution().maximumCandies(
        candies=[8009354, 3742360, 6196357, 5769413, 9681885, 2583391], k=25111650
    )
)

