import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 看到的数和为k(k=0-m)时至少需要翻动多少次 不能则返回-1
# !单调队列优化的01背包dp

if __name__ == "__main__":
    n, m = map(int, input().split())
    A, B = [], []
    for _ in range(n):
        a, b = map(int, input().split())
        A.append(a)
        B.append(b)
