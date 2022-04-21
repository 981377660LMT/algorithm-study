# 5769_经历k个点的最大花费-状压bfs
# bfs适合求路径

# 从每个点出发求某个状态的最短路即可
# 因为是经历k个点,层序扩散即可

from collections import defaultdict, deque
from typing import List


class Solution:
    def maximumCost(self, n: int, highways: List[List[int]], k: int) -> int:
        adjMap = defaultdict(lambda: defaultdict(int))
        for u, v, w in highways:
            adjMap[u][v] = w
            adjMap[v][u] = w

        dist = [[-int(1e20)] * (1 << n) for _ in range(n)]  # 负无穷，要狠一点
        queue = deque([(i, 1 << i, 0, 0) for i in range(n)])  # cur,visited,cost,count

        res = -1
        while queue:
            cur, visited, cost, count = queue.popleft()
            if dist[cur][visited] > cost:  # 非常关键，不加就TLE
                continue

            if count == k:
                res = max(res, cost)
                continue

            for next, weight in adjMap[cur].items():
                if (visited >> next) & 1:
                    continue

                distCand = cost + weight
                if distCand > dist[next][visited | (1 << next)]:
                    dist[next][visited | (1 << next)] = distCand
                    queue.append((next, visited | (1 << next), distCand, count + 1))

        return res


# 只有一条边

print(Solution().maximumCost(2, [[0, 1, 0]], 1))

