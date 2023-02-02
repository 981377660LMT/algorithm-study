# 1129. 颜色交替的最短路径-带有状态的二维dist

# 在一个有向图中，节点分别标记为 0, 1, ..., n-1。
# 图中每条边为红色或者蓝色，且存在自环或平行边。
# red_edges 中的每一个 [i, j] 对表示从节点 i 到节点 j 的红色有向边。
# 类似地，blue_edges 中的每一个 [i, j] 对表示从节点 i 到节点 j 的蓝色有向边。
# 返回长度为 n 的数组 answer，
# 其中 answer[X] 是从节点 0 到节点 X 的红色边和蓝色边交替出现的最短路径的长度。
# 如果不存在这样的路径，那么 answer[x] = -1。

# 特殊建图+bfs求无权图最短路

from collections import deque
from typing import List

INF = int(1e18)


class Solution:
    def shortestAlternatingPaths(
        self, n: int, redEdges: List[List[int]], blueEdges: List[List[int]]
    ) -> List[int]:
        adjList = [[[], []] for _ in range(n)]  # 0: red, 1: blue
        for u, v in redEdges:
            adjList[u][0].append(v)
        for u, v in blueEdges:
            adjList[u][1].append(v)

        queue = deque([(0, 0, 0), (0, 1, 0)])  # (node, color, dist) 起点为红蓝两种
        dist = [[INF] * 2 for _ in range(n)]
        dist[0][0] = 0
        dist[0][1] = 0
        while queue:
            cur, curColor, curDist = queue.popleft()
            if curDist > dist[cur][curColor]:
                continue
            for next in adjList[cur][curColor]:
                cand, nextColor = curDist + 1, curColor ^ 1
                if cand < dist[next][nextColor]:
                    dist[next][nextColor] = cand
                    queue.append((next, nextColor, cand))  # type: ignore

        res = [-1] * n
        for i, (a, b) in enumerate(dist):
            if a == b == INF:
                continue
            res[i] = a if a < b else b
        return res


print(
    Solution().shortestAlternatingPaths(
        5, [[0, 1], [1, 2], [2, 3], [3, 4]], [[1, 2], [2, 3], [3, 1]]
    )
)
