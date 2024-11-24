from typing import List
from collections import deque


class Edge:
    def __init__(self, to, rev, capacity):
        self.to = to
        self.rev = rev
        self.capacity = capacity


class MaxFlow:
    def __init__(self, N):
        self.size = N
        self.graph = [[] for _ in range(N)]

    def add(self, fr, to, capacity):
        forward = Edge(to, len(self.graph[to]), capacity)
        backward = Edge(fr, len(self.graph[fr]), 0)
        self.graph[fr].append(forward)
        self.graph[to].append(backward)

    def bfs_level(self, s, t, level):
        queue = deque()
        level[s] = 0
        queue.append(s)
        while queue:
            v = queue.popleft()
            for e in self.graph[v]:
                if e.capacity > 0 and level[e.to] < 0:
                    level[e.to] = level[v] + 1
                    queue.append(e.to)
        return level[t] != -1

    def dfs_flow(self, level, iter_, v, t, upTo):
        if v == t:
            return upTo
        for i in range(iter_[v], len(self.graph[v])):
            e = self.graph[v][i]
            if e.capacity > 0 and level[v] < level[e.to]:
                d = self.dfs_flow(level, iter_, e.to, t, min(upTo, e.capacity))
                if d > 0:
                    e.capacity -= d
                    self.graph[e.to][e.rev].capacity += d
                    return d
            iter_[v] += 1
        return 0

    def max_flow(self, s, t):
        flow = 0
        level = [-1] * self.size
        INF = float("inf")
        while True:
            level = [-1] * self.size
            if not self.bfs_level(s, t, level):
                break
            iter_ = [0] * self.size
            while True:
                f = self.dfs_flow(level, iter_, s, t, INF)
                if f == 0:
                    break
                flow += f
        return flow


class Solution:
    def maxRemoval(self, nums: List[int], queries: List[List[int]]) -> int:
        n = len(nums)
        m = len(queries)

        counts = [0] * n
        index_to_queries = [[] for _ in range(n)]
        for idx, (l, r) in enumerate(queries):
            for i in range(l, r + 1):
                counts[i] += 1
                index_to_queries[i].append(idx)

        for i in range(n):
            if counts[i] < nums[i]:
                return -1

        S = n + m
        T = n + m + 1
        size = n + m + 2
        mf = MaxFlow(size)

        for i in range(n):
            mf.add(S, i, nums[i])

        for idx, (l, r) in enumerate(queries):
            q_node = n + idx
            for i in range(l, r + 1):
                mf.add(i, q_node, 1)

        for idx in range(m):
            q_node = n + idx
            mf.add(q_node, T, 1)

        flow = mf.max_flow(S, T)

        max_deletable_queries = m - flow
        return max_deletable_queries
