from sortedcontainers import SortedList
from collections import defaultdict
from typing import Any, Callable, Generic, Iterable, Optional, Protocol, TypeVar


class SupportsComparison(Protocol):
    def __eq__(self, other: Any, /) -> bool:
        ...

    def __lt__(self, other: Any, /) -> bool:
        ...

    def __gt__(self, other: Any, /) -> bool:
        return (not self < other) and self != other

    def __le__(self, other: Any, /) -> bool:
        return self < other or self == other

    def __ge__(self, other: Any, /) -> bool:
        return not self < other


T = TypeVar("T")


class MultiSet(SortedList, Generic[T]):
    def __init__(
        self,
        iterable: Optional[Iterable[T]] = None,
        key: Optional[Callable[[Any], SupportsComparison]] = None,
    ):
        super().__init__(iterable, key)
        self._counter = defaultdict(int)
        for item in self:
            self._counter[item] += 1

    def add(self, value: T) -> None:
        raise NotImplementedError("use ``add_with_count`` instead")

    def add_with_count(self, value: T, count: int) -> None:
        assert count >= 0, "count must be non-negative"
        if count == 0:
            return
        if self._counter[value] == 0:
            super().add(value)
        self._counter[value] += count

    def remove(self, value: T) -> None:
        raise NotImplementedError("use ``remove_with_count`` instead")

    def remove_with_count(self, value: T, count: int) -> bool:
        assert count >= 0, "count must be non-negative"
        remove_count = min(count, self._counter[value])
        if remove_count == 0:
            return False
        self._counter[value] -= remove_count
        if self._counter[value] == 0:
            super().remove(value)
        return True

    def count(self, value: T) -> int:
        return self._counter[value]

    def __contains__(self, value: T) -> bool:
        return self._counter[value] > 0
