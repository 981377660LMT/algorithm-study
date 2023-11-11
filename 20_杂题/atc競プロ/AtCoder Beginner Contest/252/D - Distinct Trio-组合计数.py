"""
n,ai<=2e5
求ai aj ak 不同的三元组数 (i,j,k) 其中 i<j<k

三种的情况还好说,如果是k元组,就需要dp了
"""

from collections import defaultdict
from math import comb
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    indexMap = defaultdict(list)
    for i, num in enumerate(nums):
        indexMap[num].append(i)

    res = comb(n, 3)
    for indexes in indexMap.values():
        # 两同
        res -= comb(len(indexes), 2) * (n - len(indexes))
        # 三同
        res -= comb(len(indexes), 3)

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
