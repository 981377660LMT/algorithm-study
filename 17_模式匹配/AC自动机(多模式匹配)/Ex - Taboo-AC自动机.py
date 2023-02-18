#  一个文本串，可以将一个字符改成*，
#  保证都是小写字母组成。问最少改多少次，可以完全不包含若干模式串。
# !AC自动机(多模式串匹配)
# https://atcoder.jp/contests/abc268/submissions/34752897

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    from AutoMaton import AhoCorasick

    s = input()
    n = int(input())
    bad = [input() for _ in range(n)]
