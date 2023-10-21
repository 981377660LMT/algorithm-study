from typing import List

MOD = 10 ** 9 + 7
# 集团里有 n 名员工
# 第 i 种工作会产生 profit[i] 的利润，它要求 group[i] 名成员共同参与。
# 工作的任何至少产生 minProfit 利润的子集称为 盈利计划 。并且工作的成员总数最多为 n 。
# 有多少种计划可以选择？
# 1 <= n <= 100
# 100的数据 可能是O(n^3) 三个for循环
# https://leetcode-cn.com/problems/profitable-schemes/solution/c-python3-0-1bei-bao-wen-ti-qu-bie-pu-to-926y/


class Solution:
    def profitableSchemes(self, n: int, minProfit: int, group: List[int], profit: List[int]) -> int:
        dp = [[0 for _ in range(minProfit + 1)] for _ in range(n + 1)]

        # 前person个人，收益为0的方案有1种
        for person in range(n + 1):
            dp[person][0] = 1

        # 对每种工作(循环物品,二维费用)
        for needperson, getprof in zip(group, profit):
            # 对容量倒序(二维)
            for person in range(n, needperson - 1, -1):
                for prof in range(minProfit, -1, -1):
                    dp[person][prof] += dp[person - needperson][max(0, prof - getprof)]
                    dp[person][prof] %= MOD

        # 关于为什么是 max(0, money - prof) 而不是 for prof in range(minProfit,  getprof-1, -1):
        # 因为是至少获利minProfit, 前（people - needperson）个人获利至少是个负数的时候，说明前（people - needperson）个人可以摸鱼不干活儿（获利为0）

        return dp[n][minProfit]


print(Solution().profitableSchemes(n=10, minProfit=5, group=[2, 3, 5], profit=[6, 7, 8]))
# 输出：7
# 解释：至少产生 5 的利润，只要完成其中一种工作就行，所以该集团可以完成任何工作。
# 有 7 种可能的计划：(0)，(1)，(2)，(0,1)，(0,2)，(1,2)，以及 (0,1,2)
