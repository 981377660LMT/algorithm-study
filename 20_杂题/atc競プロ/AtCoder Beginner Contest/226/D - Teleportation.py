"""
题意:给定N个小镇的坐标,寻找一个向量集合使得任取两个小镇,
都存在一个向量,使得两镇可以互相到达

标准化向量
ワープ(warp):传送
"""

from math import gcd
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    pos = [tuple(map(int, input().split())) for _ in range(n)]
    slope = set()
    for i in range(n):
        x1, y1 = pos[i]
        for j in range(n):
            if i == j:
                continue
            x2, y2 = pos[j]
            dx, dy = x2 - x1, y2 - y1
            gcd_ = gcd(dx, dy)
            slope.add((dx // gcd_, dy // gcd_))

    print(len(slope))
