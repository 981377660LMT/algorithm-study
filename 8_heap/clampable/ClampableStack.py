class ClampableStack:
    __slots__ = ("_clampMin", "_total", "_count", "_stack")

    def __init__(self, clampMin: bool):
        """
        clampMin 为true时,支持容器内所有数与x取min;为false时,支持容器内所有数与x取max.
        """
        self._clampMin = clampMin
        self._total = 0
        self._count = 0
        self._stack = []

    def addAndClamp(self, x: int) -> None:
        newCount = 1
        if self._clampMin:
            while self._stack:
                top = self._stack[-1]
                if top[0] < x:
                    break
                self._stack.pop()
                v, c = top
                self._total -= v * c
                newCount += c
        else:
            while self._stack:
                top = self._stack[-1]
                if top[0] > x:
                    break
                self._stack.pop()
                v, c = top
                self._total -= v * c
                newCount += c
        self._total += x * newCount
        self._count += 1
        self._stack.append((x, newCount))

    def sum(self) -> int:
        return self._total

    def __len__(self) -> int:
        return self._count


if __name__ == "__main__":
    cs = ClampableStack(clampMin=True)
    cs.addAndClamp(1)
    cs.addAndClamp(2)
    cs.addAndClamp(1)
    assert cs.sum() == 3
    cs = ClampableStack(clampMin=False)
    cs.addAndClamp(1)
    cs.addAndClamp(2)
    cs.addAndClamp(1)
    assert cs.sum() == 5
    print("clamped stack passed")
