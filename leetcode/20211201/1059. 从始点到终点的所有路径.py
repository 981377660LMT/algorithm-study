# 当从始点 source 出发的所有路径都可以到达目标终点 destination 时返回 true，否则返回 false。
# 图中的边数在 0 到 10000 之间。
# 给定的图中可能带有自环和平行边
from typing import List


# 总结
# 0.终点不能有子节点
# 1.dfs终点条件:叶子节点
# 2.回溯visited
# 3.如果任意一条不成立，则return False


class Solution:
    def leadsToDestination(
        self, n: int, edges: List[List[int]], source: int, destination: int
    ) -> bool:
        def dfs(cur: int) -> bool:
            # 没有后继了
            if len(adjList[cur]) == 0:
                return cur == destination
            for next in adjList[cur]:
                if visited[next]:
                    return False
                visited[next] = True
                # 所有路径都要满足
                if not dfs(next):
                    return False
                visited[next] = False
            return True

        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)

        # 终点不能后继结点，必不行
        if len(adjList[destination]) > 0:
            return False

        visited = [False] * n
        visited[source] = True

        return dfs(source)


print(
    Solution().leadsToDestination(
        n=4, edges=[[0, 1], [0, 2], [1, 3], [2, 3]], source=0, destination=3
    )
)
