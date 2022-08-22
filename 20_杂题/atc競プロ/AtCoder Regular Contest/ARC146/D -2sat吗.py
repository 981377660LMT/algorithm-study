import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m, k = map(int, input().split())
Q = []
for _ in range(k):
    i, x, j, y = map(int, input().split())
    i, j = i - 1, j - 1
    Q.append((i, x, j, y))
