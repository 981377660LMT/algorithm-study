from typing import List
from collections import defaultdict, deque


class Solution:
    def treeDiameter(self, edges: List[List[int]]) -> int:
        n = len(edges)
        adjMap = defaultdict(list)  # 邻接表
        for x, y in edges:  # 初始化邻接表，建图
            adjMap[x].append(y)
            adjMap[y].append(x)

        queue = deque([0])
        visited = [False for _ in range(n + 1)]
        visited[0] = True
        lastVisited = 0  # 全局变量，好记录第一次BFS最后一个点的ID
        while queue:
            curLen = len(queue)
            for _ in range(curLen):
                lastVisited = queue.popleft()
                for next in adjMap[lastVisited]:
                    if not visited[next]:
                        visited[next] = True
                        queue.append(next)

        visited = [False for _ in range(n + 1)]
        queue = deque([lastVisited])  # 第一次最后一个点，作为第二次BFS的起点
        visited[lastVisited] = True
        level = -1  # 记好距离
        while queue:
            curLen = len(queue)
            for _ in range(curLen):
                cur = queue.popleft()
                for next in adjMap[cur]:
                    if not visited[next]:
                        visited[next] = True
                        queue.append(next)
            level += 1

        return level


print(Solution().treeDiameter(edges=[[0, 1], [1, 2], [2, 3], [1, 4], [4, 5]]))
# 输出：4
# 解释：
# 这棵树上最长的路径是 3 - 2 - 1 - 4 - 5，边数为 4。
