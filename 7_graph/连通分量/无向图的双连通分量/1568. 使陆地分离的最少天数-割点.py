# 一共3种情况，0，1，2， 并查集求岛屿数量如果大于2 返回0， 如果岛屿数量为1， tarjan算法求割点，
# 如果找到割点返回 1，没有割点则返回，2


from copy import deepcopy
from typing import List


def findCutVertices(n: int, graph: List[List[int]]) -> List[bool]:
    """Tarjan 算法求无向图的割点

    Args:
        n (int): 顶点数
        graph (List[List[int]]): 邻接表

    Returns:
        List[bool]: 每个点是否是割点
    """

    def dfs(cur: int, pre: int) -> int:
        nonlocal dfsId
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        childCount = 0
        for next in graph[cur]:
            if dfsOrder[next] == 0:
                childCount += 1
                nextLow = dfs(next, cur)
                if nextLow >= dfsOrder[cur]:
                    isCut[cur] = True
                if nextLow < curLow:
                    curLow = nextLow
            elif next != pre and dfsOrder[next] < curLow:
                curLow = dfsOrder[next]
        if pre == -1 and childCount == 1:  # 特判：只有一个儿子的树根，删除后并没有增加连通分量的个数，这种情况下不是割顶
            isCut[cur] = False
        return curLow

    isCut = [False] * n
    dfsOrder = [0] * n  # 值从 1 开始
    dfsId = 0
    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i, -1)

    return isCut


class Solution:
    def minDays(self, grid: List[List[int]]) -> int:
        # 特判
        oneCount = sum(row.count(1) for row in grid)
        if oneCount == 0:
            return 0
        elif oneCount == 1:
            return 1

        # 连通分量不为1的情况
        gridCopy = deepcopy(grid)
        part = self.floodFill(gridCopy)
        if part != 1:
            return 0

        # tarjan寻找割点
        ROW, COL = len(grid), len(grid[0])
        adjList = [[] for _ in range(ROW * COL)]
        for i in range(ROW):
            for j in range(COL):
                if grid[i][j] == 1:
                    cur = i * COL + j
                    if i - 1 >= 0 and grid[i - 1][j] == 1:
                        next = (i - 1) * COL + j
                        adjList[cur].append(next)
                        adjList[next].append(cur)
                    if j + 1 < COL and grid[i][j + 1] == 1:
                        next = i * COL + j + 1
                        adjList[cur].append(next)
                        adjList[next].append(cur)

        isCut = findCutVertices(ROW * COL, adjList)
        if any(isCut):
            return 1
        else:
            return 2

    def floodFill(self, grid: List[List[int]]) -> int:
        def dfs(r: int, c: int) -> None:
            if grid[r][c] == 0:
                return
            grid[r][c] = 0
            for dr, dc in [(1, 0), (-1, 0), (0, 1), (0, -1)]:
                nr, nc = r + dr, c + dc
                if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] == 1:
                    dfs(nr, nc)

        res = 0
        row, col = len(grid), len(grid[0])
        for r in range(row):
            for c in range(col):
                if grid[r][c] == 1:
                    res += 1
                    dfs(r, c)
        return res


assert Solution().minDays([[0, 1, 1, 0], [0, 1, 1, 0], [0, 0, 0, 0]]) == 2
assert (
    Solution().minDays(grid=[[0, 1, 0, 1, 1], [1, 1, 1, 1, 1], [1, 1, 1, 1, 1], [1, 1, 1, 1, 0]])
    == 1
)
