# 最大栈/最小栈


class MaxStack:
    __slots__ = ("_stack", "_maxes")

    def __init__(self):
        self._stack = []
        self._maxes = []

    def append(self, x: int) -> None:
        self._stack.append(x)
        if not self._maxes or x >= self._maxes[-1]:
            self._maxes.append(x)

    def pop(self) -> int:
        res = self._stack.pop()
        if res == self._maxes[-1]:
            self._maxes.pop()
        return res

    def top(self) -> int:
        return self._stack[-1]

    @property
    def max(self) -> int:
        return self._maxes[-1]

    def __len__(self) -> int:
        return len(self._stack)

    def __getitem__(self, i: int) -> int:
        return self._stack[i]

    def __repr__(self) -> str:
        return f"{self._stack}"


class MinStack:
    __slots__ = ("_stack", "_mins")

    def __init__(self):
        self._stack = []
        self._mins = []

    def append(self, x: int) -> None:
        self._stack.append(x)
        if not self._mins or x <= self._mins[-1]:
            self._mins.append(x)

    def pop(self) -> int:
        res = self._stack.pop()
        if res == self._mins[-1]:
            self._mins.pop()
        return res

    def top(self) -> int:
        return self._stack[-1]

    @property
    def min(self) -> int:
        return self._mins[-1]

    def __len__(self) -> int:
        return len(self._stack)

    def __getitem__(self, i: int) -> int:
        return self._stack[i]

    def __repr__(self) -> str:
        return f"{self._stack}"


if __name__ == "__main__":
    maxStack = MaxStack()
    maxStack.append(1)
    maxStack.append(3)
    assert maxStack.max == 3
    maxStack.append(2)
    assert maxStack.max == 3
    print(maxStack)
    maxStack.pop()
    assert maxStack.max == 3
    maxStack.pop()
    assert maxStack.max == 1
