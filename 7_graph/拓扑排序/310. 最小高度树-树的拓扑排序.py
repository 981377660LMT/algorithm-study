# 请你找到所有的 最小高度树 并按 任意顺序 返回它们的根节点标签列表。
# 思路：不断删除叶子节点，无向图的拓扑排序

# 310. 最小高度树-树的拓扑排序

# 叶子向中心拓扑排序
# !树的拓扑排序

# !叶子结点:度数为0(只有一个结点时)或者1的结点

from collections import deque
from typing import List


class Solution:
    def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
        # if n == 1:
        #     return [0]  # !孤立结点的情形

        adjList = [[] for _ in range(n)]
        deg = [0] * n
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
            deg[u] += 1
            deg[v] += 1

        queue, visited = deque(), [False] * n
        for i in range(n):
            if deg[i] <= 1:  # !叶子结点(包括孤立的结点)
                queue.append(i)
                visited[i] = True

        res = deque()  # 记录最后一层的节点
        while queue:
            res = queue.copy()
            len_ = len(queue)
            for _ in range(len_):  # !每次删除外层的叶子节点,一层一层删除
                cur = queue.popleft()
                for next in adjList[cur]:
                    if not visited[next]:
                        deg[next] -= 1
                        if deg[next] == 1:
                            visited[next] = True
                            queue.append(next)

        return list(res)


print(Solution().findMinHeightTrees(4, [[1, 0], [1, 2], [1, 3]]))
