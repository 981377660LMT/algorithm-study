from typing import List

MOD = 10 ** 9 + 7
# 集团里有 n 名员工
# 第 i 种工作会产生 profit[i] 的利润，它要求 group[i] 名成员共同参与。
# 工作的任何至少产生 minProfit 利润的子集称为 盈利计划 。并且工作的成员总数最多为 n 。
# 有多少种计划可以选择？
# 1 <= n <= 100
# 100的数据 可能是O(n^3) 三个for循环
# https://leetcode-cn.com/problems/profitable-schemes/solution/c-python3-0-1bei-bao-wen-ti-qu-bie-pu-to-926y/

# 二维费用的01背包
# 0. 人数为背包容量即第一维的限制条件，利润为背包第二维的限制条件，工作为待选取的物品，每个物品有他的`(cost,score)`
# 1. 只有完全背包的容量是从小到大循环的
# 2. 循环顺序 `物品 => 体积 => 决策`


MOD = int(1e9 + 7)


class Solution:
    def profitableSchemes(self, n: int, minProfit: int, group: List[int], profit: List[int]) -> int:
        goods = []
        for cost, score in zip(group, profit):
            goods.append((cost, score))

        dp = [[0 for _ in range(minProfit + 1)] for _ in range(n + 1)]
        for cap in range(n + 1):
            dp[cap][0] = 1

        for cost, score in goods:
            for i in range(n, cost - 1, -1):
                for j in range(minProfit, -1, -1):
                    # 负数表示超出需要，已经满足了条件，所以取max(0, j - score)
                    dp[i][j] += dp[i - cost][max(0, j - score)]
                    dp[i][j] %= MOD

        return dp[-1][-1]


print(Solution().profitableSchemes(n=10, minProfit=5, group=[2, 3, 5], profit=[6, 7, 8]))
# 输出：7
# 解释：至少产生 5 的利润，只要完成其中一种工作就行，所以该集团可以完成任何工作。
# 有 7 种可能的计划：(0)，(1)，(2)，(0,1)，(0,2)，(1,2)，以及 (0,1,2)
