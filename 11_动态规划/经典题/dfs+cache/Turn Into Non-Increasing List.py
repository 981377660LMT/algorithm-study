from functools import lru_cache
from typing import List

# https://binarysearch.com/problems/Turn-Into-Non-Increasing-List
class Solution:
    def solve(self, nums: List[int]) -> int:
        """Return the minimum number of operations required so that the list becomes non-increasing."""

        @lru_cache(None)
        def dfs(index: int, pre: int) -> int:
            ...


print(Solution().solve(nums=[1, 5, 3, 9, 1]))
# We can merge [1, 5] to get [6, 3, 9, 1] and then merge [6, 3] to get [9, 9, 1].


# todo
