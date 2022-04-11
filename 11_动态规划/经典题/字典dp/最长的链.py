from collections import defaultdict
from typing import List

# 首尾相同的区间可以相连 求最长的链
class Solution:
    def solve(self, blocks: List[List[int]]) -> int:
        """
        You can join two blocks if the end of one is equal to the start of another. 
        Return the length of the longest chain of blocks.
        """

        blocks.sort()
        # endswith maxLen
        dp = defaultdict(int)
        for s, e in blocks:
            dp[e] = max(dp[e], dp[s] + 1)
        return max(dp.values(), default=0)


print(Solution().solve(blocks=[[3, 4], [4, 5], [3, 7], [0, 1], [1, 3]]))
