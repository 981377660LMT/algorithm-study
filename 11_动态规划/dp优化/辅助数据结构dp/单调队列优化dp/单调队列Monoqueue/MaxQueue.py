from typing import Any, Deque, Iterable, List, Optional
from collections import deque


class MaxQueue:
    __slot__ = ("maxQueue", "rawQueue")

    def __init__(self, iterable: Optional[Iterable[int]] = None) -> None:
        self.maxQueue: Deque[List[Any]] = deque()
        self.rawQueue: Deque[int] = deque()

        if iterable is not None:
            for value in iterable:
                self.append(value)

    @property
    def max(self) -> int:
        if not self.maxQueue:
            raise ValueError("monoQueue is empty")
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError("popleft from empty queue")

        self.maxQueue[0][-1] -= 1
        if self.maxQueue[0][-1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int, *metaInfo: Any) -> None:
        """
        Args:
            value (int): 元素的值
            metaInfo: Any 当前元素附加的元信息，不会添加到原始队列
        """

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[-1]
        self.maxQueue.append([value, *metaInfo, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)

    def __getitem__(self, index: int) -> int:
        return self.rawQueue[index]

    def __repr__(self) -> str:
        return f"MaxQueue({self.rawQueue})"


class MinQueue:
    __slots__ = ("minQueue", "rawQueue")

    def __init__(self, iterable: Optional[Iterable[int]] = None) -> None:
        self.minQueue: Deque[List[Any]] = deque()
        self.rawQueue: Deque[int] = deque()

        if iterable is not None:
            for value in iterable:
                self.append(value)

    @property
    def min(self) -> int:
        if not self.minQueue:
            raise ValueError("monoQueue is empty")
        return self.minQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError("popleft from empty queue")

        self.minQueue[0][-1] -= 1
        if self.minQueue[0][-1] == 0:
            self.minQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int, *metaInfo: Any) -> None:
        """
        Args:
            value (int): 元素的值
            metaInfo: Any 当前元素附加的元信息，不会添加到原始队列
        """
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[-1]
        self.minQueue.append([value, *metaInfo, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)

    def __getitem__(self, index: int) -> int:
        return self.rawQueue[index]

    def __repr__(self) -> str:
        return f"MinQueue({self.rawQueue})"


if __name__ == "__main__":
    maxQueue = MaxQueue()
    maxQueue.append(1)
    maxQueue.append(2)
    assert maxQueue.max == 2
    maxQueue.popleft()
    assert maxQueue.max == 2

    minQueue = MinQueue()
    minQueue.append(2)
    minQueue.append(1)
    minQueue.append(3)
    assert minQueue.min == 1
    minQueue.popleft()
    assert minQueue.min == 1
    minQueue.popleft()
    assert minQueue.min == 3
