# !多用数组来存图减小开销
# !注意到行和列是正交的 因此可以分别处理拓扑序
# !二维拓扑排序

from typing import List

from topoSort import topoSort


class Solution:
    def buildMatrix(
        self, k: int, rowConditions: List[List[int]], colConditions: List[List[int]]
    ) -> List[List[int]]:
        adjList1 = [[] for _ in range((k))]
        for pre, cur in rowConditions:
            pre, cur = pre - 1, cur - 1
            adjList1[pre].append(cur)
        rowOrder, ok = topoSort(k, adjList1)
        if not ok:
            return []

        adjList2 = [[] for _ in range((k))]
        for pre, cur in colConditions:
            pre, cur = pre - 1, cur - 1
            adjList2[pre].append(cur)
        colOrder, ok = topoSort(k, adjList2)
        if not ok:
            return []

        mp1 = {num: r for r, num in enumerate(rowOrder)}
        mp2 = {num: c for c, num in enumerate(colOrder)}
        res = [[0] * k for _ in range(k)]
        for num in range(k):
            res[mp1[num]][mp2[num]] = num + 1
        return res


print(Solution().buildMatrix(k=3, rowConditions=[[1, 2], [3, 2]], colConditions=[[2, 1], [3, 2]]))
print(
    Solution().buildMatrix(
        k=3, rowConditions=[[1, 2], [2, 3], [3, 1], [2, 3]], colConditions=[[2, 1]]
    )
)
