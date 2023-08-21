from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n 表示数轴上的房屋数量，编号从 0 到 n - 1 。

# 另给你一个二维整数数组 offers ，其中 offers[i] = [starti, endi, goldi] 表示第 i 个买家想要以 goldi 枚金币的价格购买从 starti 到 endi 的所有房屋。

# 作为一名销售，你需要有策略地选择并销售房屋使自己的收入最大化。

# 返回你可以赚取的金币的最大数目。

# 注意 同一所房屋不能卖给不同的买家，并且允许保留一些房屋不进行出售。


class Solution:
    def maximizeTheProfit(self, n: int, offers: List[List[int]]) -> int:
        ...
