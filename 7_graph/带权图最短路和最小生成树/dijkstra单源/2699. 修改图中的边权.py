# https://www.luogu.com.cn/problem/solution/CF715B
# https://leetcode.cn/problems/modify-graph-edge-weights/
# 给你一个 n 个节点的 无向带权连通 图，节点编号为 0 到 n - 1 ，
# 再给你一个整数数组 edges ，其中 edges[i] = [ai, bi, wi] 表示节点 ai 和 bi 之间有一条边权为 wi 的边。
# 部分边的边权为 -1（wi = -1），其他边的边权都为 正 数（wi > 0）。
# !你需要将所有边权为 -1 的边都修改为范围 [1, 2e9] 中的 正整数 ，
# !使得从节点 source 到节点 destination 的 最短距离 为整数 target 。
# 如果有 多种 修改方案可以使 source 和 destination 之间的最短距离等于 target ，你可以返回任意一种方案。
# 如果存在使 source 到 destination 最短距离为 target 的方案，请你按任意顺序返回包含所有边的数组（包括未修改边权的边）。如果不存在这样的方案，请你返回一个 空数组 。
# 注意：你不能修改一开始边权为正数的边。

# 1 <= n <= 1e3
# 1 <= edges.length <= 1e4



# 解:
# !1. 第一次修改, 找到最短路, 如果最短路大于target, 则无解，然后把路径上的边权修改.
# !2. 修改完后, while循环里检查是否有最短路小于target, 如果有, 则把路径上的边权修改.

from collections import defaultdict
from typing import List
from dijkstra模板 import dijkstra2

NEED_MODIFY = -1  # 需要修改的边权
LOWER = 1
UPPER = int(2e9)
INF = int(4e18)


class Solution:
    def modifiedGraphEdges(
        self, n: int, edges: List[List[int]], source: int, destination: int, target: int
    ) -> List[List[int]]:
        adjList = [[] for _ in range(n)]
        eid = defaultdict(lambda: defaultdict(int))
        todo = [False] * len(edges)
        for i, (u, v, w) in enumerate(edges):
            if w == NEED_MODIFY:
                w = LOWER  # 需要修改的边权置为最小
                todo[i] = True
            adjList[u].append((v, w))
            adjList[v].append((u, w))
            eid[u][v] = i
            eid[v][u] = i

        dist, path = dijkstra2(n, adjList, source, destination)
        if dist > target:
            return []

        eids = [eid[u][v] for u, v in zip(path, path[1:])][::-1]
        diff = target - dist
        for e in eids:
            if not todo[e]:
                continue
            canChange = min(diff, UPPER - LOWER)
            edges[e][2] = LOWER + canChange
            diff -= canChange
        if diff != 0:
            return []
        for i in range(len(edges)):
            if edges[i][2] == NEED_MODIFY:
                edges[i][2] = UPPER

        while True:
            # check
            newAdjList = [[] for _ in range(n)]
            for u, v, w in edges:
                newAdjList[u].append((v, w))
                newAdjList[v].append((u, w))
            dist, path = dijkstra2(n, newAdjList, source, destination)
            eids = [eid[u][v] for u, v in zip(path, path[1:])]
            if all(not todo[e] for e in eids) and dist < target:
                return []

            if dist >= target:
                break
            diff = target - dist
            for e in eids:
                if not todo[e]:
                    continue
                canChange = min(diff, UPPER - edges[e][2])
                edges[e][2] = edges[e][2] + canChange
                diff -= canChange
                if diff == 0:
                    break

        return edges


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, m, target, source, destination = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        edges.append([u, v, w])
    res = Solution().modifiedGraphEdges(n, edges, source, destination, target)
    if not res:
        print("NO")
    else:
        print("YES")
        for u, v, w in res:
            print(u, v, w)
