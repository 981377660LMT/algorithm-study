from typing import List, Tuple


# 1 <= numLaps <= 1000
# 1 <= tires.length <= 105

# 总结：
# 这道题一开始想贪心的解法(贪心ptsd)，sortedList弄了好久，
# 最后才意识到是dp 状态由圈数唯一决定 但是怎么求每个圈的最小时间花费呢?


# 总结:很明显贪心不对的(举反例),就不要贪心了,考虑别的解法,一般是dp,找dfs的自变量是什么,怎么转移,初始值是什么
# 实际上是20个完全背包,凑成numLaps的容量,看最少花费

# 1 <= tires.length <= 105
# 1 <= numLaps <= 1000
# 2 <= ri <= 105
class Solution:
    def minimumFinishTime(self, tires: List[List[int]], changeTime: int, numLaps: int) -> int:
        """tires[i] = [fi, ri] 表示第 i 种轮胎如果连续使用，第 x 圈需要耗时 fi * ri(x-1) 秒"""
        """每一圈后，你可以选择耗费 changeTime 秒 换成 任意一种轮胎（也可以换成当前种类的新轮胎）。"""
        # 预处理出不换轮胎,连续使用同一个轮胎跑 xx 圈的最小耗时 即每个物品的价格
        # 状态转移 每个圈为状态 转移为下一次连续用多少个轮胎
        prices = [int(1e20)] * 20
        for a0, q in tires:
            curSum, curItem = 0, a0
            for i in range(len(prices)):
                curSum += curItem
                if curSum > int(1e5):
                    break
                prices[i] = min(prices[i], curSum)
                curItem *= q

        # dp[i]表示第i圈的最小耗时
        dp = [int(1e20)] * (numLaps + 1)
        dp[0] = 0
        for i in range(1, numLaps + 1):
            for j in range(len(prices)):
                if i - (j + 1) >= 0:
                    dp[i] = min(dp[i], dp[i - (j + 1)] + prices[j] + changeTime)

        # 减去最后一次的换轮胎耗时
        return dp[-1] - changeTime


# 21 25
print(Solution().minimumFinishTime(tires=[[2, 3], [3, 4]], changeTime=5, numLaps=4))
print(Solution().minimumFinishTime(tires=[[1, 10], [2, 2], [3, 4]], changeTime=6, numLaps=5))

