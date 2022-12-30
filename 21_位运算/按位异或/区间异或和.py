# 区间异或和

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def preXor(upper: int) -> int:
    """[0, upper]内所有数的异或"""
    mod = upper % 4
    if mod == 0:
        return upper
    if mod == 1:
        return 1
    if mod == 2:
        return upper + 1
    return 0


def rangeXor(lower: int, upper: int) -> int:
    """[lower, upper]内所有数的异或"""
    return preXor(upper) ^ preXor(lower - 1)


if __name__ == "__main__":
    lower, upper = map(int, input().split())
    print(rangeXor(lower, upper))
