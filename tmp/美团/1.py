# 线性规划

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


T = int(input())
for _ in range(T):
    x, y = map(int, input().split())  # x个A点心 y个B点心
    upper = (x + y) // 3
    print(min(x, y, upper))
