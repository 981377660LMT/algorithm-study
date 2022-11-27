# 6251-长度为5的回文子序列个数
# 1<=len(s)<=1e4
# !枚举中心+前后缀分解
# O(100*n 统计每个前后缀处的回文节个数)

from typing import DefaultDict, List
from itertools import product
from collections import defaultdict
import string

MOD = int(1e9 + 7)
HALF = ["".join(pair) for pair in product(string.digits, repeat=2)]


class Solution:
    def countPalindromes(self, s: str) -> int:
        def makeDp(s: str) -> List[DefaultDict[str, int]]:
            dp = [defaultdict(int) for _ in range(n + 1)]
            for i in range(1, n + 1):
                cur = s[i - 1]
                dp[i] = dp[i - 1].copy()
                dp[i][cur] += 1  # 单独选
                for pre in string.digits:  # 组成长为2的前缀
                    tmp = pre + cur
                    dp[i][tmp] += dp[i - 1][pre]
            return dp

        n = len(s)
        res, preDp, sufDp = 0, makeDp(s), makeDp(s[::-1])[::-1]
        for i in range(n):  # 枚举分割点
            for h in HALF:
                res += preDp[i][h] * sufDp[i + 1][h]
                res %= MOD
        return res


# print(Solution().countPalindromes(s="000000000" * 1000))
# print(Solution().countPalindromes(s="103301"))
print(Solution().countPalindromes(s="9999900000"))
