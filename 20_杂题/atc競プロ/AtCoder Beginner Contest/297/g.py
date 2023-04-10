from functools import lru_cache
from math import ceil
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, L, R = map(int, input().split())
    stones = list(map(int, input().split()))
    # # NIM 找规律???
    # L, R = 3, 5  # 00011122
    # L, R = 4, 5  # 000011112
    # L, R = 4, 6  # 0000111122
    # L, R = 5, 6  # 00000111112
    # L, R = 5, 7  # 000001111122
    # L, R = 6, 7  # 0 x 6 + 1 x 6 + 2 x 1
    # L, R = 6, 9  # 0 x 6 + 1 x 6 + 2 x 3
    # !猜测周期为
    # !(L-1)个0,然后是L个1,2,...(div-1),(R-L)个div
    # 然后是L个0,1,2,...(div-1),(R-L)个div

    div = ceil(R / L)

    def getMex2(state: int) -> int:
        state -= 1
        if state < L - 1:
            return 0
        sum_ = L * (div - 1) + R - 1
        if state < sum_:
            mod_ = (state - L - 1) % L
            return mod_ + 1

        state -= L * (div - 1) + R - 1
        state %= L * (div - 1) + R
        mod_ = state % L
        return mod_ + 1

    f = getMex1 if L == R else getMex2
    grundy = 0
    for stone in stones:
        grundy ^= f(stone)
    print("First" if grundy else "Second")
