# 元数组只有o和x
# 求同时包含o 和 x 的子数组数
# n<=1e6
# 所有子数组数-只包含o的子数组数-只包含x的子数组数
from itertools import groupby
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n = int(input())
s = input()
groups = [len(list(group)) for _, group in groupby(s)]
res = n * (n + 1) // 2
for len_ in groups:
    res -= len_ * (len_ + 1) // 2
print(res)

