# 与环相关的工具函数.


from typing import List, Tuple


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Cycle:
    __slots__ = "n"

    def __init__(self, n: int):
        self.n = n

    def dist(self, u: int, v: int) -> int:
        """环上两点距离."""
        d = abs(u - v)
        return min(d, self.n - d)

    def path(self, u: int, v: int) -> List[int]:
        """环上两点路径."""
        if self.distLeft(u, v) <= self.distRight(u, v):
            return self.pathLeft(u, v)
        return self.pathRight(u, v)

    def segment(self, u: int, v: int) -> List[Tuple[int, int]]:
        if self.distLeft(u, v) <= self.distRight(u, v):
            return self.segmentLeft(u, v)
        return self.segmentRight(u, v)

    def segmentLeft(self, from_: int, to: int) -> List[Tuple[int, int]]:
        if from_ >= to:
            return [(from_, to)]
        return [(from_, 0), (self.n - 1, to)]

    def segmentRight(self, from_: int, to: int) -> List[Tuple[int, int]]:
        if to >= from_:
            return [(from_, to)]
        return [(from_, self.n - 1), (0, to)]

    def pathLeft(self, from_: int, to: int) -> List[int]:
        """逆时针从from_到to的路径经过的点."""
        if from_ >= to:
            return list(range(from_, to - 1, -1))
        # from_ -> 0 -> N-1 -> to
        return list(range(from_, -1, -1)) + list(range(self.n - 1, to - 1, -1))

    def pathRight(self, from_: int, to: int) -> List[int]:
        if to >= from_:
            return list(range(from_, to + 1))
        # from_ -> N-1 -> 0 -> to
        return list(range(from_, self.n)) + list(range(to + 1))

    def distLeft(self, from_: int, to: int) -> int:
        """逆时针从from_到to的距离."""
        return from_ - to if from_ >= to else from_ + self.n - to

    def distRight(self, from_: int, to: int) -> int:
        """顺时针从from_到to的距离."""
        return to - from_ if to >= from_ else to + self.n - from_

    def onPathLeft(self, from_: int, to: int, x: int) -> bool:
        """x是否在from_到to的逆时针路径上."""
        if x < to:
            x += self.n
        if from_ < to:
            from_ += self.n
        return to <= x <= from_

    def onPathRight(self, from_: int, to: int, x: int) -> bool:
        """x是否在from_到to的顺时针路径上."""
        if from_ > to:
            to += self.n
        if from_ > x:
            x += self.n
        return from_ <= x <= to

    def jumpLeft(self, from_: int, steps: int) -> int:
        """逆时针从from_跳steps步后的位置."""
        return (from_ - steps) % self.n

    def jumpRight(self, from_: int, steps: int) -> int:
        """顺时针从from_跳steps步后的位置."""
        return (from_ + steps) % self.n


if __name__ == "__main__":
    # B - Hands on Ring (Easy)
    # https://atcoder.jp/contests/abc376/tasks/abc376_b
    # n的环形格子。两个棋子，初始位于0,1。
    # 给定q个指令，每个指令指定一个棋子移动到某个格子上，期间不能移动另外一个棋子。
    # 依次执行这些指令，问移动的次数。
    N, Q = map(int, input().split())
    H, T = [], []
    for _ in range(Q):
        h, t = input().split()
        H.append(h)
        T.append(int(t) - 1)

    C = Cycle(N)
    posL, posR = 0, 1
    res = 0
    for t, h in zip(T, H):
        if h == "L":
            from_, to = posL, t
            if not C.onPathLeft(from_, to, posR):
                res += C.distLeft(from_, to)
            else:
                res += C.distLeft(to, from_)
            posL = to
        else:
            from_, to = posR, t
            if not C.onPathRight(from_, to, posL):
                res += C.distRight(from_, to)
            else:
                res += C.distRight(to, from_)
            posR = to
    print(res)
