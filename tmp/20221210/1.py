from curses.ascii import isdigit
from heapq import nlargest
from typing import List


INF = int(1e20)
# 给你一个 n 个点的无向图，节点从 0 到 n - 1 编号。给你一个长度为 n 下标从 0 开始的整数数组 vals ，其中 vals[i] 表示第 i 个节点的值。
# 同时给你一个二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示节点 ai 和 bi 之间有一条双向边。
# 星图 是给定图中的一个子图，它包含一个中心节点和 0 个或更多个邻居。换言之，星图是给定图中一个边的子集，且这些边都有一个公共节点。
# 下图分别展示了有 3 个和 4 个邻居的星图，蓝色节点为中心节点。
# 星和 定义为星图中所有节点值的和。
# 给你一个整数 k ，请你返回 至多 包含 k 条边的星图中的 最大星和 。


class Solution:
    def maxStarSum(self, vals: List[int], edges: List[List[int]], k: int) -> int:
        n = len(vals)
        adjList = [[0] for _ in range(n)]
        for u, v in edges:
            if vals[v] > 0:
                adjList[u].append(vals[v])
            if vals[u] > 0:
                adjList[v].append(vals[u])

        # 枚举中心结点
        res = -INF
        for i in range(n):
            topk = sum(nlargest(k, adjList[i]))  # 注意边数限制 所以可以这样统计 (如果没有边数统计怎么做?)
            res = max(res, topk + vals[i])
        return res
