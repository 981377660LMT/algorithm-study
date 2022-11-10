from math import gcd
from typing import Callable, Generic, TypeVar


E = TypeVar("E")


class SlidingWindowAggrigation(Generic[E]):
    """SlidingWindowAggrigation

    Api:
    1. append value to tail,O(1).
    2. pop value from head,O(1).
    3. query aggrigated value in window,O(1).
    """

    __slots__ = ["_stack0", "_stack1", "_stack2", "_stack3", "_e0", "_e1", "_size", "_op", "_e"]

    def __init__(self, e: Callable[[], E], _op: Callable[[E, E], E]):
        """
        Args:
            e: unit element
            op: merge function
        """
        self._stack0 = []
        self._stack1 = []
        self._stack2 = []
        self._stack3 = []
        self._e = e
        self._e0 = e()
        self._e1 = e()
        self._size = 0
        self._op = _op

    def append(self, value: E) -> None:
        if not self._stack0:
            self._push0(value)
            self._transfer()
        else:
            self._push1(value)
        self._size += 1

    def popleft(self) -> None:
        if not self._size:
            return
        if not self._stack0:
            self._transfer()
        self._stack0.pop()
        self._stack2.pop()
        self._e0 = self._stack2[-1] if self._stack2 else self._e()
        self._size -= 1

    def query(self) -> E:
        return self._op(self._e0, self._e1)

    def _push0(self, value):
        self._stack0.append(value)
        self._e0 = self._op(value, self._e0)
        self._stack2.append(self._e0)

    def _push1(self, value):
        self._stack1.append(value)
        self._e1 = self._op(self._e1, value)
        self._stack3.append(self._e1)

    def _transfer(self):
        while self._stack1:
            self._push0(self._stack1.pop())
        while self._stack3:
            self._stack3.pop()
        self._e1 = self._e()

    @property
    def size(self):
        return self._size


if __name__ == "__main__":
    windowGcd = SlidingWindowAggrigation(lambda: 0, gcd)
    print(windowGcd.query())
    windowGcd.append(4)
    print(windowGcd.query())
    windowGcd.append(6)
    print(windowGcd.query())
    windowGcd.popleft()
    print(windowGcd.query())
    windowGcd.popleft()
    print(windowGcd.query())
    windowGcd.popleft()
    print(windowGcd.size)
