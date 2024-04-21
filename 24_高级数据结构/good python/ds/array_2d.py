# from titan_pylib.data_structures.array.array_2D import Array2D


from typing import TypeVar, Generic, List

T = TypeVar("T")


class Array2D(Generic[T]):
    def __init__(self, n: int, m: int, val: T) -> None:
        self.n: int = n
        self.m: int = m
        self.val: T = val
        self.a: List[T] = [val for _ in range(n * m)]

    def set(self, i: int, j: int, val: T) -> None:
        assert (
            0 <= i < self.n and 0 <= j < self.m
        ), f"IndexError: {self.__class__.__name__}.set({i}, {j}, {val})"
        self.a[i * self.m + j] = val

    def get(self, i: int, j: int) -> T:
        assert (
            0 <= i < self.n and 0 <= j < self.m
        ), f"IndexError: {self.__class__.__name__}.get({i}, {j})"
        return self.a[i * self.m + j]

    def tolist(self) -> List[List[T]]:
        a = [[self.val] * self.m for _ in range(self.n)]
        for i in range(self.n):
            for j in range(self.m):
                a[i][j] = self.get(i, j)
        return a

    def __str__(self) -> str:
        return str(self.tolist())

    __repr__ = __str__
