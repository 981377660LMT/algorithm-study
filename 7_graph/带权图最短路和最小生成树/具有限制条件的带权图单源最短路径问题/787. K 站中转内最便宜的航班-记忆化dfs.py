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
from functools import lru_cache

INF = int(1e20)


class Solution:
    def findCheapestPrice(
        self, n: int, flights: List[List[int]], src: int, dst: int, k: int
    ) -> int:
        @lru_cache(None)
        def dfs(city: int, remain: int) -> int:
            if remain < 0:
                return INF
            if city == dst:
                return 0

            res = INF
            for next, weight in adjList[city]:
                res = min(res, dfs(next, remain - 1) + weight)
            return res

        adjList = [[] for _ in range(n)]
        for u, v, w in flights:
            adjList[u].append((v, w))
        res = dfs(src, k + 1)
        dfs.cache_clear()
        return res if res != INF else -1


print(
    Solution().findCheapestPrice(
        n=3, flights=[[0, 1, 100], [1, 2, 100], [0, 2, 500]], src=0, dst=2, k=1
    )
)
