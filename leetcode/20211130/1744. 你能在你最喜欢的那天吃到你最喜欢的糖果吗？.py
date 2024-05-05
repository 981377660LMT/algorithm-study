from typing import List
from itertools import accumulate

# 你在吃完 所有 第 i - 1 类糖果之前，不能 吃任何一颗第 i 类糖果。
# 在吃完所有糖果之前，你必须每天 至少 吃 一颗 糖果。
# answer[i] 为 true 的条件是：在每天吃 不超过 dailyCapi 颗糖果的前提下，你可以在第 favoriteDayi 天吃到第 favoriteTypei 类糖果；


class Solution:
    def canEat(self, candiesCount: List[int], queries: List[List[int]]) -> List[bool]:
        pre = [0] + list(accumulate(candiesCount))
        return [pre[type] // cap <= day < pre[type + 1] for type, day, cap in queries]


print(
    Solution().canEat(
        candiesCount=[7, 4, 5, 3, 8], queries=[[0, 2, 2], [4, 2, 4], [2, 13, 1000000000]]
    )
)
# 输出：[true,false,true]
# 提示：
# 1- 在第 0 天吃 2 颗糖果(类型 0），第 1 天吃 2 颗糖果（类型 0），第 2 天你可以吃到类型 0 的糖果。
# 2- 每天你最多吃 4 颗糖果。即使第 0 天吃 4 颗糖果（类型 0），第 1 天吃 4 颗糖果（类型 0 和类型 1），你也没办法在第 2 天吃到类型 4 的糖果。换言之，你没法在每天吃 4 颗糖果的限制下在第 2 天吃到第 4 类糖果。
# 3- 如果你每天吃 1 颗糖果，你可以在第 13 天吃到类型 2 的糖果。

