# https://atcoder.jp/contests/abc240/tasks/abc240_f
# 三维空间中n次行走 每次需要选择一个方向走一步(六种方向)
# !到达(x,y,z) 有多少种方法
# n<=1e7
# -1e7<=x,y,z<=1e7

# !生成函数


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n, x, y, z = map(int, input().split())
    if ((x + y + z) & 1) ^ (n & 1):
        print(0)
        exit(0)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
