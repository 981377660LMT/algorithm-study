# 三维空间中n次行走 每次需要选择一个方向走一步(六种方向)
# !到达(x,y,z) 有多少种方法
# 这道题会卡fft求卷积

# !需要先求二维 再求三维
# https://atcoder.jp/contests/abc240/tasks/abc240_f

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n, x, y, z = map(int, input().split())


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
