# @param {number} n
# @param {number[][]} flights
# @param {number} src
# @param {number} dst
# @param {number} k  最多经过 k 站中转的路线
# @return {number}
# 找到出一条最多经过 k 站中转的路线，使得从 src 到 dst 的 价格最便宜 ，
# 并返回该价格。 如果不存在这样的路线，则输出 -1。
# @summary
# 带限制的最短路径


from typing import List
from heapq import heappop, heappush


# pq超时 需要用bellman ford
class Solution:
    def findCheapestPrice(
        self, n: int, flights: List[List[int]], src: int, dst: int, k: int
    ) -> int:
        adjList = [[] for _ in range(n)]
        for u, v, w in flights:

            adjList[u].append((v, w))

        pq = [(0, src, k + 1)]
        while pq:
            cur_cost, cur_id, cur_k = heappop(pq)
            if cur_id == dst:
                return cur_cost
            if cur_k > 0:
                for next_id, next_cost in adjList[cur_id]:
                    heappush(pq, (cur_cost + next_cost, next_id, cur_k - 1))

        return -1


print(
    Solution().findCheapestPrice(
        n=3, flights=[[0, 1, 100], [1, 2, 100], [0, 2, 500]], src=0, dst=2, k=1
    )
)
