# G - Big Banned Grid(左上角到右下角是否可达)
# https://atcoder.jp/contests/abc413/tasks/abc413_g
# 在 H×W 的整格网格中，有 K 个被障碍物占据的格子。
# 可在四方向移动至相邻的非障碍物格子。判断从 (1,1) 能否到达 (H,W)。
#
# !最大流->最小割->平面图最小割等于对偶图最短路


from typing import Callable, Optional


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(
        self, key1: int, key2: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        if beforeUnion is not None:
            beforeUnion(root1, root2)
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def unionTo(self, parent: int, child: int) -> bool:
        """定向合并, 将child合并到parent所在的连通分量中."""
        root1, root2 = self.find(parent), self.find(child)
        if root1 == root2:
            return False
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        root = key
        while self._data[root] >= 0:
            root = self._data[root]
        while key != root:
            parent = self._data[key]
            self._data[key] = root
            key = parent
        return root

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


DIR8 = ((0, 1), (1, 0), (0, -1), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1))
if __name__ == "__main__":
    H, W, K = map(int, input().split())
    obstacles = []
    for _ in range(K):
        r, c = map(int, input().split())
        obstacles.append((r - 1, c - 1))

    mp = {(r, c): i for i, (r, c) in enumerate(obstacles)}
    LEFT_BOTTOM = K
    RIGHT_TOP = K + 1

    uf = UnionFindArraySimple(K + 2)
    for r, c in obstacles:
        cur = mp[(r, c)]
        for dr, dc in DIR8:
            nr, nc = r + dr, c + dc
            if 0 <= nr < H and 0 <= nc < W:
                next = mp.get((nr, nc), None)
                if next is not None:
                    uf.union(cur, next)
            else:
                if nr < 0 or nc >= W:
                    uf.union(cur, RIGHT_TOP)
                elif nr >= H or nc < 0:
                    uf.union(cur, LEFT_BOTTOM)

    canReach = uf.find(RIGHT_TOP) != uf.find(LEFT_BOTTOM)
    print("Yes" if canReach else "No")
