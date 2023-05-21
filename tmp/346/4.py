from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 n 个节点的 无向带权连通 图，节点编号为 0 到 n - 1 ，再给你一个整数数组 edges ，其中 edges[i] = [ai, bi, wi] 表示节点 ai 和 bi 之间有一条边权为 wi 的边。

# 部分边的边权为 -1（wi = -1），其他边的边权都为 正 数（wi > 0）。

# 你需要将所有边权为 -1 的边都修改为范围 [1, 2 * 109] 中的 正整数 ，使得从节点 source 到节点 destination 的 最短距离 为整数 target 。如果有 多种 修改方案可以使 source 和 destination 之间的最短距离等于 target ，你可以返回任意一种方案。


# 如果存在使 source 到 destination 最短距离为 target 的方案，请你按任意顺序返回包含所有边的数组（包括未修改边权的边）。如果不存在这样的方案，请你返回一个 空数组 。
# 1 <= n <= 100
# 1 <= edges.length <= n * (n - 1) / 2
# edges[i].length == 3
# 0 <= ai, bi < n
# wi = -1 或者 1 <= wi <= 107
# ai != bi
# 0 <= source, destination < n
# source != destination
# 1 <= target <= 109
# 输入的图是连通图，且没有自环和重边。


# -1当作INF

from typing import List, Sequence, Tuple
from heapq import heappop, heappush


def dijkstra2(
    n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int, end: int
) -> Tuple[int, List[int]]:
    """dijkstra求出起点到end的(最短距离,路径) 时间复杂度O((V+E)logV)"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    pre = [-1] * n  # 记录每个点的前驱

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
                pre[next] = cur

    path = []
    cur = end
    while pre[cur] != -1:
        path.append(cur)
        cur = pre[cur]
    path.append(start)
    return dist[end], path[::-1]


class Solution:
    def modifiedGraphEdges(
        self, n: int, edges: List[List[int]], source: int, destination: int, target: int
    ) -> List[List[int]]:
        adjList = [[] for _ in range(n)]
        eid = defaultdict(lambda: defaultdict(int))
        isBad = set()
        for i, (u, v, w) in enumerate(edges):
            if w == -1:
                w = 1
                isBad.add(i)
            adjList[u].append((v, w))
            adjList[v].append((u, w))
            eid[u][v] = i
            eid[v][u] = i

        dist, path = dijkstra2(n, adjList, source, destination)
        if dist > target:
            return []

        # print(path)
        eids = [eid[u][v] for u, v in zip(path, path[1:])][::-1]
        diff = target - dist
        for e in eids:
            if e not in isBad:
                continue
            canChange = min(diff, int(2e9) - 1)
            edges[e][2] = 1 + canChange
            diff -= canChange
        if diff != 0:
            return []
        for i in range(len(edges)):
            if edges[i][2] == -1:
                edges[i][2] = int(2e9)

        while True:
            # check
            newAdjList = [[] for _ in range(n)]
            for u, v, w in edges:
                newAdjList[u].append((v, w))
                newAdjList[v].append((u, w))
            dist, path = dijkstra2(n, newAdjList, source, destination)
            eids = [eid[u][v] for u, v in zip(path, path[1:])]
            if all(e not in isBad for e in eids) and dist < target:
                return []

            if dist >= target:
                break
            diff = target - dist
            for e in eids:
                if e not in isBad:
                    continue
                canChange = min(diff, int(2e9) - edges[e][2])
                edges[e][2] = edges[e][2] + canChange
                diff -= canChange
                if diff == 0:
                    break
        # newAdjList = [[] for _ in range(n)]
        # for u, v, w in edges:
        #     newAdjList[u].append((v, w))
        #     newAdjList[v].append((u, w))
        # dist, path = dijkstra2(n, newAdjList, source, destination)
        # print(dist)
        return edges


# # n = 5, edges = [[4,1,-1],[2,0,-1],[0,3,-1],[4,3,-1]], source = 0, destination = 1, target = 5
# print(
#     Solution().modifiedGraphEdges(
#         n=5,
#         edges=[[4, 1, -1], [2, 0, -1], [0, 3, -1], [4, 3, -1]],
#         source=0,
#         destination=1,
#         target=5,
#     )
# )
# # n = 4, edges = [[1,0,4],[1,2,3],[2,3,5],[0,3,-1]], source = 0, destination = 2, target = 6
# print(
#     Solution().modifiedGraphEdges(
#         n=4,
#         edges=[[1, 0, 4], [1, 2, 3], [2, 3, 5], [0, 3, -1]],
#         source=0,
#         destination=2,
#         target=6,
#     )
# )
# 5
# [[1,4,1],[2,4,-1],[3,0,2],[0,4,-1],[1,3,10],[1,0,10]]
# 0
# 2
# 15
print(
    Solution().modifiedGraphEdges(
        n=5,
        edges=[[1, 4, 1], [2, 4, -1], [3, 0, 2], [0, 4, -1], [1, 3, 10], [1, 0, 10]],
        source=0,
        destination=2,
        target=15,
    )
)
# [[1,4,1],[2,4,4],[3,0,2],[0,4,14],[1,3,10],[1,0,10]]
# 4
# [[0,1,-1],[2,0,2],[3,2,6],[2,1,10],[3,0,-1]]
# 1
# 3
# 12
print(
    Solution().modifiedGraphEdges(
        n=4,
        edges=[[0, 1, -1], [2, 0, 2], [3, 2, 6], [2, 1, 10], [3, 0, -1]],
        source=1,
        destination=3,
        target=12,
    )
)
# 4
# [[0,1,-1],[1,2,-1],[3,1,-1],[3,0,2],[0,2,5]]
# 2
# 3
# 8
print(
    Solution().modifiedGraphEdges(
        n=4,
        edges=[[0, 1, -1], [1, 2, -1], [3, 1, -1], [3, 0, 2], [0, 2, 5]],
        source=2,
        destination=3,
        target=8,
    )
)
