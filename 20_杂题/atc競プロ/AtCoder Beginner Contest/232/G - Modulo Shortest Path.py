# 给定一张N个点的有向完全图，其中，i到j的有向边边权为(Ai+ Bj) mod M。
# 问0到N-1的最短路。
# n<=2e5 M<=1e9
# https://atcoder.jp/contests/abc232/tasks/abc232_g

# 技巧:
# !将点和边控制在O(n)级别

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
