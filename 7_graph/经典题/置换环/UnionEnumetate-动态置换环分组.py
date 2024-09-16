# https://noshi91.github.io/Library/data_structure/union_enumerate.cpp
# 比并查集更加轻量的分组

from typing import List


class UnionEnumerate:
    __slots__ = ("_next",)

    def __init__(self, n: int) -> None:
        self._next = list(range(n))

    def getAll(self) -> List[List[int]]:
        n = len(self._next)
        res = []
        visited = [False] * n
        for i in range(n):
            if not visited[i]:
                cur = self.collectGroup(i)
                for j in cur:
                    visited[j] = True
                res.append(cur)
        return res

    def union(self, x: int, y: int) -> None:
        """合并x和y所在的集合,需要保证`x和y不在同一个集合中`."""
        self._next[x], self._next[y] = self._next[y], self._next[x]

    def collectGroup(self, x: int) -> List[int]:
        res = []
        cur = x
        while True:
            res.append(cur)
            cur = self._next[cur]
            if cur == x:
                break
        return res

    def __len__(self) -> int:
        return len(self._next)

    def __repr__(self) -> str:
        res = uf.getAll()
        return f"UnionEnumerate({res})"


if __name__ == "__main__":
    uf = UnionEnumerate(5)
    print(uf)
    uf.union(0, 1)
    print(uf)
