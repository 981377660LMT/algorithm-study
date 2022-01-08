# 依赖于前和下 可以从下到上、从前到后 或者 从前到后、从下到上
# 也可以 长度从短到长，相当于斜着来
from typing import List
from functools import lru_cache


class Solution:
    def PredictTheWinner(self, nums: List[int]) -> bool:
        if len(nums) == 1:
            return True

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left == right:
                return nums[left]
            return max(nums[left] - dfs(left + 1, right), nums[right] - dfs(left, right - 1))

        return dfs(0, len(nums) - 1) >= 0


print(Solution().PredictTheWinner([1, 5, 2]))
print(Solution().PredictTheWinner([1, 5, 233, 7]))
