# 可删除元素的队列.(需要保证队列中不能有重复元素)

from collections import deque
from typing import Generic, Optional, Sequence, TypeVar


V = TypeVar("V")


class RemovableQueue(Generic[V]):
    __slots__ = ("_queue", "_removedQueue")

    def __init__(self, values: Optional[Sequence[V]] = None) -> None:
        self._queue = deque(values) if values is not None else deque()
        self._removedQueue = deque()

    def append(self, value: V) -> None:
        self._queue.append(value)

    def popLeft(self) -> V:
        self._refresh()
        return self._queue.popleft()

    def top(self) -> V:
        self._refresh()
        return self._queue[0]

    def remove(self, value: V) -> None:
        """删除前必须保证value存在于队列."""
        self._removedQueue.append(value)

    def empty(self) -> bool:
        return len(self) == 0

    def __len__(self) -> int:
        return len(self._queue) - len(self._removedQueue)

    def _refresh(self) -> None:
        while self._removedQueue and self._removedQueue[0] == self._queue[0]:
            self._removedQueue.popleft()
            self._queue.popleft()


if __name__ == "__main__":
    q = RemovableQueue([1, 2, 1, 3, 4, 5])
    q.remove(1)
    q.remove(3)
    q.remove(5)
    print(q.popLeft())
    print(q.popLeft())
