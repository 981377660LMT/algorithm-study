from typing import List


class Doubling:
    __slots__ = ("_n", "_log", "_to")

    def __init__(self, n: int, maxStep: int) -> None:
        self._n = n
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


if __name__ == "__main__":
    # https://leetcode.cn/problems/prison-cells-after-n-days/
    def move(preState: int) -> int:
        s1, s2 = preState >> 1, preState << 1
        nextState = s1 ^ s2 ^ 0b11111111  # 两个相邻的房间都被占用或都是空的，那么该牢房就会被占用
        nextState &= 0b01111110  # 行中的第一个和最后一个房间无法有两个相邻的房间
        return nextState

    class Solution:
        def prisonAfterNDays(self, cells: List[int], k: int) -> List[int]:
            n = len(cells)
            db = Doubling(1 << n, int(1e9 + 10))
            for i in range(1 << n):
                db.add(i, move(i))
            db.build()
            start = int("".join(map(str, cells)), 2)
            res = db.jump(start, k)
            return [int(res >> i & 1) for i in range(n - 1, -1, -1)]

    print(Solution().prisonAfterNDays([0, 1, 0, 1, 1, 0, 0, 1], 7))
