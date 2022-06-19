from collections import defaultdict
from typing import List

# 你可以拔开任意两台直连计算机之间的线缆，
# 并用它连接一对未直连的计算机。
# 请你计算并返回使所有计算机都连通所需的最少操作次数。
# 如果不可能，则返回 -1 。
# 1 <= n <= 10^5

# !n台计算机连接成一个网络至少要n-1根线缆
# !如果一个节点数为 m 的连通分量中的边数超过 m - 1，就一定可以找到一条多余的边


class Solution:
    def makeConnected(self, n: int, connections: List[List[int]]) -> int:
        """dfs求连通分量"""
        if len(connections) < n - 1:
            return -1

        def dfs(cur: int) -> None:
            for next in adjMap[cur]:
                if next in visited:
                    continue
                visited.add(next)
                dfs(next)

        adjMap = defaultdict(set)
        for u, v in connections:
            adjMap[u].add(v)
            adjMap[v].add(u)

        count = 0
        visited = set()
        for start in range(n):
            if start not in visited:
                visited.add(start)
                dfs(start)
                count += 1
        return count - 1


print(Solution().makeConnected(n=6, connections=[[0, 1], [0, 2], [0, 3], [1, 2]]))
# 输出：1
# 解释：拔下计算机 1 和 2 之间的线缆，并将它插到计算机 1 和 3 上。
