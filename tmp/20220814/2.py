from typing import List, Tuple, Optional
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def edgeScore2(self, edges: List[int]) -> int:
        n = len(edges)
        score = defaultdict(int)
        for cur, next in enumerate(edges):
            score[next] += cur
        return max(range(n), key=lambda x: (score[x], -x))

    # 返回 边积分 最高的节点。如果多个节点的 边积分 相同，返回编号 最小 的那个。
    def edgeScore(self, edges: List[int]) -> int:
        """因为返回编号最小的,可以直接index"""
        res = [0] * (len(edges))
        for cur, next in enumerate(edges):
            res[next] += cur
        return res.index(max(res))
