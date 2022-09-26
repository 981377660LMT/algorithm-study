from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 为了缓解「力扣嘉年华」期间的人流压力，组委会在活动期间开设了一些交通专线。
# path[i] = [a, b] 表示有一条从地点 a通往地点 b 的 单向 交通专线。
# 若存在一个地点，满足以下要求，我们则称之为 交通枢纽：

# 所有地点（除自身外）均有一条 单向 专线 直接 通往该地点；
# 该地点不存在任何 通往其他地点 的单向专线。
# 请返回交通专线的 交通枢纽。若不存在，则返回 -1。


class Solution:
    def transportationHub(self, path: List[List[int]]) -> int:
        adjMap = defaultdict(set)
        indeg = defaultdict(int)
        outdeg = defaultdict(int)
        allVertex = set()
        for a, b in path:
            adjMap[a].add(b)
            outdeg[a] += 1
            indeg[b] += 1
            allVertex.add(a)
            allVertex.add(b)

        n = len(allVertex)
        for v in allVertex:
            if indeg[v] == n - 1 and outdeg[v] == 0:
                return v

        return -1
