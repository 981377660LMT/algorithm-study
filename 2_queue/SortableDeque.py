# SortableQueue/SortableDeque
# !可排序队列/可排序双端队列


from collections import deque
from typing import Any, Iterable, Optional
from sortedcontainers import SortedList


class SortableDeque:
    __slots__ = ("_queue", "_sl", "_size")

    def __init__(self, iterable: Optional[Iterable[Any]] = None) -> None:
        self._queue = deque()
        self._sl = SortedList()
        self._size = 0
        for item in iterable or []:
            self.append(item)

    def sort(self) -> None:
        self._sl.update(self._queue)
        self._queue.clear()

    def append(self, x: int) -> None:
        self._queue.append(x)
        self._size += 1

    def appendleft(self, x: int) -> None:
        self._queue.appendleft(x)
        self._size += 1

    def pop(self) -> int:
        self._size -= 1
        if self._sl:
            return self._sl.pop()
        return self._queue.pop()

    def popleft(self) -> int:
        self._size -= 1
        if self._sl:
            return self._sl.pop(0)
        return self._queue.popleft()

    def front(self) -> int:
        return self._sl[0] if self._sl else self._queue[0]

    def back(self) -> int:
        return self._sl[-1] if self._sl else self._queue[-1]

    def __len__(self) -> int:
        return self._size


if __name__ == "__main__":
    # 3 1 4 1 5 9 2 6
    nums = [3, 1, 4, 1, 5, 9, 2, 6]

    sd = SortableDeque(nums)
    sd.sort()
    print(sd.front())
    print(sd.back())
    print(sd.pop())
    print(sd.popleft())
    print(sd.front())
    print(sd.back())
    print(sd.pop())
    print(sd.popleft())
    print(sd.front())
    print(sd.back())
    print(sd.pop())
