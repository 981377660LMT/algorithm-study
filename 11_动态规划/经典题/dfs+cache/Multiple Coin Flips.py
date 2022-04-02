# k个硬币朝上的概率


from functools import lru_cache


class Solution:
    def solve(self, chances, k):
        @lru_cache(None)
        def dfs(index: int, count: int) -> int:
            if count > k:
                return 0
            if index == len(chances):
                return int(count == k)

            skip = (1 - chances[index]) * dfs(index + 1, count)
            choose = chances[index] * dfs(index + 1, count + 1)
            return skip + choose

        return dfs(0, 0)


print(Solution().solve(chances=[0.5, 0.4], k=2))
