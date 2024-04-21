# from titan_pylib.data_structures.deque.foldable_deque import FoldableDeque
# from titan_pylib.data_structures.stack.foldable_stack import FoldableStack
from typing import Generic, Iterable, TypeVar, Callable, Union, List

T = TypeVar("T")


class FoldableStack(Generic[T]):
    def __init__(self, n_or_a: Union[int, Iterable[T]], op: Callable[[T, T], T], e: T) -> None:
        self._op = op
        self._e = e
        if isinstance(n_or_a, int):
            self._n = n_or_a
            self._a = [e] * self._n
        else:
            n_or_a = list(n_or_a)
            self._n = len(n_or_a)
            self._a = list(n_or_a)
        self._data = [e] * (self._n + 1)
        for i in range(self._n):
            self._data[i + 1] = op(self._data[i], self._a[i])

    def append(self, key: T) -> None:
        self._a.append(key)
        self._data.append(self._op(self._data[-1], key))

    def top(self) -> T:
        return self._a[-1]

    def pop(self) -> T:
        self._data.pop()
        return self._a.pop()

    def all_prod(self) -> T:
        return self._data[-1]

    def prod(self, r: int) -> T:
        return self._data[r]

    def tolist(self) -> List[T]:
        return list(self._a)

    def __len__(self):
        return len(self._a)


from typing import Generic, Iterable, TypeVar, Callable, Union, List

T = TypeVar("T")


class FoldableDeque(Generic[T]):
    def __init__(self, n_or_a: Union[int, Iterable[T]], op: Callable[[T, T], T], e: T) -> None:
        self._op = op
        self._e = e
        self.front: FoldableStack[T] = FoldableStack(0, op, e)
        self.back: FoldableStack[T] = FoldableStack(n_or_a, op, e)

    def _rebuild(self) -> None:
        new = self.front.tolist()[::-1] + self.back.tolist()
        self.front = FoldableStack(new[: len(new) // 2][::-1], self._op, self._e)
        self.back = FoldableStack(new[len(new) // 2 :], self._op, self._e)

    def tolist(self) -> List[T]:
        return self.front.tolist()[::-1] + self.back.tolist()

    def append(self, v: T) -> None:
        self.back.append(v)

    def appendleft(self, v: T) -> None:
        self.front.append(v)

    def pop(self) -> T:
        if not self.back:
            self._rebuild()
        return self.back.pop() if self.back else self.front.pop()

    def popleft(self) -> T:
        if not self.front:
            self._rebuild()
        return self.front.pop() if self.front else self.back.pop()

    def all_prod(self) -> T:
        return self._op(self.front.all_prod(), self.back.all_prod())

    def __len__(self):
        return len(self.front) + len(self.back)
