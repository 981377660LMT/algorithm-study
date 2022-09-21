# 涂色
# 一个格子是黑色当且仅当位于斜率为±1的对角线上
# !即abs(x-a)==abs(y-b)

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, a, b = map(int, input().split())
    x1, x2, y1, y2 = map(int, input().split())
    for r in range(x1, x2 + 1):
        row = []
        for c in range(y1, y2 + 1):
            row.append("#" if abs(r - a) == abs(c - b) else ".")
        print("".join(row))
