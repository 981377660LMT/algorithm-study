# CycleMemo/PeriodicFunctionPower

from typing import Callable, TypeVar, Generic, Hashable

T = TypeVar("T", bound=Hashable)


class PeriodicFunctionPowerunctionPower(Generic[T]):
    __slots__ = ("cycleStart", "cycleLen", "preCycle", "cycle")

    def __init__(self, initial: T, moveTo: Callable[[T], T]) -> None:
        history = []
        memo = set()
        now = initial
        while now not in memo:
            history.append(now)
            memo.add(now)
            now = moveTo(now)
        cycleStart = history.index(now)
        cycleLen = len(history) - cycleStart
        preCycle = history[:cycleStart]
        cycle = history[cycleStart:]
        self.cycleStart = cycleStart
        self.cycleLen = cycleLen
        self.preCycle = preCycle
        self.cycle = cycle

    def kth(self, k: int) -> T:
        if k < self.cycleStart:
            return self.preCycle[k]
        k -= self.cycleStart
        k %= self.cycleLen
        return self.cycle[k]


if __name__ == "__main__":
    P = PeriodicFunctionPowerunctionPower(0, lambda x: 1 + (x + 1) % 3)
    print(P.cycleStart, P.cycleLen, P.preCycle, P.cycle)
