from functools import lru_cache


class Solution:
    def translateNum(self, num: int) -> int:
        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= n:
                return int(index == n)
            res = dfs(index + 1)
            num2 = int(s[index : index + 2])
            if 10 <= num2 <= 25:
                res += dfs(index + 2)

            return res

        s = str(num)
        n = len(s)
        return dfs(0)


print(Solution().translateNum(12258))
