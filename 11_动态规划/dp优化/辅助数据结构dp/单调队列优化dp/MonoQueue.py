from typing import Callable, TypeVar, Generic
from collections import deque

T = TypeVar("T")


class MonoQueue(Generic[T]):
    """
    单调队列维护滑动窗口最小值.
    单调队列队头元素为当前窗口最小值，队尾元素为当前窗口最大值.
    """

    __slots__ = ("minQueue", "_minQueueCount", "_less", "_len")

    def __init__(self, less: Callable[[T, T], bool]) -> None:
        self.minQueue = deque()
        self._minQueueCount = deque()
        self._less = less
        self._len = 0

    def append(self, value: T) -> "MonoQueue[T]":
        count = 1
        while self.minQueue and self._less(value, self.minQueue[-1]):
            self.minQueue.pop()
            count += self._minQueueCount.pop()
        self.minQueue.append(value)
        self._minQueueCount.append(count)
        self._len += 1
        return self

    def popleft(self) -> None:
        if not self._len:
            raise IndexError("popleft from empty queue")
        self._minQueueCount[0] -= 1
        if self._minQueueCount[0] == 0:
            self.minQueue.popleft()
            self._minQueueCount.popleft()
        self._len -= 1

    def head(self) -> T:
        if not self._len:
            raise ValueError("monoQueue is empty")
        return self.minQueue[0]

    @property
    def min(self) -> "T":
        return self.head()

    def __len__(self) -> int:
        return self._len

    def __repr__(self) -> str:
        class Item:
            __slots__ = ("value", "count")

            def __init__(self, value: T, count: int):
                self.value = value
                self.count = count

            def __repr__(self) -> str:
                return f"[value: {self.value}, count: {self.count}]"

        res = []
        for i in range(len(self.minQueue)):
            res.append(Item(self.minQueue[i], self._minQueueCount[i]))
        return f"MonoQueue({res})"


if __name__ == "__main__":
    monoQueue = MonoQueue[int](lambda x, y: x < y)
    monoQueue.append(1).append(2).append(3).append(4).append(5)
    print(monoQueue, monoQueue.min)
