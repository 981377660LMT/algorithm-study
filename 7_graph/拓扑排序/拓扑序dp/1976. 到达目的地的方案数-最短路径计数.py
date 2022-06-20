from typing import List
from collections import defaultdict
from heapq import heappop, heappush

# 你想知道花费 最少时间 从路口 0 出发到达路口 n - 1 的方案数。
# 请返回花费 最少时间 到达目的地的 路径数目 。由于答案可能很大，将结果对 109 + 7 取余 后返回。
# 1 <= n <= 200

# 总结：
# 1.If we meet candidate == dist[neib], it means we found one more way to reach node with minimal cost.
# 2.If candidate < dist[neib], it means that we found better candidate, so we update distance and put cnt[neib] = cnt[idx].

MOD = int(1e9 + 7)
INF = int(1e20)

# 单源最短路 +  DP


class Solution:
    def countPaths(self, n: int, roads: List[List[int]]) -> int:
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in roads:
            adjMap[u][v] = w
            adjMap[v][u] = w

        dist = defaultdict(lambda: INF)
        dist[0] = 0
        pq = [(0, 0)]

        # 1.注意这个count数组表示到id的最短路径数
        count = defaultdict(int)
        count[0] = 1

        while pq:
            curDist, cur = heappop(pq)
            if cur == n - 1:
                return count[cur]

            for next in adjMap[cur]:
                cand = adjMap[cur][next] + curDist

                # 2.相等加count
                if cand == dist[next]:
                    count[next] += count[cur]
                    count[next] %= MOD
                # 3.更优直接覆盖count
                elif cand < dist[next]:
                    dist[next] = cand
                    heappush(pq, (cand, next))
                    count[next] = count[cur]
                    count[next] %= MOD
        return -1


print(
    Solution().countPaths(
        n=7,
        roads=[
            [0, 6, 7],
            [0, 1, 2],
            [1, 2, 3],
            [1, 3, 3],
            [6, 3, 3],
            [3, 5, 1],
            [6, 5, 1],
            [2, 5, 1],
            [0, 4, 5],
            [4, 6, 2],
        ],
    )
)

# 输出：4
# 解释：从路口 0 出发到路口 6 花费的最少时间是 7 分钟。
# 四条花费 7 分钟的路径分别为：
# - 0 ➝ 6
# - 0 ➝ 4 ➝ 6
# - 0 ➝ 1 ➝ 2 ➝ 5 ➝ 6
# - 0 ➝ 1 ➝ 3 ➝ 5 ➝ 6

