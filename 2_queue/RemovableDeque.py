from collections import deque
from typing import Generic, Optional, Sequence, TypeVar


V = TypeVar("V")


class RemovableDeque(Generic[V]):
    __slots__ = ("_queue", "_removed", "_len")

    def __init__(self, values: Optional[Sequence[V]] = None) -> None:
        self._queue = deque(values) if values is not None else deque()
        self._removed = dict()

    def append(self, value: V) -> None:
        self._queue.append(value)

    def appendleft(self, value: V) -> None:
        self._queue.appendleft(value)

    def pop(self) -> V:
        ...

    def popLeft(self) -> V:
        self._refresh()
        return self._queue.popleft()

    def get(self, index: int) -> V:
        self._refresh()
        return self._queue[index]

    def remove(self, value: V) -> None:
        """删除前必须保证value存在于队列."""
        self._removed[value] = self._removed.get(value, 0) + 1

    def empty(self) -> bool:
        return self._len == 0

    def __len__(self) -> int:
        return self._len

    def _refresh(self) -> None:
        while self._removedQueue and self._removedQueue[0] == self._queue[0]:
            self._removedQueue.popleft()
            self._queue.popleft()


if __name__ == "__main__":
    q = RemovableDeque([1, 2, 3, 4, 5])
    q.remove(1)
    q.remove(3)
    q.remove(5)
    print(q.popLeft())
    print(q.popLeft())
