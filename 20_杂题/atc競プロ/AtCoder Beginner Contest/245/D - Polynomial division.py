# 多项式除法 Polynomial division
# 给定两个多项式 A C 分别是n次多项式 和 n+m 次多项式
# 求 B=C/A 的系数

import sys
import os

import numpy as np

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    # !多项式除法 numpy
    poly1 = np.poly1d(list(map(int, input().split()))[::-1])
    poly2 = np.poly1d(list(map(int, input().split()))[::-1])
    div, _mod = poly2 / poly1
    print(*map(int, div.coef.tolist()[::-1]))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
