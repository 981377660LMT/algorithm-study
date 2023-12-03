# 用两个栈模拟一个双端队列，注意到每个栈的前缀最小值是可以 O(1) 修改和询问的
# 当某个栈空的时候将另一个栈分一半给这个已经空的栈（即暴力重构）

# 参考 SlidingWindowAggregationDeque


from typing import Tuple


E = Tuple[int, ...]  # (value, index,...)


class MaxDeque:
    __slots__ = ("_left", "_right")

    def __init__(self):
        self._left = _MaxStack()
        self._right = _MaxStack()

    def append(self, x: "E") -> None:
        self._right.append(x)

    def appendleft(self, x: "E") -> None:
        self._left.append(x)

    def pop(self) -> "E":
        if self._right:
            return self._right.pop()
        tmp = []
        n = len(self._left)
        for _ in range(n):
            tmp.append(self._left.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._left.append(tmp[i])
        for i in range(half, n):
            self._right.append(tmp[i])
        return self._right.pop()

    def popleft(self) -> "E":
        if self._left:
            return self._left.pop()
        tmp = []
        n = len(self._right)
        for _ in range(n):
            tmp.append(self._right.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._right.append(tmp[i])
        for i in range(half, n):
            self._left.append(tmp[i])
        return self._left.pop()

    @property
    def max(self) -> int:
        if not self._left:
            return self._right.max
        if not self._right:
            return self._left.max
        return self._left.max if self._left.max > self._right.max else self._right.max

    def __len__(self) -> int:
        return len(self._left) + len(self._right)

    def __getitem__(self, i: int) -> "E":
        n = len(self)
        if i < 0:
            i += n
        if i < 0 or i >= n:
            raise IndexError("deque index out of range")
        if i < len(self._left):
            return self._left[-i - 1]
        return self._right[i - len(self._left)]

    def __repr__(self) -> str:
        sb = []
        for i in range(len(self)):
            sb.append(str(self[i]))
        return f"MaxDeque({', '.join(sb)})"


class _MaxStack:
    __slots__ = ("_stack", "_maxes")

    def __init__(self):
        self._stack = []
        self._maxes = []

    def append(self, x: "E") -> None:
        self._stack.append(x)
        if not self._maxes or x[0] >= self._maxes[-1]:
            self._maxes.append(x[0])

    def pop(self) -> "E":
        res = self._stack.pop()
        if res[0] == self._maxes[-1]:
            self._maxes.pop()
        return res

    def top(self) -> "E":
        return self._stack[-1]

    @property
    def max(self) -> int:
        return self._maxes[-1]

    def __len__(self) -> int:
        return len(self._stack)

    def __getitem__(self, i: int) -> "E":
        return self._stack[i]

    def __repr__(self) -> str:
        return f"{self._stack}"


E = Tuple[int, ...]  # (value, index,...)


class MinDeque:
    __slots__ = ("_left", "_right")

    def __init__(self):
        self._left = _MinStack()
        self._right = _MinStack()

    def append(self, x: "E") -> None:
        self._right.append(x)

    def appendleft(self, x: "E") -> None:
        self._left.append(x)

    def pop(self) -> "E":
        if self._right:
            return self._right.pop()
        tmp = []
        n = len(self._left)
        for _ in range(n):
            tmp.append(self._left.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._left.append(tmp[i])
        for i in range(half, n):
            self._right.append(tmp[i])
        return self._right.pop()

    def popleft(self) -> "E":
        if self._left:
            return self._left.pop()
        tmp = []
        n = len(self._right)
        for _ in range(n):
            tmp.append(self._right.pop())
        half = n // 2
        for i in range(half - 1, -1, -1):
            self._right.append(tmp[i])
        for i in range(half, n):
            self._left.append(tmp[i])
        return self._left.pop()

    @property
    def min(self) -> int:
        if not self._left:
            return self._right.min
        if not self._right:
            return self._left.min
        return self._left.min if self._left.min < self._right.min else self._right.min

    def __len__(self) -> int:
        return len(self._left) + len(self._right)

    def __getitem__(self, i: int) -> "E":
        n = len(self)
        if i < 0:
            i += n
        if i < 0 or i >= n:
            raise IndexError("deque index out of range")
        if i < len(self._left):
            return self._left[-i - 1]
        return self._right[i - len(self._left)]

    def __repr__(self) -> str:
        sb = []
        for i in range(len(self)):
            sb.append(str(self[i]))
        return f"MinDeque({', '.join(sb)})"


class _MinStack:
    __slots__ = ("_stack", "_mins")

    def __init__(self):
        self._stack = []
        self._mins = []

    def append(self, x: "E") -> None:
        self._stack.append(x)
        if not self._mins or x[0] <= self._mins[-1]:
            self._mins.append(x[0])

    def pop(self) -> "E":
        res = self._stack.pop()
        if res[0] == self._mins[-1]:
            self._mins.pop()
        return res

    def top(self) -> "E":
        return self._stack[-1]

    @property
    def min(self) -> int:
        return self._mins[-1]

    def __len__(self) -> int:
        return len(self._stack)

    def __getitem__(self, i: int) -> "E":
        return self._stack[i]

    def __repr__(self) -> str:
        return f"{self._stack}"


if __name__ == "__main__":
    maxDeque = MaxDeque()
    maxDeque.append((1, 0))
    maxDeque.append((3, 1))
    maxDeque.append((2, 2))
    maxDeque.append((4, 3))
    maxDeque.append((5, 4))
    print(maxDeque)
    res = maxDeque.popleft()
    assert res == (1, 0)
    print(res)
    print(maxDeque)
    maxDeque.appendleft((6, 5))
    print(maxDeque)
    res = maxDeque.popleft()
    assert res == (6, 5)
    print(res)
    print(maxDeque)
    assert maxDeque.max == 5
    assert maxDeque[0] == (3, 1)
    assert len(maxDeque) == 4

    minDeque = MinDeque()
    minDeque.append((1, 0))
    minDeque.append((3, 1))
    minDeque.append((2, 2))
    minDeque.append((4, 3))
    minDeque.append((5, 4))
    print(minDeque)
    res = minDeque.popleft()
    assert res == (1, 0)
    print(res)
    print(minDeque)
    minDeque.appendleft((6, 5))
    print(minDeque)
    res = minDeque.popleft()
    assert res == (6, 5)
    print(res)
    print(minDeque)
    assert minDeque.min == 2
    assert minDeque[0] == (3, 1)
    assert len(minDeque) == 4
