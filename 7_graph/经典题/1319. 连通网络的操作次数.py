from typing import List

# 你可以拔开任意两台直连计算机之间的线缆，
# 并用它连接一对未直连的计算机。
# 请你计算并返回使所有计算机都连通所需的最少操作次数。
# 如果不可能，则返回 -1 。
# 1 <= n <= 10^5

# n台计算机连接成一个网络至少要n-1根线缆；
# n台计算机连接成一个网络，等用于整个网络只有一个连通分量


class Solution:
    def build_graph(self, n, connections):
        graph = [[] for _ in range(n)]

        for n1, n2 in connections:
            graph[n1].append(n2)
            graph[n2].append(n1)
        return graph

    def dfs(self, graph, visited, node):
        if node in visited:
            return

        visited.add(node)
        for neighbor in graph[node]:
            self.dfs(graph, visited, neighbor)

    def makeConnected(self, n: int, connections: List[List[int]]) -> int:
        if len(connections) < n - 1:
            return -1

        # 建图
        graph = self.build_graph(n, connections)

        visited = set()
        connected_components = 0
        # 遍历顶点
        for node in range(n):
            # 如果没有访问过
            if node not in visited:
                connected_components += 1
                self.dfs(graph, visited, node)

        return connected_components - 1


print(Solution().makeConnected(n=4, connections=[[0, 1], [0, 2], [1, 2]]))
# 输出：1
# 解释：拔下计算机 1 和 2 之间的线缆，并将它插到计算机 1 和 3 上。
