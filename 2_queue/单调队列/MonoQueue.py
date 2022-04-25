from typing import Iterable, Optional, Tuple
from collections import deque

# MonoQueue是一个多了 O(1)求min和max 这两个api的 deque


class MonoQueue:
    def __init__(self, iterable: Optional[Iterable[int]] = None) -> None:
        self.minQueue = deque()
        self.maxQueue = deque()
        self.rawQueue = deque()
        self.index = 0

        if iterable is not None:
            for value in iterable:
                self.append(value)

    @property
    def min(self) -> int:
        if not self.minQueue:
            raise ValueError('monoQueue is empty')
        return self.minQueue[0][0]

    @property
    def max(self) -> int:
        if not self.maxQueue:
            raise ValueError('monoQueue is empty')
        return self.maxQueue[0][0]

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError('popleft from empty queue')

        self.minQueue[0][1] -= 1
        if self.minQueue[0][1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][1] -= 1
        if self.maxQueue[0][1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()[0]

    def append(self, value: int) -> None:
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[1]
        self.minQueue.append([value, count, self.index])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[1]
        self.maxQueue.append([value, count, self.index])

        self.rawQueue.append((value, self.index))
        self.index += 1

    def __len__(self) -> int:
        return len(self.rawQueue)

    def __getitem__(self, index: int) -> int:
        return self.rawQueue[index][0]


if __name__ == '__main__':
    monoQueue = MonoQueue()
    monoQueue.append(1)
    monoQueue.append(2)
    monoQueue.append(3)
    assert monoQueue.min == 1
    monoQueue.popleft()
    assert monoQueue.min == 2
    assert monoQueue.max == 3
    monoQueue.append(0)
    assert len(monoQueue) == 3
    assert monoQueue[0] == 2

