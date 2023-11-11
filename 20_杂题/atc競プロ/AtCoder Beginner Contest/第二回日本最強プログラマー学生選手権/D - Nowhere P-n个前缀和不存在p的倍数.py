# 一个长为n的数组 每个元素可以取[1,p-1]中的某个数
# 要求n个前缀和不存在p的倍数 求这样的数组个数
# n<=1e9 2<=p<=1e9

# !尝试发现当第一个数固定时,后面每个数恰好有p-2种取法


import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def nowhereP(n: int, p: int) -> int:
    return pow(p - 2, n - 1, MOD) * (p - 1) % MOD


if __name__ == "__main__":
    n, p = map(int, input().split())
    print(nowhereP(n, p))
