# https://atcoder.jp/contests/abc334/tasks/abc334_g
# G - Christmas Color Grid 2-网格图并查集删点
# 给定一个01矩阵.
# !对每个1，将其变为0后(删点)，求图中1组成的联通分量个数.

from typing import List, Tuple
from UnionFind import UnionFindArray


def christmasColorGrid2(grid: List[List[int]]) -> List[int]:
    ROW, COL = len(grid), len(grid[0])
    adjList = [[] for _ in range(ROW * COL)]
    uf = UnionFindArray(ROW * COL)
    deg = [0] * (ROW * COL)
    onesCount = 0
    for x in range(ROW):
        for y in range(COL):
            if grid[x][y] == 0:
                continue
            cur = x * COL + y
            if x > 0 and grid[x - 1][y] == 1:
                next_ = (x - 1) * COL + y
                adjList[cur].append(next_)
                adjList[next_].append(cur)
                uf.union(cur, next_)
                deg[cur] += 1
                deg[next_] += 1
            if y > 0 and grid[x][y - 1] == 1:
                next_ = x * COL + y - 1
                adjList[cur].append(next_)
                adjList[next_].append(cur)
                uf.union(cur, next_)
                deg[cur] += 1
                deg[next_] += 1
            onesCount += 1
    zeroCount = ROW * COL - onesCount

    basePart = uf.part - zeroCount
    _, belong, isCut = findVBCC(adjList)
    res = []
    for x in range(ROW):
        for y in range(COL):
            if grid[x][y] == 0:
                continue
            cur = x * COL + y
            if isCut[cur]:  # !割点
                res.append(basePart + len(belong[cur]) - 1)
            elif deg[cur] == 0:  # !孤立点
                res.append(basePart - 1)
            else:  # !联通
                res.append(basePart)

    return res


def findVBCC(graph: List[List[int]]) -> Tuple[List[List[int]], List[List[int]], List[bool]]:
    """
    !Tarjan 算法求无向图的 v-BCC

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        Tuple[List[List[int]], List[List[int], List[bool]]:
        每个 v-BCC 组里包含哪些点，每个点所在 v-BCC 的编号(从0开始,割点有多个)，每个顶点是否为割点(便于缩点成树)

    Notes:
        - 原图的割点`至少`在两个不同的 v-BCC 中
        - 原图不是割点的点都`只存在`于一个 v-BCC 中
        - v-BCC 形成的子图内没有割点
    """

    def dfs(cur: int, pre: int) -> int:
        nonlocal dfsId, idCount
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        childCount = 0
        for _, next in enumerate(graph[cur]):
            edge = (cur, next)
            if dfsOrder[next] == 0:
                stack.append(edge)
                childCount += 1
                nextLow = dfs(next, cur)
                if nextLow >= dfsOrder[cur]:
                    isCut[cur] = True
                    idCount += 1
                    group = []
                    while True:
                        topEdge = stack.pop()
                        v1, v2 = topEdge[0], topEdge[1]
                        if vbccId[v1] != idCount:
                            vbccId[v1] = idCount
                            group.append(v1)
                        if vbccId[v2] != idCount:
                            vbccId[v2] = idCount
                            group.append(v2)
                        if v1 == cur and v2 == next:
                            break
                    groups.append(group)
                if nextLow < curLow:
                    curLow = nextLow
            elif next != pre and dfsOrder[next] < dfsOrder[cur]:
                stack.append(edge)
                if dfsOrder[next] < curLow:
                    curLow = dfsOrder[next]
        if pre == -1 and childCount == 1:
            isCut[cur] = False
        return curLow

    n = len(graph)
    dfsId = 0
    dfsOrder = [0] * n
    vbccId = [0] * n
    idCount = 0
    isCut = [False] * n
    stack = []  # (u, v, eid)
    groups = []

    for i, order in enumerate(dfsOrder):
        if order == 0:
            if len(graph[i]) == 0:  # 零度，即孤立点（isolated vertex）
                idCount += 1
                vbccId[i] = idCount
                groups.append([i])
                continue
            dfs(i, -1)

    belong = [[] for _ in range(n)]
    for i, group in enumerate(groups):
        for v in group:
            belong[v].append(i)
    return groups, belong, isCut


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    MOD = 998244353

    h, w = map(int, input().split())
    grid = []
    onesCount = 0
    for _ in range(h):
        s = input()
        onesCount += s.count("#")
        grid.append([0 if c == "." else 1 for c in s])

    res = christmasColorGrid2(grid)
    print(sum(res) * pow(onesCount, MOD - 2, MOD) % MOD)
