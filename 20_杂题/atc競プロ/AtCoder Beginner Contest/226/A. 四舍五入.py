# 题意：给定一个小数，四舍五入到附近的整数
# !python 四舍五入 py的round函数精度低，尽量避免使用 (python里round和js里round不一样)
# print(round(99.500)) 100
# print(round(42.500)) 42
# print(floor(42.500+0.5)) 43

from math import floor
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    # !四舍五入小数使用floor(f+0.5)
    f = float(input())
    print(floor(f + 0.5))
