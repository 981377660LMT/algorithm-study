import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)


N, Q = map(int, input().split())
H, T = [], []
for _ in range(Q):
    h, t = input().split()
    H.append(h)
    T.append(int(t) - 1)


def distLeft(from_: int, to: int) -> int:
    return from_ - to if from_ >= to else from_ + N - to


def distRight(from_: int, to: int) -> int:
    return to - from_ if to >= from_ else to + N - from_


def min2(a: int, b: int) -> int:
    return a if a < b else b


def onPathLeft(from_: int, to: int, x: int) -> bool:
    if from_ == to:
        return False
    if x < to:
        x += N
    if from_ < to:
        from_ += N
    return to <= x <= from_


def onPathRight(from_: int, to: int, x: int) -> bool:
    if from_ == to:
        return False
    if from_ > to:
        to += N
    if from_ > x:
        x += N
    return from_ <= x <= to


def moveLeft(cur: int, to: int, other: int) -> Tuple[int, int]:
    """cur移动到to.

    返回距离,新的cur,新的other.
    """

    if not onPathLeft(cur, to, other):
        return distLeft(cur, to), other

    otherTo = to - 1 if to > 0 else N - 1
    return distLeft(cur, to) + distLeft(other, otherTo), otherTo


def moveRight(cur: int, to: int, other: int) -> Tuple[int, int]:
    """cur移动到to.

    返回距离,新的cur,新的other.
    """
    if not onPathRight(cur, to, other):
        return distRight(cur, to), other
    otherTo = to + 1 if to < N - 1 else 0
    return distRight(cur, to) + distRight(other, otherTo), otherTo


memo = dict()


def dfs(index: int, posL: int, posR: int) -> int:
    if index == Q:
        return 0
    hash_ = index << 32 | posL << 16 | posR
    if hash_ in memo:
        return memo[hash_]
    to = T[index]
    res = INF
    if H[index] == "L":
        d1, r1 = moveLeft(posL, to, posR)
        res = min2(res, d1 + dfs(index + 1, to, r1))
        d2, r2 = moveRight(posL, to, posR)
        res = min2(res, d2 + dfs(index + 1, to, r2))
    else:
        d1, l1 = moveLeft(posR, to, posL)
        res = min2(res, d1 + dfs(index + 1, l1, to))
        d2, l2 = moveRight(posR, to, posL)
        res = min2(res, d2 + dfs(index + 1, l2, to))
    memo[hash_] = res
    return res


res = dfs(0, 0, 1)

print(res)
