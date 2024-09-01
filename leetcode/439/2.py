from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 有一个无限大的二维平面。

# 给你一个正整数 k ，同时给你一个二维数组 queries ，包含一系列查询：

# queries[i] = [x, y] ：在平面上坐标 (x, y) 处建一个障碍物，数据保证之前的查询 不会 在这个坐标处建立任何障碍物。
# 每次查询后，你需要找到离原点第 k 近 障碍物到原点的 距离 。

# 请你返回一个整数数组 results ，其中 results[i] 表示建立第 i 个障碍物以后，离原地第 k 近障碍物距离原点的距离。如果少于 k 个障碍物，results[i] == -1 。

# 注意，一开始 没有 任何障碍物。


# 坐标在 (x, y) 处的点距离原点的距离定义为 |x| + |y| 。
class Solution:
    def resultsArray(self, queries: List[List[int]], k: int) -> List[int]:
        sl = SortedList()
        res = []
        for x, y in queries:
            sl.add(abs(x) + abs(y))
            if len(sl) < k:
                res.append(-1)
            else:
                res.append(sl[k - 1])
        return res
