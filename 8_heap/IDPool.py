from heapq import heappop, heappush


class IDPool:
    __slots__ = ("_reused", "_nextId")

    def __init__(self, startId=0) -> None:
        self._reused = []
        self._nextId = startId

    def alloc(self) -> int:
        if self._reused:
            return heappop(self._reused)
        res = self._nextId
        self._nextId += 1
        return res

    def release(self, id: int) -> None:
        heappush(self._reused, id)

    def reset(self) -> None:
        self._reused.clear()
        self._nextId = 0

    def __len__(self) -> int:
        return self._nextId - len(self._reused)


if __name__ == "__main__":
    pool = IDPool()
    print(pool.alloc())
    print(pool.alloc())
    print(pool.alloc())
    pool.release(1)
    print(pool.alloc())
    print(pool.alloc())
    pool.release(0)
    pool.release(2)
    print(pool.alloc())
    print(pool.alloc())
