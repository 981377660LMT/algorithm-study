#  一个文本串，可以将一个字符改成*，
#  保证都是小写字母组成。问最少改多少次，可以完全不包含若干模式串。
# !AC自动机(多模式串匹配)

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

s = input()
n = int(input())
words = [input() for _ in range(n)]


# TODO
