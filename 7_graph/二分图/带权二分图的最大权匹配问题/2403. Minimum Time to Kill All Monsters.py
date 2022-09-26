# 击败所有怪兽的最小时间
# 二分图最小权匹配
# https://leetcode.cn/problems/minimum-time-to-kill-all-monsters/solution/er-fen-tu-by-skyflaming-uq0x/

from math import ceil
from typing import List
from KM算法模板 import KM


class Solution:
    def minimumTime(self, power: List[int]) -> int:
        n = len(power)
        adjMatrix = [[0] * n for _ in range(n)]  # !左边点代表poweri,右边点代表排列里的元素j
        for i in range(n):
            for j in range(n):
                adjMatrix[i][j] = -ceil(power[i] / (j + 1))  # !当前击败怪物花费 ceil(power[i] / (j + 1)) 时间
        km = KM(adjMatrix)
        return -km.getResult()
