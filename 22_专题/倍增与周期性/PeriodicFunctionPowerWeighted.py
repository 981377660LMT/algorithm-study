from typing import Callable, TypeVar, Generic, Hashable, Tuple, List

T = TypeVar("T", bound=Hashable)
E = TypeVar("E")


class PeriodicFunctionPowerWeighted(Generic[T, E]):
    __slots__ = (
        "cycleStart",
        "cycleLen",
        "preStates",
        "cycleStates",
        "preWeights",
        "cycleWeights",
        "cycleSum",
        "_e",
        "_op",
        "_pow",
    )

    def __init__(
        self,
        e: Callable[[], E],
        op: Callable[[E, E], E],
        pow: Callable[[E, int], E],
        s0: T,
        next: Callable[[T], Tuple[T, E]],
    ) -> None:
        self._e = e
        self._op = op
        self._pow = pow

        states: List[T] = []
        weights: List[E] = []
        seen: dict[T, int] = dict()
        now = s0

        while now not in seen:
            seen[now] = len(states)
            states.append(now)
            nxt, w = next(now)
            weights.append(w)
            now = nxt

        self.cycleStart = seen[now]
        self.cycleLen = len(states) - self.cycleStart

        self.preStates = states[: self.cycleStart]
        self.cycleStates = states[self.cycleStart :]

        self.preWeights: List[E] = [e()]
        acc = e()
        for i in range(self.cycleStart):
            acc = op(acc, weights[i])
            self.preWeights.append(acc)

        self.cycleWeights: List[E] = [e()]
        acc = e()
        for i in range(self.cycleLen):
            acc = op(acc, weights[self.cycleStart + i])
            self.cycleWeights.append(acc)
        self.cycleSum = acc

    def kth(self, k: int) -> Tuple[T, E]:
        """
        返回第 k 步后的 (state, totalWeight)
        """
        if k < self.cycleStart:
            return self.preStates[k], self.preWeights[k]

        k -= self.cycleStart
        full, rem = k // self.cycleLen, k % self.cycleLen
        total = self.preWeights[self.cycleStart]
        if full > 0:
            total = self._op(total, self._pow(self.cycleSum, full))
        total = self._op(total, self.cycleWeights[rem])

        state = self.cycleStates[rem]
        return state, total


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc241/tasks/abc241_e
    def puttingCandies():
        N, K = map(int, input().split())
        A = list(map(int, input().split()))

        e = lambda: 0
        op = lambda x, y: x + y
        pow = lambda x, y: x * y
        s0 = 0

        def next(x: int) -> Tuple[int, int]:
            weight = A[x % N]
            return (x + weight) % N, weight

        P = PeriodicFunctionPowerWeighted(e=e, op=op, pow=pow, s0=s0, next=next)
        _, total = P.kth(K)
        print(total)

    from typing import List

    class Solution:
        def wordsTyping(self, sentence: List[str], rows: int, cols: int) -> int:
            """
            计算给定句子在 rows x cols 屏幕上能完整显示的次数。
            时间复杂度 O(n + rows)，空间复杂度 O(n)。
            """
            n = len(sentence)
            nextIndex = [0] * n  # 当前行以 sentence[i] 开头时，下一行应从哪个单词开始
            times = [0] * n  # 当前行以 sentence[i] 开头时，完整放下整句的次数

            for i in range(n):
                count, width = 0, 0  # 完整放下整句的次数，当前行已用宽度
                j = i  # 当前单词下标
                # 不断向当前行装入单词 + 一个空格，直至装不下为止
                while width + len(sentence[j]) <= cols:
                    width += len(sentence[j]) + 1
                    j += 1
                    if j == n:
                        j = 0
                        count += 1
                nextIndex[i] = j
                times[i] = count

            e = lambda: 0
            op = lambda x, y: x + y
            pow = lambda x, y: x * y
            P = PeriodicFunctionPowerWeighted(e, op, pow, 0, lambda x: (nextIndex[x], times[x]))
            _, total = P.kth(rows)
            return total
