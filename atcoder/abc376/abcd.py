import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    N, Q = map(int, input().split())
    H, T = [], []
    for _ in range(Q):
        h, t = input().split()
        H.append(h)
        T.append(int(t) - 1)

    # def getPath(from_: int, to: int) -> Tuple[List[int], List[int]]:
    #     if from_ == to:
    #         return [from_], [from_]
    #     if from_ > to:
    #         from_, to = to, from_
    #     path1 = list(range(from_, to + 1))
    #     path2 = list(range(to, N)) + list(range(from_ + 1))
    #     return path1, path2

    def dist(from_: int, to: int) -> int:
        if from_ == to:
            return 0
        if from_ > to:
            to += N
        return to - from_

    def onPath(from_: int, to: int, x: int) -> bool:
        """x是否在from_到to的路径上."""
        if from_ == to:
            return False
        if from_ > to:
            to += N
        if from_ > x:
            x += N
        return from_ <= x <= to

    posL, posR = 0, 1
    res = 0
    for i in range(Q):
        to = T[i]
        if H[i] == "L":
            from_, to = posL, to
            if not onPath(from_, to, posR):
                res += dist(from_, to)
            else:
                res += dist(to, from_)
            posL = to

        else:
            from_, to = posR, to
            if not onPath(from_, to, posL):
                res += dist(from_, to)
            else:
                res += dist(to, from_)
            posR = to
    print(res)
