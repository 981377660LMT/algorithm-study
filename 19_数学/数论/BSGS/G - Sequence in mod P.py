import sys
from math import ceil, sqrt

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def bsgs(base: int, target: int, mod: int) -> int:
    base %= mod
    target %= mod
    if target == 1 or mod == 1:
        return 0
    mp = dict()
    t = ceil(sqrt(mod))
    val = 1
    for i in range(t):
        mp[target * val % mod] = i
        val = val * base % mod

    step = val
    val = step
    for i in range(1, t + 1):
        if val in mp:
            return i * t - mp[val]
        val = val * step % mod
    return -1


def solve(P: int, A: int, B: int, S: int, G: int) -> int:
    if S == G:
        return 0
    if A == 0:
        return 1 if B == G else -1
    if A == 1:
        if B == 0:
            return -1
        return (G - S) * pow(B, -1, P) % P

    invA1 = pow(A - 1, -1, P)
    C = B * invA1 % P
    target = (G + C) % P
    base = (S + C) % P
    if base == 0:
        return 0 if target == 0 else -1
    target = target * pow(base, -1, P) % P
    return bsgs(A, target, P)


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        P, A, B, S, G = map(int, input().split())
        print(solve(P, A, B, S, G))
