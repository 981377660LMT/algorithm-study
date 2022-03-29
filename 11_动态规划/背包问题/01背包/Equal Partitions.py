# return whether you can partition nums into two groups where the sum of the elements in both groups is equal.
# 是否能将数组分为两半，使得两半的和相等。
# n ≤ 250
# 1 ≤ nums[i] ≤ 100


# dfs写法
from functools import lru_cache


class Solution:
    def solve(self, nums):
        sum_ = sum(nums)
        if sum_ & 1:
            return False
        target = sum_ >> 1

        @lru_cache(None)
        def dfs(index, cur) -> bool:
            if cur > target:
                return False
            if index == len(nums):
                return not not cur == target

            res = False
            res = res or dfs(index + 1, cur + nums[index])
            res = res or dfs(index + 1, cur)
            return res

        return dfs(0, 0)


print(Solution().solve(nums=[1, 2, 5, 4]))

# dp标准写法


# 滚动集合写法

