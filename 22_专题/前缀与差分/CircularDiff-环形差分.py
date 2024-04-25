# CircularDiff-环形差分


from typing import List


class CircularDiff:
    __slots__ = "_diff", "_n"

    def __init__(self, n: int) -> None:
        self._n = n
        self._diff = [0] * (n + 1)

    def add(self, start: int, end: int, x: int) -> None:
        if start < 0:
            start = 0
        if end > self._n:
            end = self._n
        if start >= end:
            return
        self._diff[start] += x
        self._diff[end] -= x

    def addCircle(self, start: int, end: int, x: int) -> None:
        if start >= end:
            return
        n = self._n
        loop = (end - start) // n
        if loop > 0:
            self.add(0, n, x * loop)
        if (end - start) % n == 0:
            return
        start %= n
        end %= n
        if start < end:
            self.add(start, end, x)
        else:
            self.add(start, n, x)
            if end > 0:
                self.add(0, end, x)

    def build(self) -> None:
        for i in range(1, self._n + 1):
            self._diff[i] += self._diff[i - 1]

    def get(self, i: int) -> int:
        return self._diff[i]

    def getAll(self) -> List[int]:
        return self._diff[:-1]


if __name__ == "__main__":
    D = CircularDiff(10)
    D.addCircle(4, 22, 1)
    D.build()
    print(D.getAll())
