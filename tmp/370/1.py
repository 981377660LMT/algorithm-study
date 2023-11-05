from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 一场比赛中共有 n 支队伍，按从 0 到  n - 1 编号。每支队伍也是 有向无环图（DAG） 上的一个节点。

# 给你一个整数 n 和一个下标从 0 开始、长度为 m 的二维整数数组 edges 表示这个有向无环图，其中 edges[i] = [ui, vi] 表示图中存在一条从 ui 队到 vi 队的有向边。

# 从 a 队到 b 队的有向边意味着 a 队比 b 队 强 ，也就是 b 队比 a 队 弱 。

# 在这场比赛中，如果不存在某支强于 a 队的队伍，则认为 a 队将会是 冠军 。


# 如果这场比赛存在 唯一 一个冠军，则返回将会成为冠军的队伍。否则，返回 -1 。
class Solution:
    def findChampion(self, n: int, edges: List[List[int]]) -> int:
        adjList = [[] for _ in range(n)]
        indeg = [0] * n
        for u, v in edges:
            adjList[u].append(v)
            indeg[v] += 1
        zeroCount = indeg.count(0)
        if zeroCount != 1:
            return -1
        return indeg.index(0)
