from typing import List, Tuple
from heapq import heappush, heappop

# 给定一个无向图和打折次数 discounts，
# 求城市 0 到城市 n - 1 不超过打折次数的最小花费，不可达返回 -1。
# Dijkstra 求二维最短路，增加「打折次数」维度，当且仅当打折次数不超过限制且能刷新最短路时，进行松弛。

# 0 <= discounts <= 500
# 2 <= n <= 1000

# 1928. 规定时间内到达终点的最小花费.py
class Solution:
    def minimumCost(self, n: int, highways: List[List[int]], discounts: int) -> int:

        adjList = [[] for _ in range(n)]
        for u, v, w in highways:
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        # 花费，id，打折次数
        pq: List[Tuple[int, int, int]] = [(0, 0, 0)]

        # 经过每个点的打折次数对应的cost visited[点坐标][打折次数]=cost
        dist = [[0x7FFFFFFF] * 505 for _ in range(n)]

        while pq:
            cur_cost, cur_id, cur_discount = heappop(pq)
            if cur_discount > discounts:
                continue

            if cur_id == n - 1:
                return cur_cost

            if cur_cost < dist[cur_id][cur_discount]:
                dist[cur_id][cur_discount] = cur_cost
                for next_id, weight in adjList[cur_id]:
                    heappush(pq, (cur_cost + weight, next_id, cur_discount))
                    heappush(pq, (cur_cost + weight // 2, next_id, cur_discount + 1))
        return -1

    def minimumCost2(self, n: int, highways: List[List[int]], discounts: int) -> int:

        adjList = [[] for _ in range(n)]
        for u, v, w in highways:
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        # 花费，id，打折次数
        pq: List[Tuple[int, int, int]] = [(0, 0, 0)]

        # 每个点的limit
        dist = [0x7FFFFFFF for _ in range(n)]

        while pq:
            cur_cost, cur_id, cur_discount = heappop(pq)
            if cur_discount > discounts:
                continue

            if cur_id == n - 1:
                return cur_cost

            if cur_discount < dist[cur_id]:
                dist[cur_id] = cur_discount
                for next_id, weight in adjList[cur_id]:
                    heappush(pq, (cur_cost + weight, next_id, cur_discount))
                    heappush(pq, (cur_cost + weight // 2, next_id, cur_discount + 1))
        return -1


print(
    Solution().minimumCost(
        n=5, highways=[[0, 1, 4], [2, 1, 3], [1, 4, 11], [3, 2, 3], [3, 4, 2]], discounts=1
    )
)

