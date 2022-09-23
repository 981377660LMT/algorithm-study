# 输出一个[A,B]之间的C的倍数
# 如果不存在输出-1

from math import ceil, floor
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    A, B, C = map(int, input().split())
    lower, upper = ceil(A / C) * C, floor(B / C) * C
    if lower <= upper:
        print(upper)
    else:
        print(-1)
