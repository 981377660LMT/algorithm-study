from functools import lru_cache

INF = int(1e20)
MOD = int(1e9 + 7)


# 给你整数 zero ，one ，low 和 high ，我们从空字符串开始构造一个字符串，每一步执行下面操作中的一种：
# 将 '0' 在字符串末尾添加 zero  次。
# 将 '1' 在字符串末尾添加 one 次。
# 以上操作可以执行任意次。
# !如果通过以上过程得到一个 长度 在 low 和 high 之间（包含上下边界）的字符串，那么这个字符串我们称为 好 字符串。
# 请你返回满足以上要求的 不同 好字符串数目。由于答案可能很大，请将结果对 109 + 7 取余 后返回。

# !爬楼梯加强版


class Solution:
    def countGoodStrings(self, low: int, high: int, zero: int, one: int) -> int:
        """dp+前缀和"""

        def cal(upper: int) -> int:
            """长度不超过upper的字符串的个数"""

            @lru_cache(None)
            def dfs(curLength: int) -> int:
                if curLength >= upper:
                    return 1 if curLength == upper else 0
                return (1 + dfs(curLength + zero) + dfs(curLength + one)) % MOD

            return dfs(0)

        return (cal(high) - cal(low - 1)) % MOD

    def countGoodStrings2(self, low: int, high: int, zero: int, one: int) -> int:
        dp = [0] * (high + 1)
        dp[0] = 1
        res = 0
        for i in range(1, high + 1):
            if i - zero >= 0:
                dp[i] = (dp[i] + dp[i - zero]) % MOD
            if i - one >= 0:
                dp[i] = (dp[i] + dp[i - one]) % MOD
            if i >= low:
                res = (res + dp[i]) % MOD
        return res


print(Solution().countGoodStrings(low=3, high=3, zero=1, one=1))
print(Solution().countGoodStrings2(low=3, high=3, zero=1, one=1))
