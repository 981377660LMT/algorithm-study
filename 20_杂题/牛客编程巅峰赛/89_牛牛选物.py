from typing import List

# n<=20
# 体积，重量<=10^9
# 现在有n个物品，每个物品有一个体积v[i]和重量g[i],选择其中总体积恰好为V的若干个物品，
# 想使这若干个物品的总重量最大，求最大总重量为多少。（如果不存在合法方案，返回-1）
class Solution:
    def Maximumweight(self, v: List[int], g: List[int], V: int) -> int:
        """n<=20 枚举即可"""
        ...
