"""竞赛图找环/基环树找环/基环树森林找环"""


from collections import deque
from functools import lru_cache
from typing import Iterable, List, Mapping, Sequence, Tuple, Union


SequenceGraph = Sequence[Iterable[int]]
MappingGraph = Mapping[int, Iterable[int]]
Graph = Union[SequenceGraph, MappingGraph]


def cyclePartition(
    n: int, graph: Graph, directed: bool
) -> Tuple[List[List[int]], List[bool], List[int], List[int]]:
    """返回基环树森林的环分组信息(环的大小>=2)以及每个点在拓扑排序中的最大深度.

    Args:
        - n: 图的节点数.
        - graph: 图的邻接表表示.
        - directed: 图是否有向.

    Returns:
        - groups: 环分组,每个环的大小>=2.
        - inCycle: 每个点是否在环中.
        - belong: 每个点所在的环的编号.如果不在环中,则为-1.
        - depth: 每个点在拓扑排序中的最大深度,最外层的点深度为0.
    """

    def max(a: int, b: int) -> int:
        return a if a > b else b

    deg = [0] * n
    if directed:
        for u in range(n):
            for v in graph[u]:
                deg[v] += 1
    else:
        for u in range(n):
            for v in graph[u]:
                if u < v:
                    deg[u] += 1
                    deg[v] += 1

    startDeg = 0 if directed else 1
    queue = deque([i for i in range(n) if deg[i] == startDeg])
    visited = [False] * n
    depth = [0] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next_ in graph[cur]:
            depth[next_] = max(depth[next_], depth[cur] + 1)
            deg[next_] -= 1
            if deg[next_] == startDeg:
                queue.append(next_)

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in graph[cur]:
            dfs(next, path)

    groups = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        groups.append(path)

    inCycle, belong = [False] * n, [-1] * n
    for gid, group in enumerate(groups):
        for node in group:
            inCycle[node] = True
            belong[node] = gid

    return groups, inCycle, belong, depth


if __name__ == "__main__":
    # 2360. 图中的最长环
    # https://leetcode.cn/problems/longest-cycle-in-a-graph/description/
    # !求内向基环树(每个点出度最多为1)的最大环
    class Solution:
        def longestCycle(self, edges: List[int]) -> int:
            """
            每个节点至多有一条出边
            外向基环树最大环
            """
            n = len(edges)
            adjList = [[] for _ in range(n)]
            for u, v in enumerate(edges):
                if v == -1:
                    continue
                adjList[u].append(v)

            cycle, *_ = cyclePartition(n, adjList, directed=True)
            return max((len(g) for g in cycle), default=-1)

    # 457. 环形数组是否存在循环
    # https://leetcode.cn/problems/circular-array-loop/
    class Solution2:
        def circularArrayLoop(self, nums: List[int]) -> bool:
            def getNext(i: int) -> int:
                return (i + nums[i]) % n

            n = len(nums)
            adjList = [[] for _ in range(n)]
            for i in range(n):
                j = getNext(i)
                if i == j:
                    continue
                if nums[i] * nums[j] > 0:
                    adjList[i].append(j)

            cycles, *_ = cyclePartition(n, adjList, directed=True)
            return any(len(g) > 1 for g in cycles)

    # 2127. 参加会议的最多员工数
    # https://leetcode.cn/problems/maximum-employees-to-be-invited-to-a-meeting/
    class Solution3:
        def maximumInvitations(self, favorite: List[int]) -> int:
            n = len(favorite)
            adjList = [[] for _ in range(n)]
            for u, v in enumerate(favorite):
                adjList[u].append(v)

            cycleGroup, *_, depth = cyclePartition(n, adjList, directed=True)
            # 两种情况:1.所有的二元基环树里的最长链之和;2.唯一的最长环的长度
            cand1 = sum((1 + depth[i]) for i in range(n) if favorite[favorite[i]] == i)
            cand2 = max(len(cycle) for cycle in cycleGroup)
            return max(cand1, cand2)

    # 2204. 无向图中到环的距离
    # https://leetcode.cn/problems/distance-to-a-cycle-in-undirected-graph/
    class Solution4:
        def distanceToCycle(self, n: int, edges: List[List[int]]) -> List[int]:
            """从基环出发，求所有树枝上的点的深度."""
            adjList = [[] for _ in range(n)]
            for u, v in edges:
                adjList[u].append(v)
                adjList[v].append(u)
            cycles, *_ = cyclePartition(n, adjList, directed=False)

            def bfsMultiStart(starts: Iterable[int], adjList: List[List[int]]) -> List[int]:
                """多源bfs"""
                n = len(adjList)
                dist = [int(1e18)] * n
                queue = deque(starts)
                for start in starts:
                    dist[start] = 0
                while queue:
                    cur = queue.popleft()
                    for next in adjList[cur]:
                        cand = dist[cur] + 1
                        if cand < dist[next]:
                            dist[next] = cand
                            queue.append(next)
                return dist

            return bfsMultiStart(cycles[0], adjList)

    # 100075. 有向图访问计数
    # https://leetcode.cn/problems/count-visited-nodes-in-a-directed-graph/description/
    class Solution5:
        def countVisitedNodes(self, edges: List[int]) -> List[int]:
            n = len(edges)
            adjList = [[] for _ in range(n)]
            for u, v in enumerate(edges):
                adjList[u].append(v)
            cycles, inCycle, belong, *_ = cyclePartition(n, adjList, directed=True)

            @lru_cache(None)
            def dfs(cur: int) -> int:
                if inCycle[cur]:
                    return len(cycles[belong[cur]])
                return 1 + dfs(edges[cur])

            return [dfs(i) for i in range(n)]
