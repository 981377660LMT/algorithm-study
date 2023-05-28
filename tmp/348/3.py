from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始、长度为 n 的二进制字符串 s ，你可以对其执行两种操作：

# 选中一个下标 i 并且反转从下标 0 到下标 i（包括下标 0 和下标 i ）的所有字符，成本为 i + 1 。
# 选中一个下标 i 并且反转从下标 i 到下标 n - 1（包括下标 i 和下标 n - 1 ）的所有字符，成本为 n - i 。
# 返回使字符串内所有字符 相等 需要的 最小成本 。

# 反转 字符意味着：如果原来的值是 '0' ，则反转后值变为 '1' ，反之亦然。


# !前缀相等+后缀相等
# dp??


class Solution:
    def minimumCost(self, s: str) -> int:
        if len(s) == 1:
            return 0
        res = self.solve(s)
        sb = ["0" if c == "1" else "1" for c in s]
        s = "".join(sb)
        res = min(res, self.solve(s))
        return res

    def solve(self, s: str) -> int:
        """全变为1的最小成本."""

        def getDp(s: str) -> List[int]:
            """dp[i]表示s[:i+1]变为全1的最小成本."""
            n = len(s)
            dp = [0] * n
            dp[0] = 0 if s[0] == "1" else 1
            for i in range(1, n):
                if s[i] == "1":
                    dp[i] = dp[i - 1]
                else:
                    if s[i - 1] == "1":
                        dp[i] = dp[i - 1] + 2 * i + 1
                    else:
                        dp[i] = dp[i - 1] + 1
            return dp

        n = len(s)
        dp1 = getDp(s)
        dp2 = getDp(s[::-1])[::-1]
        res = min(dp1[i] + dp2[i + 1] for i in range(n - 1))
        res = min(res, dp1[-1])
        res = min(res, dp2[0])
        return res


# s = "0011"
# s = "010101"
print(Solution().minimumCost("0011"))
print(Solution().minimumCost("010101"))
print(Solution().minimumCost("10101"))
print(Solution().minimumCost("1"))
