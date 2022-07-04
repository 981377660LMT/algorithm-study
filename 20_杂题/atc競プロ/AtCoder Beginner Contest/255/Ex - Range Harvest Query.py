"""
每天早上第i棵树上会长i个果实 (1 <= i <= n)
Q次摘果子 第i次摘果的时间是Di天晚上 摘掉[lefti,righti]里所有的果实
求一共摘多少果子 模998244353

Q<=2e5
N,Di<=1e18

线段树+坐标压缩
线段树每个节点记录上一次被更新的时间
涨了 (t2-t1)天  区间[L,R]每天涨 (L+R)(R-L+1)/2

复杂度QlogQ

https://atcoder.jp/contests/abc255/submissions/32404823
"""
# n<=1e18


import sys
import os
from typing import List

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353


def main() -> None:
    n, q = map(int, input().split())
    bit = BIT2(int(1e18) + 10)
    res = 0
    for _ in range(q):
        day, left, right = map(int, input().split())
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", "") == "caomeinaixi":
        while True:
            try:
                main()
            except (EOFError, ValueError):
                break
    else:
        main()

# TODO
