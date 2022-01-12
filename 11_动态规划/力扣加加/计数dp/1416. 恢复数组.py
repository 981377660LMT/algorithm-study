from functools import lru_cache

# 某个程序本来应该输出一个整数数组。但是这个程序忘记输出空格了以致输出了一个数字字符串，
# 我们所知道的信息只有：数组中所有整数都在 [1, k] 之间，且数组中的数字都没有前导 0 。
# 1 <= s.length <= 10^5.

# lru 超时 改用手动memo ？？？

MOD = int(1e9 + 7)


class Solution:
    def numberOfArrays(self, s: str, k: int) -> int:
        # 以i为字符串起点的方案数
        @lru_cache(typed=False, maxsize=len(s))
        def dfs(i):
            if i == len(s):
                return 1
            if s[i] == '0':
                return 0

            res = 0
            for j in range(i, len(s)):
                if int(s[i : j + 1]) > k:
                    break
                res += dfs(j + 1) % MOD
            return res

        res = dfs(0)
        dfs.cache_clear()

        return res % MOD


print(Solution().numberOfArrays(s="1317", k=2000))
# 输出：8
# 解释：可行的数组方案为 [1317]，[131,7]，[13,17]，[1,317]，[13,1,7]，[1,31,7]，[1,3,17]，[1,3,1,7]

