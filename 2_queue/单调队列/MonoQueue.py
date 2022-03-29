from collections import deque


class MonoQueue:
    def __init__(self):
        self.minQueue = deque()
        self.maxQueue = deque()
        self.rawQueue = deque()

    @property
    def min(self) -> int:
        return self.minQueue[0][0] if self.minQueue else None

    @property
    def max(self) -> int:
        return self.maxQueue[0][0] if self.maxQueue else None

    def popleft(self) -> int:
        if not self.rawQueue:
            raise IndexError('popleft from empty queue')

        self.minQueue[0][1] -= 1
        if self.minQueue[0][1] == 0:
            self.minQueue.popleft()

        self.maxQueue[0][1] -= 1
        if self.maxQueue[0][1] == 0:
            self.maxQueue.popleft()

        return self.rawQueue.popleft()

    def append(self, value: int) -> None:
        count = 1
        while self.minQueue and self.minQueue[-1][0] > value:
            count += self.minQueue.pop()[1]
        self.minQueue.append([value, count])

        count = 1
        while self.maxQueue and self.maxQueue[-1][0] < value:
            count += self.maxQueue.pop()[1]
        self.maxQueue.append([value, count])

        self.rawQueue.append(value)

    def __len__(self) -> int:
        return len(self.rawQueue)


if __name__ == '__main__':
    monoQueue = MonoQueue()
    monoQueue.append(1)
    monoQueue.append(2)
    monoQueue.append(3)
    assert monoQueue.min == 1
    monoQueue.popleft()
    assert monoQueue.min == 2
    assert monoQueue.max == 3
