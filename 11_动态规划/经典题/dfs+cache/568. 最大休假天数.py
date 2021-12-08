from typing import List
from functools import lru_cache

# ，flights[i][j] 代表城市i到城市j的航空状态;这些航班用 N*N 矩阵表示(有向图)
# 您每天最多只能乘坐一次航班，并且只能在每周的星期一上午乘坐航班。
# days[i][j] 代表您在第j个星期在城市i能休假的最长天数。
# 给定 flights 矩阵和 days 矩阵，您需要输出 K 周内可以休假的最长天数。
# N 和 K 都是正整数，在 [1, 100] 范围内。
class Solution:
    def maxVacationDays(self, flights: List[List[int]], days: List[List[int]]) -> int:
        weight = flights
        n = len(weight)  # 结点数
        k = len(days[0])  # 星期数

        @lru_cache(None)
        def dfs(city: int, day: int) -> int:
            if day == k:
                return 0

            res = 0
            for nextCity in range(n):
                tmp = 0
                if weight[city][nextCity] == 1 or city == nextCity:
                    tmp = days[nextCity][day] + dfs(nextCity, day + 1)
                res = max(res, tmp)
            return res

        return dfs(0, 0)


print(
    Solution().maxVacationDays(
        flights=[[0, 1, 1], [1, 0, 1], [1, 1, 0]], days=[[1, 3, 1], [6, 0, 3], [3, 3, 3]]
    )
)
# 输出: 12
# 解释:
# Ans = 6 + 3 + 3 = 12.

# 最好的策略之一：
# 第一个星期 : 星期一从城市0飞到城市1，玩6天，工作1天。
# （虽然你是从城市0开始，但因为是星期一，我们也可以飞到其他城市。）
# 第二个星期 : 星期一从城市1飞到城市2，玩3天，工作4天。
# 第三个星期 : 呆在城市2，玩3天，工作4天。

