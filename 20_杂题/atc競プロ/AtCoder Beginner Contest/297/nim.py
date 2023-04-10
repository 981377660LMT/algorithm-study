from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    # n, L, R = map(int, input().split())
    # stones = list(map(int, input().split()))
    # # NIM 找规律???
    # L, R = 3, 5  # 00011122
    # L, R = 4, 5  # 000011112
    # L, R = 4, 6  # 0000111122
    # L, R = 5, 6  # 00000111112
    # L, R = 5, 7  # 000001111122
    # L, R = 6, 7  # 0 x 6 + 1 x 6 + 2 x 1
    L, R = 3, 3  # 001122334
    # L, R = 2, 7  # 001122334
    # !猜测周期为
    # !(L-1)个0+L个1+(R-L)个2
    # 然后是L个0,1,2,...(div-1),(R-L)个div

    def getMex1(state: int) -> int:
        state -= 1
        if state < L - 1:
            return 0
        if L - 1 <= state < (L - 1 + L):
            return 1
        if (L - 1 + L) <= state < (L - 1 + L + R - L):
            return 2
        state -= L + R - 1

        state %= L + R
        if state < L:
            return 0
        if L <= state < (L + L):
            return 1
        if (L + L) <= state < (L + L + R - L):
            return 2

    @lru_cache(None)
    def getMex(state: int) -> int:
        if state < L:
            return 0
        nexts = set()
        for cur in range(L, min(state, R) + 1):
            nexts.add(getMex(state - cur))
        mex = 0
        while mex in nexts:
            mex += 1
        return mex

    for i in range(1, 50):
        print(i, getMex(i), "a")
        # print(i, getMex1(i), "b")
        # assert getMex(i) == getMex1(i), f"{i} {getMex(i)} {getMex1(i)}"
