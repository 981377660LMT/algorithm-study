from typing import List, Tuple
from heapq import heappop, heappush

Cost, ID, Time = int, int, int
#  edges[i] = [xi, yi, timei] 表示城市 xi 和 yi 之间有一条双向道路，耗费时间为 timei 分钟
# 你在城市 0 ，你想要在 maxTime 分钟以内 （包含 maxTime 分钟）到达城市 n - 1
# 请你返回完成旅行的 最小费用 ，如果无法在 maxTime 分钟以内完成旅行，请你返回 -1 。

# 带限制的单源最短路径问题


class Solution:
    def minCost2(self, maxTime: int, edges: List[List[int]], passingFees: List[int]) -> int:
        n = len(passingFees)
        adjList = [[] for _ in range(n)]
        for u, v, t in edges:
            adjList[u].append((v, t))
            adjList[v].append((u, t))

        pq: List[Tuple[Cost, ID, Time]] = [(passingFees[0], 0, 0)]
        visited = [0x7FFFFFFF for _ in range(n)]

        while pq:
            cur_cost, cur_id, cur_time = heappop(pq)
            if cur_time > maxTime:
                continue
            if cur_id == n - 1:
                return cur_cost
            if cur_time < visited[cur_id]:
                visited[cur_id] = cur_time
                for next_id, next_time in adjList[cur_id]:
                    heappush(pq, (cur_cost + passingFees[next_id], next_id, cur_time + next_time))
        return -1

    def minCost(self, maxTime: int, edges: List[List[int]], passingFees: List[int]) -> int:
        n = len(passingFees)
        adjList = [[] for _ in range(n)]
        for u, v, t in edges:
            adjList[u].append((v, t))
            adjList[v].append((u, t))

        pq: List[Tuple[Cost, ID, Time]] = [(passingFees[0], 0, 0)]
        dist = [[0x7FFFFFFF] * (maxTime + 1) for _ in range(n)]

        while pq:
            cur_cost, cur_id, cur_time = heappop(pq)
            if cur_time > maxTime:
                continue

            if cur_id == n - 1:
                return cur_cost

            if cur_cost < dist[cur_id][cur_time]:
                dist[cur_id][cur_time] = cur_cost
                for next_id, next_time in adjList[cur_id]:
                    heappush(pq, (cur_cost + passingFees[next_id], next_id, cur_time + next_time))
        return -1


print(
    Solution().minCost(
        maxTime=29,
        edges=[[0, 1, 10], [1, 2, 10], [2, 5, 10], [0, 3, 1], [3, 4, 10], [4, 5, 15]],
        passingFees=[5, 1, 2, 20, 20, 3],
    )
)
