from functools import lru_cache
from math import comb


# 由于古董键盘年久失修，键盘上只有 26 个字母 a~z 可以按下，
# 且每个字母最多仅能被按 k 次。
# 小扣随机按了 n 次按键，请返回小扣总共有可能按出多少种内容

# 1 <= k <= 5
# 1 <= n <= 26*k

MOD = int(1e9 + 7)
comb = lru_cache(comb)


class Solution:
    def keyboard(self, k: int, n: int) -> int:
        """对每个字母，看选几个，放在那几个位置"""

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain < 0:
                return 0
            if index == 26:
                return 1 if remain == 0 else 0

            res = 0
            for select in range(min(k, remain) + 1):
                res += comb(remain, select) * dfs(index + 1, remain - select)
                res %= MOD
            return res

        res = dfs(0, n)
        dfs.cache_clear()
        return res


print(Solution().keyboard(k=1, n=2))

# 输出：650
# 解释：由于只能按两次按键，且每个键最多只能按一次，所有可能的字符串（按字典序排序）为 "ab", "ac", ... "zy"
