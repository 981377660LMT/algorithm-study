# 求x除以10向下取整 -1e18<=x<=1e18
# !注意除法精度 不能用floor(x/10)
# !可以用 x//10 或者 floor(Fraction(x, 10))

# !向下取整:floor(注意精度) 或者 地板除//
# !向零取整:int

from fractions import Fraction
from math import floor, sqrt
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    x = int(input())
    print(x // 10)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
