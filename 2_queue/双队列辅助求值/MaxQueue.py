from collections import deque


class MaxQueue:
    def __init__(self):
        self.monoQueue = deque()
        self.data = deque()

    def getMax(self) -> int:
        return self.monoQueue[0] if self.monoQueue else -1

    def append(self, value: int) -> None:
        while self.monoQueue and self.monoQueue[-1] < value:
            self.monoQueue.pop()
        self.monoQueue.append(value)
        self.data.append(value)

    def popleft(self) -> int:
        if not self.monoQueue:
            return -1
        res = self.data.popleft()
        if res == self.monoQueue[0]:
            self.monoQueue.popleft()
        return res

