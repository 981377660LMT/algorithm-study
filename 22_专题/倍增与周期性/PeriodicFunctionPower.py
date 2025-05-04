# CycleMemo/PeriodicFunctionPower

from typing import Callable, TypeVar, Generic, Hashable


T = TypeVar("T", bound=Hashable)


class PeriodicFunctionPower(Generic[T]):
    """
    >>> P = PeriodicFunctionPower(0, lambda x: 1 + (x + 1) % 3)
    >>> P.cycleStart
    1
    >>> P.cycleLen
    3
    >>> P.preCycle
    [0]
    >>> P.cycle
    [2, 1, 3]
    >>> P.kth(0)
    0
    >>> P.kth(4)
    2
    """

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
    import doctest

    doctest.testmod()

    # https://atcoder.jp/contests/abc167/tasks/abc167_d?lang=ja
    def teleporter():
        _, K = map(int, input().split())
        A = list(map(int, input().split()))

        P = PeriodicFunctionPower(0, lambda x: A[x] - 1)
        print(P.kth(K) + 1)

    # https://atcoder.jp/contests/typical90/tasks/typical90_bf
    def originalCalculator():
        N, K = map(int, input().split())
        MOD = int(1e5)

        def next(x: int) -> int:
            weight = sum(int(d) for d in str(x))
            return (x + weight) % MOD

        P = PeriodicFunctionPower(N, next)
        print(P.kth(K))

    # https://atcoder.jp/contests/abc258/tasks/abc258_e
    def packingPotatoes():
        n, q, limit = map(int, input().split())
        weights = [int(num) for num in input().split()]

        # !滑窗查找每个土豆开始的组能放几个土豆
        div, mod = divmod(limit, sum(weights))  # 注意先要模
        size = [div * n] * n
        right, curSum = 0, 0
        for left in range(n):
            while (tmp := curSum + (weights[right % n])) < mod:
                curSum = tmp
                right += 1
            size[left] += right - left + 1
            curSum -= weights[left]

        P = PeriodicFunctionPower(0, lambda x: (x + size[x]) % n)
        for _ in range(q):
            k = int(input()) - 1
            print(size[P.kth(k)])

    packingPotatoes()
