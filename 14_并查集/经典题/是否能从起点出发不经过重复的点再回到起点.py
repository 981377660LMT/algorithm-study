# !给定Flood Fill中，给定一张 n×m 的网格图，
# 其中有障碍和可行点，其中一个点是起点。
# 问是否能从起点出发不经过重复的点再回到起点。
# !起点是否在大小>=4的环上

# !dfs从源点找环或者并查集
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


DIR4 = [(1, 0), (-1, 0), (0, 1), (0, -1)]


class UnionFindArray:
    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)


def solve2(grid: List[str]) -> bool:
    """并查集"""
    ROW, COL = len(grid), len(grid[0])
    sr, sc = -1, -1
    uf = UnionFindArray(ROW * COL)
    for r in range(ROW):
        for c in range(COL):
            if grid[r][c] == "S":
                sr, sc = r, c
            elif grid[r][c] == ".":
                if r > 0 and grid[r - 1][c] == ".":
                    uf.union(r * COL + c, (r - 1) * COL + c)
                if c > 0 and grid[r][c - 1] == ".":
                    uf.union(r * COL + c, r * COL + c - 1)

    res = set()  # !看周围点的分组 (是否至少有两个在一个组)
    count = 0
    for dr, dc in DIR4:
        nr, nc = sr + dr, sc + dc
        if 0 <= nr < ROW and 0 <= nc < COL:
            res.add(uf.find(nr * COL + nc))
            count += 1
    return len(res) < count


def solve1(grid: List[str]) -> bool:
    """dfs"""
    ROW, COL = len(grid), len(grid[0])
    sr, sc = -1, -1
    for r in range(ROW):
        for c in range(COL):
            if grid[r][c] == "S":
                sr, sc = r, c
                break
        if sr != -1:
            break

    def dfs(cur: int, depth: int) -> None:
        global res
        if visited[cur]:
            if cur == start and depth >= 4:
                res = True
            return
        visited[cur] = True
        for dr, dc in DIR4:
            nr, nc = cur // COL + dr, cur % COL + dc
            if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] != "#":
                dfs(nr * COL + nc, depth + 1)

    res = False
    start = sr * COL + sc
    visited = [False] * (ROW * COL)
    dfs(start, 0)
    return res


if __name__ == "__main__":
    _ROW, _ = map(int, input().split())
    grid = [input() for _ in range(_ROW)]
    print("Yes" if solve1(grid) else "No")
