from functools import lru_cache


class Solution:
    @lru_cache(None)
    def solve(self, n, k):
        if n <= 1:
            return int(k == 1)
        # (Either we put the smallest to the front and increase the number of visible blocks by one, or we put it after any of the blocks, to get the same number of visible blocks.)
        return self.solve(n - 1, k - 1) + (n - 1) * self.solve(n - 1, k)


print(Solution().solve(n=13, k=2))  # 5
