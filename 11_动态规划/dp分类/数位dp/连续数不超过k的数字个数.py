# return the number of strings of equal size as s (only consisting of lowercase letters)
# that are lexicographically smaller or equal to s and have at most k consecutive equal characters.
from functools import lru_cache


MOD = int(1e9 + 7)


class Solution:
    def solve(self, s: str, k: int) -> int:
        @lru_cache(None)
        def dfs(pos: int, isLimit: bool, pre: int, repeat: int) -> int:
            """当前在第pos位，isLimit表示是否贴合上界，pre表示之前的选择，repeat表示重复选择了几次"""
            if pos == n:
                return 1

            res = 0
            up = nums[pos] if isLimit else 25
            for cur in range(up + 1):
                if repeat == k and cur == pre:
                    continue
                res += dfs(pos + 1, (isLimit and cur == up), cur, 1 if cur != pre else repeat + 1)
                res %= MOD

            return res % MOD

        if k == 0:
            return 0

        n = len(s)
        nums = [ord(c) - ord('a') for c in s]
        return dfs(0, True, -1, 0)


print(Solution().solve("abb", 2))
print(Solution().solve("aaaz", 0))
