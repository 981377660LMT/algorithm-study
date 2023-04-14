"""基环树找环/基环树森林找环"""

from collections import deque
from typing import List, Tuple


def findCycleAndCalDepth(
    n: int, adjList: List[List[int]], deg: List[int], *, isDirected: bool
) -> Tuple[List[List[int]], List[int]]:
    """无/有向基环树森林找环上的点,并记录每个点在拓扑排序中的最大深度,最外层的点深度为0"""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    depth = [0] * n
    startDeg = 0 if isDirected else 1
    queue = deque([i for i in range(n) if deg[i] == startDeg])
    visited = [False] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next in adjList[cur]:
            depth[next] = max(depth[next], depth[cur] + 1)
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in adjList[cur]:
            dfs(next, path)

    cycleGroup = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        cycleGroup.append(path)

    return cycleGroup, depth


if __name__ == "__main__":
    # 2360. 图中的最长环
    # !求内向基环树(每个点出度最多为1)的最大环
    class Solution:
        def longestCycle(self, edges: List[int]) -> int:
            """
            每个节点至多有一条出边
            外向基环树最大环
            """
            n = len(edges)
            adjList = [[] for _ in range(n)]
            deg = [0] * n
            for u, v in enumerate(edges):
                if v == -1:
                    continue
                adjList[u].append(v)
                deg[v] += 1

            cycle, _ = findCycleAndCalDepth(n, adjList, deg, isDirected=True)
            return max((len(g) for g in cycle), default=-1)

    # https://leetcode.cn/problems/circular-array-loop/
    # 457. 环形数组是否存在循环
    class Solution2:
        def circularArrayLoop(self, nums: List[int]) -> bool:
            def getNext(i: int) -> int:
                return (i + nums[i]) % n

            n = len(nums)
            adjList = [[] for _ in range(n)]
            deg = [0] * n
            for i in range(n):
                j = getNext(i)
                if i == j:
                    continue
                if nums[i] * nums[j] > 0:
                    adjList[i].append(j)
                    deg[j] += 1
            cycles, _ = findCycleAndCalDepth(n, adjList, deg, isDirected=True)
            return any(len(g) > 1 for g in cycles)
