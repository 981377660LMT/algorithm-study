# 继承deque的MaxQueue

from typing import Any
from collections import deque


class MaxQueue(deque):
    @property
    def max(self) -> int:
        if not self:
            raise ValueError('maxQueue is empty')
        return self[0][0]

    def append(self, value: int, *metaInfo: Any) -> None:
        count = 1
        while self and self[-1][0] < value:
            count += self.pop()[-1]
        super().append([value, *metaInfo, count])

    def popleft(self) -> None:
        if not self:
            raise IndexError('popleft from empty queue')

        self[0][-1] -= 1
        if self[0][-1] == 0:
            super().popleft()


class MinQueue(deque):
    @property
    def min(self) -> int:
        if not self:
            raise ValueError('minQueue is empty')
        return self[0][0]

    def append(self, value: int, *metaInfo: Any) -> None:
        count = 1
        while self and self[-1][0] > value:
            count += self.pop()[-1]
        super().append([value, *metaInfo, count])

    def popleft(self) -> None:
        if not self:
            raise IndexError('popleft from empty queue')

        self[0][-1] -= 1
        if self[0][-1] == 0:
            super().popleft()


if __name__ == '__main__':
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

