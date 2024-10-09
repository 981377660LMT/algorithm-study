from collections import deque


class SizedQueue:
    """带容量限制的队列."""

    __slots__ = ("_sum", "_size", "_queue")

    def __init__(self, size: int):
        self._sum = 0
        self._size = size
        self._queue = deque()

    def append(self, v: int):
        self._queue.append(v)
        self._sum += v
        if len(self._queue) > self._size:
            self._sum -= self._queue.popleft()

    def popleft(self) -> int:
        res = self._queue.popleft()
        self._sum -= res
        return res

    def sum(self) -> int:
        return self._sum

    def __getitem__(self, i: int) -> int:
        return self._queue[i]

    def __len__(self) -> int:
        return len(self._queue)

    def __str__(self) -> str:
        return f"SizedQueue{list(self._queue)}"


if __name__ == "__main__":
    Q = SizedQueue(3)
    Q.append(1)
    Q.append(2)
    Q.append(3)
    Q.append(4)
    print(Q.sum())  # 9
    print(Q)
