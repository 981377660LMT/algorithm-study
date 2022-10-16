from collections import defaultdict, deque
from decimal import MAX_EMAX


class Point:
    def __init__(self, a=0, b=0):
        self.x = a
        self.y = b


# 给出一张带权无向图，图中任意两点间有且仅有一条路径。
# 计算从任意点出发并访问完所有节点经过边的权值之和的最小值。
# n<=100000

# 最好weight之和，最坏两倍weight之和
# !最大化只有一遍的路径权重之和=>找`树的一端到另一端最长的路径`
class Solution:
    def solve(self, n, Edge, val):
        def bfs1(start: int) -> int:
            """从任一点出发找到直径的一端"""
            queue = deque([(start, 0)])
            visited = set([start])
            res, maxDist = start, 0
            while queue:
                cur, dist = queue.popleft()
                if dist > maxDist:
                    maxDist = dist
                    res = cur
                for next in adjMap[cur]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append((next, dist + adjMap[cur][next]))
            return res

        def bfs2(start: int) -> int:
            """找到直径长度"""
            queue = deque([(start, 0)])
            visited = set([start])
            res = 0
            while queue:
                cur, dist = queue.popleft()
                res = max(res, dist)
                for next in adjMap[cur]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append((next, dist + adjMap[cur][next]))
            return res

        # write code here
        res = 0
        adjMap = defaultdict(lambda: defaultdict(int))
        for i in range(n - 1):
            u, v, w = Edge[i].x, Edge[i].y, val[i]
            adjMap[u][v] = w
            adjMap[v][u] = w
            res += 2 * w
        start = bfs1(1)
        diameter = bfs2(start)
        return res - diameter
