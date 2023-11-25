from typing import Generic, List, Optional, TypeVar


T = TypeVar("T")


class RollbackArray(Generic[T]):
    __slots__ = ("_n", "_data", "_history")

    def __init__(self, data: Optional[List[T]] = None) -> None:
        if data is None:
            data = []
        self._n = len(data)
        self._data = data[:]
        self._history = []  # (index, value)

    def getTime(self) -> int:
        return len(self._history)

    def rollback(self, time: int) -> None:
        while len(self._history) > time:
            i, v = self._history.pop()
            self._data[i] = v

    def undo(self) -> bool:
        if not self._history:
            return False
        i, v = self._history.pop()
        self._data[i] = v
        return True

    def __getitem__(self, i: int) -> T:
        return self._data[i]

    def __setitem__(self, i: int, x: T) -> None:
        self._history.append((i, self._data[i]))
        self._data[i] = x

    def __len__(self) -> int:
        return self._n

    def __repr__(self) -> str:
        return f"RollbackArray({self._data})"

    def __iter__(self):
        return iter(self._data)


if __name__ == "__main__":
    arr = RollbackArray([1, 2, 3, 4, 5])
    print(arr)
    arr[0] = 10
    print(arr)
    arr[1] = 20
    print(arr)
    state = arr.getTime()
    arr[2] = 30
    print(arr)
    arr.rollback(state)
    print(arr)
