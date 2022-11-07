# 击败所有怪兽的最小时间
# 二分图最小权匹配
# https://leetcode.cn/problems/minimum-time-to-kill-all-monsters/solution/er-fen-tu-by-skyflaming-uq0x/

from math import ceil
from typing import List
from KM算法模板 import KM

INF = int(1e18)


class Solution:
    def minimumTime(self, power: List[int]) -> int:
        n = len(power)
        costMatrix = [[0] * n for _ in range(n)]  # !左边点代表poweri,右边点代表排列里的元素j
        for i in range(n):
            for j in range(n):
                costMatrix[i][j] = -ceil(
                    power[i] / (j + 1)
                )  # !当前击败怪物花费 ceil(power[i] / (j + 1)) 时间
        km = KM(costMatrix)
        return -km.getResult()[0]

    # 状压的dp写法
    def minimumTime2(self, power: List[int]) -> int:
        n = len(power)
        target = (1 << n) - 1
        dp = [INF] * (target + 1)
        dp[0] = 0

        for state in range(target + 1):
            gain = state.bit_count() + 1
            for i in range(n):
                if state & (1 << i):
                    continue
                cost = ceil(power[i] / gain)
                dp[state | (1 << i)] = min(dp[state | (1 << i)], cost + dp[state])

        return dp[target]
