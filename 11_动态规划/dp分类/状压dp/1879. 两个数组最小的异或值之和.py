from typing import List
from functools import lru_cache

# 1 <= n <= 14
# Time complexity is O(n*2^n), because we have 2^n states and O(n) transactions frome one state to others. Space complexity is O(2^n).


class Solution:
    def minimumXORSum(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)

        @lru_cache(None)
        def dfs(index: int, state: int) -> int:
            if index == n:
                return 0

            res = int(1e15)
            for next in range(n):
                if not state & (1 << next):
                    res = min(
                        res, (nums1[index] ^ nums2[next]) + dfs(index + 1, state | (1 << next))
                    )
            return res

        return dfs(0, 0)


print(Solution().minimumXORSum(nums1=[1, 2], nums2=[2, 3]))
# 输出：2
# 解释：将 nums2 重新排列得到 [3,2] 。
# 异或值之和为 (1 XOR 3) + (2 XOR 2) = 2 + 0 = 2 。
