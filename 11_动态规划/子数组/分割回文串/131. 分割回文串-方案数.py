from functools import lru_cache


class Solution:
    @lru_cache(None)
    def solve(self, s):
        n = len(s)
        if n <= 1:
            return 1

        res = 0
        for i in range(1, n + 1):
            substring = s[:i]
            if substring == substring[::-1]:
                res += self.solve(s[i:])

        return res
