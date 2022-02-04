# 给出两个正整数x，y，另给出若干个数对[ai,bi,ci]，
# 请挑选若干数对使得挑出的数对ai的和不小于x，bi的和不小于y，
# 计算挑出数对的ci的和的最小值

# breadNum,beverageNum<=2000 暗示状态就是这两维
# 可以用一个二维数组来表示状态，
# 其中第一维表示面包的数量，第二维表示饮料的数量


INF = 0x3F3F3F3F


class Solution:
    def minCost(self, breadNum, beverageNum, packageSum):
        n = len(packageSum)
        dp = [[INF] * (beverageNum + 1) for _ in range(breadNum + 1)]
        dp[0][0] = 0

        for k in range(n):
            num1, num2, cost = packageSum[k]
            for i in range(breadNum, -1, -1):
                for j in range(beverageNum, -1, -1):
                    dp[i][j] = min(dp[i][j], dp[max(0, i - num1)][max(0, j - num2)] + cost)

        return dp[-1][-1]


print(
    Solution().minCost(
        5, 60, [[3, 36, 120], [10, 25, 129], [5, 50, 250], [1, 45, 130], [4, 20, 119]]
    )
)
