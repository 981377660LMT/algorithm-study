from random import randint
from typing import List, Tuple


DIR4 = [(-1, 0), (0, 1), (1, 0), (0, -1)]  # 上右下左


def generateMaze(row: int, col: int) -> Tuple[List[Tuple[int, int]], "UnionFindArraySimple"]:
    """并查集生成迷宫.

    生成的迷宫是一个树.
    返回树的 row*col - 1 条边以及并查集.
    """
    uf = UnionFindArraySimple(row * col)
    edges = []
    for r in range(row):
        for c in range(col):
            cur = r * col + c
            curRoot = uf.find(cur)
            other = []
            for dr, dc in DIR4:
                nr, nc = r + dr, c + dc
                next_ = nr * col + nc
                if 0 <= nr < row and 0 <= nc < col and curRoot != uf.find(next_):
                    other.append(next_)
            if other:
                rand = randint(0, len(other) - 1)
                uf.union(cur, other[rand])
                edges.append((cur, other[rand]))
    return edges, uf


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(self, key1: int, key2: int) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


if __name__ == "__main__":
    row, col = 3, 3
    maze = generateMaze(row, col)
    print(maze)
