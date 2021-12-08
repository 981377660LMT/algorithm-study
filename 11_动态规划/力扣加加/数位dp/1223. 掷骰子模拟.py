from typing import List

# 投掷骰子时，连续 掷出数字 i 的次数不能超过 rollMax[i]
# 计算掷 n 次骰子可得到的不同点数序列的数量。

MOD = int(1e9 + 7)


# dp[i][j][k] 表示第 i 轮掷骰子掷出数字 j 时 j 连续出现 k 次的组合数
class Solution:
    def dieSimulator(self, n: int, rollMax: List[int]) -> int:
        maxFreq = max(rollMax)
        dp = [[[0 for _ in range(maxFreq + 1)] for _ in range(7)] for _ in range(n + 1)]
        for j in range(1, 7):
            dp[1][j][1] = 1

        for i in range(2, n + 1):
            for j in range(1, 7):
                # pre is j
                for k in range(2, rollMax[j - 1] + 1):
                    dp[i][j][k] = dp[i - 1][j][k - 1]

                # pre is not j
                uniCount = 0
                for pre in range(1, 7):
                    if pre == j:
                        continue
                    for k in range(1, rollMax[pre - 1] + 1):
                        uniCount += dp[i - 1][pre][k]
                        uniCount %= MOD
                dp[i][j][1] = uniCount

        res = 0
        for j in range(1, 7):
            res += sum(dp[n][j])
            res %= MOD
        return res


print(Solution().dieSimulator(n=2, rollMax=[1, 1, 2, 2, 2, 3]))
# 输出：34
# 解释：我们掷 2 次骰子，如果没有约束的话，共有 6 * 6 = 36 种可能的组合。但是根据 rollMax 数组，数字 1 和 2 最多连续出现一次，所以不会出现序列 (1,1) 和 (2,2)。因此，最终答案是 36-2 = 34。

