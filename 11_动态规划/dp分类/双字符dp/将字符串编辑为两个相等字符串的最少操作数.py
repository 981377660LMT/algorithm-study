# Return the minimum number of operations required such that s = t + t for any t.
# 0 ≤ n ≤ 100
class Solution:
    def solve(self, s):
        def countEditDist(s1, s2) -> int:
            """编辑距离+枚举分割点"""
            n1, n2 = len(s1), len(s2)
            dp = [[0] * (n2 + 1) for _ in range(n1 + 1)]
            for i in range(n1 + 1):
                dp[i][0] = i
            for j in range(n2 + 1):
                dp[0][j] = j

            for i in range(1, n1 + 1):
                for j in range(1, n2 + 1):
                    if s1[i - 1] == s2[j - 1]:
                        dp[i][j] = dp[i - 1][j - 1]
                    else:
                        dp[i][j] = min(dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1]) + 1

            return dp[-1][-1]

        res = len(s)
        # 枚举分割点
        for i in range(len(s)):
            res = min(countEditDist(s[:i], s[i:]), res)
        return res


print(Solution().solve(s="abczbdc"))
