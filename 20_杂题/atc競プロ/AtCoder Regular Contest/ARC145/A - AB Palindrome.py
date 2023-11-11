"""
有一个长度为N的字符串S,
字符串中只有 A,B两种字符,
每次都可以将S中两个相邻的字符替换为AB。
问能否将S变成回文字符串。

分类讨论
https://zhuanlan.zhihu.com/p/548093657
"""
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
s = input()

print("No" if (s[0] == "A" and s[-1] == "B" or s == "BA") else "Yes")
