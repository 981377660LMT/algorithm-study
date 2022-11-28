import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 200 という整数が大好きなりんごさんのために、次の問題を解いてください。
# N 個の正整数からなる数列 A が与えられるので、以下の条件をすべて満たす整数の組 (i,j) の個数を求めてください。

# 1≤i<j≤N
# A
# i
# ​
#  −A
# j
# ​
#   は 200 の倍数である。


def findAll(string: str, target: str) -> List[int]:
    """找到所有匹配的字符串起始位置"""
    start = 0
    res = []
    while True:
        pos = string.find(target, start)
        if pos == -1:
            break
        else:
            res.append(pos)
            start = pos + 1

    return res


if __name__ == "__main__":
    s = input()
    print(s.count("ZONe"))
