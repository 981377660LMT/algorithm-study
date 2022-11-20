# 1977. 划分数字的方案数
# 你写下了若干 正整数 ，并将它们连接成了一个字符串 num 。
# 但是你忘记给这些数字之间加逗号了。
# !你只记得这一列数字是 非递减 的且 没有 任何数字有前导 0 。
# 请你返回有多少种可能的 正整数数组 可以得到字符串 num 。
# 由于答案可能很大，将结果对 1e9 + 7 取余 后返回。
# 1 <= num.length <= 3500
# num 只含有数字 '0' 到 '9' 。


# https://leetcode.cn/problems/number-of-ways-to-separate-numbers/solution/by-iancn-90gl/
# dp[i][j]代表num[:i]中最后一个数长度为j的方案数量
# !因为范围为长度,所以把长度作为dp数组的第二维
MOD = int(1e9 + 7)


class Solution:
    def numberOfCombinations(self, s: str) -> int:
        n = len(s)
        # !dp[index][len] 表示 前index个数中,最后一个数长度为len的方案数量
        dp = [[0] * (i + 1) for i in range(n + 1)]
        dp[0][0] = 1
        dpSum = [[0] * (i + 1) for i in range(n + 1)]
        dpSum[0][0] = 1

        for i in range(1, n + 1):
            for j in range(1, i + 1):
                if s[i - j] != "0":  # !这一组的开头
                    # 这组的长度小于上一组
                    # dp[i][j] = sum(dp[i - j][k] for k in range(min(j, i - j + 1)))
                    dp[i][j] = dpSum[i - j][min(j - 1, i - j)]
                    # 这组的长度等于上一组
                    if i - 2 * j >= 0 and s[i - j : i] >= s[i - 2 * j : i - j]:
                        dp[i][j] += dp[i - j][j]
                    dp[i][j] %= MOD
                dpSum[i][j] = (dpSum[i][j - 1] + dp[i][j]) % MOD

        return dpSum[-1][-1]
