# from titan_pylib.algorithm.cycle_memo import CycleMemo

from typing import Callable, TypeVar, Generic, List, Set, Hashable

T = TypeVar("T", bound=Hashable)


class CycleMemo(Generic[T]):
    def __init__(self, initial: T, move_to: Callable[[T], T]) -> None:
        history: List[T] = []
        memo: Set[T] = set()
        now: T = initial
        while now not in memo:
            history.append(now)
            memo.add(now)
            now = move_to(now)
        cycle_start = history.index(now)
        cycle_len = len(history) - cycle_start
        pre = history[:cycle_start]
        cycle = history[cycle_start:]
        self.cycle_start = cycle_start
        self.cycle_len = cycle_len
        self.pre = pre
        self.cycle = cycle

    def kth(self, k: int) -> T:
        if k < self.cycle_start:
            return self.pre[k]
        k -= self.cycle_start
        k %= self.cycle_len
        return self.cycle[k]
