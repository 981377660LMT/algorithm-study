# E - Permute K times
# https://atcoder.jp/contests/abc367/tasks/abc367_e
# !给定数组x,a，进行k次操作。
# 每次操作，求数组bi=axi，然后b=a。问k次操作后的数组a。


from typing import Callable, Tuple


class Doubling:
    __slots__ = ("_n", "_log", "_to")

    def __init__(self, n: int, maxStep: int) -> None:
        self._n = n
        if maxStep <= 0:
            maxStep = 1
        self._log = maxStep.bit_length()
        self._to = [[-1] * n for _ in range(self._log)]

    def add(self, from_: int, to: int) -> None:
        """初始状态:从 `from` 状态可转移到 `to` 状态.

        0 <= from,to < n
        """
        if to < -1 or to >= self._n:
            raise Exception("to is out of range")
        self._to[0][from_] = to

    def build(self) -> None:
        for k in range(self._log - 1):
            for v in range(self._n):
                w = self._to[k][v]
                if w == -1:
                    self._to[k + 1][v] = -1
                    continue
                self._to[k + 1][v] = self._to[k][w]

    def jump(self, from_: int, step: int) -> int:
        """从 `from` 状态开始，执行 `step` 次操作，返回最终到达的状态.

        0 <= from < n
        如果最终状态不存在，返回 -1
        """
        if step >= 1 << self._log:
            raise Exception("step is over max step")
        to = from_
        for k in range(self._log):
            if to == -1:
                break
            if step & (1 << k):
                to = self._to[k][to]
        return to

    def maxStep(self, from_: int, check: Callable[[int], bool]) -> Tuple[int, int]:
        """求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号."""
        step = 0
        to = from_
        for k in range(self._log - 1, -1, -1):
            tmp = self._to[k][to]
            if to == -1:
                continue
            if check(tmp):
                step |= 1 << k
                to = tmp
        return step, to


if __name__ == "__main__":
    N, K = map(int, input().split())
    P = list(map(int, input().split()))
    A = list(map(int, input().split()))
    for i in range(len(P)):
        P[i] -= 1

    D = Doubling(N, K + 1)
    for i in range(N):
        D.add(i, P[i])
    D.build()
    order = [D.jump(i, K) for i in range(N)]
    res = [A[order[i]] for i in range(N)]
    print(" ".join(map(str, res)))
