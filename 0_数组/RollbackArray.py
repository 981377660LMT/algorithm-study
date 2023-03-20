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

    def time(self) -> int:
        return len(self._history)

    def rollback(self, t: int) -> None:
        for i in range(self.time() - 1, t - 1, -1):
            index, value = self._history[i]
            self._data[index] = value
        self._history = self._history[:t]

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
    state = arr.time()
    arr[2] = 30
    print(arr)
    arr.rollback(state)
    print(arr)
