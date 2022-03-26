from typing import List

INF = int(1e20)


class Solution:
    def findCheapestPrice(
        self, n: int, flights: List[List[int]], src: int, dst: int, k: int
    ) -> int:
        dist = [INF] * (n)
        dist[src] = 0

        for _ in range(k + 1):
            preDist = dist[:]
            for u, v, w in flights:
                if preDist[u] + w < dist[v]:
                    dist[v] = preDist[u] + w

        return -1 if dist[dst] == INF else dist[dst]


print(
    Solution().findCheapestPrice(
        n=4,
        flights=[[0, 1, 100], [1, 2, 100], [2, 0, 100], [1, 3, 600], [2, 3, 200]],
        src=0,
        dst=3,
        k=1,
    )
)

print(Solution().findCheapestPrice(4, [[0, 1, 1], [0, 2, 5], [1, 2, 1], [2, 3, 1]], 0, 3, 1,))
print(
    Solution().findCheapestPrice(
        10,
        [
            [3, 4, 4],
            [2, 5, 6],
            [4, 7, 10],
            [9, 6, 5],
            [7, 4, 4],
            [6, 2, 10],
            [6, 8, 6],
            [7, 9, 4],
            [1, 5, 4],
            [1, 0, 4],
            [9, 7, 3],
            [7, 0, 5],
            [6, 5, 8],
            [1, 7, 6],
            [4, 0, 9],
            [5, 9, 1],
            [8, 7, 3],
            [1, 2, 6],
            [4, 1, 5],
            [5, 2, 4],
            [1, 9, 1],
            [7, 8, 10],
            [0, 4, 2],
            [7, 2, 8],
        ],
        6,
        0,
        7,
    )
)
print(
    Solution().findCheapestPrice(
        11,
        [
            [0, 3, 3],
            [3, 4, 3],
            [4, 1, 3],
            [0, 5, 1],
            [5, 1, 100],
            [0, 6, 2],
            [6, 1, 100],
            [0, 7, 1],
            [7, 8, 1],
            [8, 9, 1],
            [9, 1, 1],
            [1, 10, 1],
            [10, 2, 1],
            [1, 2, 100],
        ],
        0,
        2,
        4,
    )
)
