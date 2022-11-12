from typing import List
from collections import deque


# !多用数组来存图减小开销
# !注意到行和列是正交的 因此可以分别处理拓扑序
# !二维拓扑排序
class Solution:
    def buildMatrix(
        self, k: int, rowConditions: List[List[int]], colConditions: List[List[int]]
    ) -> List[List[int]]:
        def topoSort(adjList: List[List[int]], deg: List[int]) -> List[int]:
            res = []
            queue = deque([i for i in range(1, k + 1) if deg[i] == 0])
            while queue:
                cur = queue.popleft()
                res.append(cur)
                for next in adjList[cur]:
                    deg[next] -= 1
                    if deg[next] == 0:
                        queue.append(next)
            if any(deg[i] != 0 for i in range(1, k + 1)):
                return []
            return res

        adjList1 = [[] for _ in range((k + 1))]
        deg1 = [0] * (k + 1)
        for pre, cur in rowConditions:
            adjList1[pre].append(cur)
            deg1[cur] += 1

        adjList2 = [[] for _ in range((k + 1))]
        deg2 = [0] * (k + 1)
        for pre, cur in colConditions:
            adjList2[pre].append(cur)
            deg2[cur] += 1

        rowOrder = topoSort(adjList1, deg1)
        if not rowOrder:
            return []

        colOrder = topoSort(adjList2, deg2)
        if not colOrder:
            return []

        mp1 = {num: r for r, num in enumerate(rowOrder)}
        mp2 = {num: c for c, num in enumerate(colOrder)}
        res = [[0] * k for _ in range(k)]
        for num in range(1, k + 1):
            res[mp1[num]][mp2[num]] = num
        return res


print(Solution().buildMatrix(k=3, rowConditions=[[1, 2], [3, 2]], colConditions=[[2, 1], [3, 2]]))
print(
    Solution().buildMatrix(
        k=3, rowConditions=[[1, 2], [2, 3], [3, 1], [2, 3]], colConditions=[[2, 1]]
    )
)
