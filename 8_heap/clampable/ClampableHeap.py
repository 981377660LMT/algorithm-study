from heapq import heappop, heappush


class ClampableHeap:
    __slots__ = ("_clampMin", "_total", "_count", "_heap")

    def __init__(self, clampMin: bool):
        """
        clampMin:
          为true时,调用Clamp(x)后,容器内所有数最小值被截断(小于x的数变成x);
          为false时,调用Clamp(x)后,容器内所有数最大值被截断(大于x的数变成x).
          如果需要同时支持两种操作,可以使用双端堆.
        """
        self._clampMin = clampMin
        self._total = 0
        self._count = 0
        self._heap = []

    def add(self, x: int) -> None:
        heappush(self._heap, (x if self._clampMin else -x, 1))
        self._total += x
        self._count += 1

    def clamp(self, x: int) -> None:
        newCount = 0
        if self._clampMin:
            while self._heap:
                v, c = self._heap[0]
                if v > x:
                    break
                heappop(self._heap)
                self._total -= v * c
                newCount += c
            self._total += x * newCount
            heappush(self._heap, (x, newCount))
        else:
            while self._heap:
                v, c = self._heap[0]
                v = -v
                if v < x:
                    break
                heappop(self._heap)
                self._total -= v * c
                newCount += c
            self._total += x * newCount
            heappush(self._heap, (-x, newCount))

    def sum(self) -> int:
        return self._total

    def __len__(self) -> int:
        return self._count


if __name__ == "__main__":
    ch = ClampableHeap(clampMin=False)
    ch.add(1)
    ch.add(2)
    ch.add(3)
    assert ch.sum() == 6
    ch.clamp(2)
    assert ch.sum() == 5
    ch.clamp(1)
    ch.add(2)
    assert ch.sum() == 5

    ch = ClampableHeap(clampMin=True)
    ch.add(1)
    ch.add(2)
    ch.add(3)
    assert ch.sum() == 6
    ch.clamp(2)
    assert ch.sum() == 7
    ch.clamp(3)
    ch.add(2)
    assert ch.sum() == 11
