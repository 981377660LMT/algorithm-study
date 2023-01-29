import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 最大流
# 最小覆盖
if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]
