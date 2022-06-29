# 2 以上 N 以下の整数のうち、K 種類以上の素因数を持つものの個数を求めてください。
# !2-N中 有多少个数 质因子个数大于k
# 2≤N≤1e7
# 1≤K≤8

# 遍历超时 筛法枚举因子 `nloglogn`
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n, k = map(int, input().split())

counter = [0] * (n + 1)  # 每个数的质因数个数 枚举素因子
for i in range(2, n + 1):
    if counter[i] != 0:
        continue
    for j in range(i, n + 1, i):
        counter[j] += 1

print(sum(1 for c in counter if c >= k))
