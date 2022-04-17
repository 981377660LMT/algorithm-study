from typing import List
from collections import defaultdict

# 要求保留连通性，很容易想到点双连通分量和割点。
# 为了防止魔物暴动，勇者在每一次夺回据点后（包括花费资源夺回据点后），
# 需要保证剩余的所有魔物据点之间是相连通的（不经过「已夺回据点」）。


class Solution:
    def minimumCost(self, cost: List[int], roads: List[List[int]]) -> int:
        n = len(cost)
        adjMap = defaultdict(set)
        for u, v in roads:
            adjMap[u].add(v)
            adjMap[v].add(u)

