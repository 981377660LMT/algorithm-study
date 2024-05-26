from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 limit 和一个大小为 n x 2 的二维数组 queries 。

# 总共有 limit + 1 个球，每个球的编号为 [0, limit] 中一个 互不相同 的数字。一开始，所有球都没有颜色。queries 中每次操作的格式为 [x, y] ，你需要将球 x 染上颜色 y 。每次操作之后，你需要求出所有球中 不同 颜色的数目。

# 请你返回一个长度为 n 的数组 result ，其中 result[i] 是第 i 次操作以后不同颜色的数目。


# 注意 ，没有染色的球不算作一种颜色。


class Solution:
    def queryResults(self, limit: int, queries: List[List[int]]) -> List[int]:
        counter = defaultdict(int)
        freqCounter = defaultdict(int)
