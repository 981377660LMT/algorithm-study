# 懒删除堆/可删除堆

from heapq import heapify, heappop, heappush
from typing import Generic, Iterable, Literal, Optional, TypeVar

T = TypeVar("T")


class ErasableHeap(Generic[T]):
    __slots__ = ("_data", "_erased")

    def __init__(self, items: Optional[Iterable[T]] = None) -> None:
        self._erased = []
        self._data = [] if items is None else list(items)
        if self._data:
            heapify(self._data)

    def push(self, value: T) -> None:
        heappush(self._data, value)
        self._normalize()

    def pop(self) -> T:
        value = heappop(self._data)
        self._normalize()
        return value

    def peek(self) -> T:
        return self._data[0]

    def remove(self, value: T) -> None:
        """从堆中删除一个元素,要保证堆中存在该元素."""
        heappush(self._erased, value)
        self._normalize()

    def clear(self) -> None:
        self._data.clear()
        self._erased.clear()

    def __len__(self) -> int:
        return len(self._data)

    def __getitem__(self, index: Literal[0]) -> T:
        return self._data[index]

    def _normalize(self) -> None:
        while self._data and self._erased and self._data[0] == self._erased[0]:
            heappop(self._data)
            heappop(self._erased)


if __name__ == "__main__":
    pq = ErasableHeap((3, 5, 4, 1, 2))
    print(pq.pop())  # 1
    print(pq.pop())  # 2
    pq.remove(3)
    print(pq.pop())  # 4
    pq.remove(5)
    print(len(pq))  # 1
    pq.push(5)
    num = 0
    print(pq[num])
