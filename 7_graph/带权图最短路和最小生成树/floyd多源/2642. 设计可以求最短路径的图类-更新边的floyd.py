# https://leetcode.cn/problems/design-graph-with-shortest-path-calculator/solution/geng-da-shu-ju-fan-wei-diao-yong-shortes-ibmj/

# 请你实现一个 Graph 类：

# !Graph(int n, int[][] edges) 初始化图有 n 个节点，并输入初始边。

# !addEdge(int[] edge) 向边集中添加一条边，其中 edge = [from, to, edgeCost] 。
# 数据保证添加这条边之前对应的两个节点之间没有有向边。

# !int shortestPath(int node1, int node2) 返回从节点 node1 到 node2 的路径 最小 代价。
# 如果路径不存在，返回 -1 。一条路径的代价是路径中所有边代价之和。

# 1 <= n <= 100
# 0 <= edges.length <= n * (n - 1)
# !调用 shortestPath 1e5 次


from typing import List, Tuple

INF = int(1e18)


class Graph:
    __slots__ = "_dist"

    def __init__(self, n: int, edges: List[Tuple[int, int, int]]):
        dist = [[INF] * n for _ in range(n)]
        for i in range(n):
            dist[i][i] = 0
        for u, v, w in edges:
            dist[u][v] = w

        for k in range(n):
            for i in range(n):
                for j in range(n):
                    cand = dist[i][k] + dist[k][j]
                    if dist[i][j] > cand:
                        dist[i][j] = cand
        self._dist = dist

    def addEdge(self, edge: Tuple[int, int, int]) -> None:
        """
        向边集中添加一条边,保证添加这条边之前对应的两个节点之间没有有向边.
        加边时,枚举每个点对,根据是否经过edge来更新最短路.
        """
        u, v, w = edge
        n = len(self._dist)
        for i in range(n):
            for j in range(n):
                cand = self._dist[i][u] + w + self._dist[v][j]
                if self._dist[i][j] > cand:
                    self._dist[i][j] = cand

    def shortestPath(self, start: int, target: int) -> int:
        """返回从节点 node1 到 node2 的最短路.如果路径不存在，返回 -1."""
        return self._dist[start][target] if self._dist[start][target] < INF else -1


# Your Graph object will be instantiated and called as such:
# obj = Graph(n, edges)
# obj.addEdge(edge)
# param_2 = obj.shortestPath(node1,node2)
