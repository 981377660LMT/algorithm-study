# In one operation you can either

# Decrement n by one
# If n is even, decrement by n / 2
# If n is divisible by 3, decrement by 2 * (n / 3)

# 求到0的最少操作次数
import functools


class Solution:
    @functools.lru_cache(None)
    def solve(self, n):
        if n <= 1:
            return n
        return min(1 + int(n % 3) + self.solve(n // 3), 1 + int(n % 2) + self.solve(n // 2))


# 也可以bfs
