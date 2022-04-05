# 统计数组中子序列'010'的个数和'101'的个数
class Solution:
    def numberOfWays(self, s: str) -> int:
        res = 0
        s0 = s1 = s10 = s01 = 0
        for char in s:
            if char == '1':
                s01 += s0
                s1 += 1
                res += s10
            else:
                s10 += s1
                s0 += 1
                res += s01
        return res

    def numberOfWays2(self, s: str) -> int:
        def numDistinct(s: str, t: str) -> int:
            """求s中有多少个子序列为t，时间复杂度O(st)"""

            if not t:
                return 0

            dp = [0] * (len(t) + 1)  # endswith dp
            dp[0] = 1

            for i in range(len(s)):
                # 注意要倒着推，避免有相同字母
                for j in reversed(range(len(t))):
                    if s[i] == t[j]:
                        dp[j + 1] += dp[j]

            return dp[-1]

        return numDistinct(s, '010') + numDistinct(s, '101')


print(Solution().numberOfWays2(s="001101"))
