# 给定一个整数 n，返回 下标从 1 开始 的数组 nums = [1, 2, ..., n]的 排列数，使其满足 自整除 条件。

# 如果对于每个 1 <= i <= n，至少 满足以下条件之一，数组 nums 就是 自整除 的：


# nums[i] % i == 0
# i % nums[i] == 0
# 数组的 排列 是对数组元素的重新排列的数量，例如，下面是数组 [1, 2, 3] 的所有排列：
from functools import lru_cache


class Solution:
    def selfDivisiblePermutationCount(self, n: int) -> int:
        @lru_cache(None)
        def dfs(index: int, visited: int) -> int:
            if index == n + 1:
                return 1

            res = 0
            for i in range(1, n + 1):
                if ((visited >> i) & 1 == 0) and (
                    ((((i) % (index)) == 0)) or (((index) % (i)) == 0)
                ):
                    res += dfs(index + 1, visited | (1 << i))

            return res

        res = dfs(1, 0)
        dfs.cache_clear()
        return res
