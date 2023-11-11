# 有0有9的数组个数(模1e9+7) 数组每个元素 0<=a[i]<=9
# !容斥原理 减去没有0的情况 减去没有9的情况 加上没有0也没有9的情况
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


if __name__ == "__main__":
    n = int(input())
    print((pow(10, n, MOD) - pow(9, n, MOD) - pow(9, n, MOD) + pow(8, n, MOD)) % MOD)
