from typing import List
from math import sqrt


# 保证递增
class Solution:
    def getFactors(self, n: int) -> List[List[int]]:
        """n所有的因子组合,因子必须大于 1 并且小于 n"""

        def dfs(lower: int, upper: int) -> List[List[int]]:
            """枚举每个组合的开头元素cur 需要遍历因子 所以需要携带upper参数"""
            res = []
            for cur in range(lower, int(sqrt(upper)) + 1):
                if upper % cur == 0:
                    res.append([cur, upper // cur])
                    for next in dfs(cur, upper // cur):
                        res.append([cur, *next])
            return res

        return dfs(2, n)


print(Solution().getFactors(12))
# [
#   [2, 6],
#   [2, 2, 3],
#   [3, 4]
# ]
