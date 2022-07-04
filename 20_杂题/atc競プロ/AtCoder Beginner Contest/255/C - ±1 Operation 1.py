# 有一个首项A 公差D 项数N 的等差数列
# 给你一个整数X 求将X变为数列某项的最小操作次数 每次操作可以是+1或-1
# n<=1e12
# -1e18<=X,a<=1e18


import sys
import os

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    """看X落在那个段内 找左右邻居"""
    x, a, d, n = map(int, input().split())
    if d < 0:
        x, a, d = -x, -a, -d

    if d == 0 or x <= a:
        print(abs(x - a))
        exit(0)

    diff = x - a
    pos = diff // d  # 确定在哪一段 然后找左右
    left = abs(a + d * min(n - 1, pos) - x)
    right = abs(a + d * min(n - 1, (pos + 1)) - x)
    print(min(left, right))


if os.environ.get("USERNAME", "") == "caomeinaixi":
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()
