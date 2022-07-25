# 支持即时删除的堆
# TODO


from typing import Any, Generic, List, Protocol, TypeVar


class SupportsComparison(Protocol):
    __slots__ = ()

    def __eq__(self, other: Any, /) -> bool:
        ...

    def __lt__(self, other: Any, /) -> bool:
        ...


T = TypeVar("T", bound=SupportsComparison)


class RemovableHeap(Generic[T]):
    """TODO"""

    def __init__(self):
        self._pq: List[T] = []

    def heappush(self, value: T) -> None:
        ...

    def heappop(self) -> T:
        ...

    def remove(self, value: T) -> None:
        ...

    def __len__(self) -> int:
        ...

    def __contains__(self, value: T) -> bool:
        ...

    def __getitem__(self, index: int) -> T:
        ...

    def __repr__(self) -> str:
        return self._pq.__repr__()
