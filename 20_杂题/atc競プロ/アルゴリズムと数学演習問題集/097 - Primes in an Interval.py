# 闭区间[left,right]内的质数个数
# 1≤L≤R≤1e12
# R−L≤500000
# !在范围里筛素数 闭区间内的质数个数
# 区间内的质数个数

from math import ceil, sqrt


L, R = map(int, input().split())
isPrime = [True] * (R - L + 1)  # P[i] := i+Lが素数か？
if L == 1:
    isPrime[0] = False

last = int(sqrt(R))
for fac in range(2, last + 1):
    start = fac * max(ceil(L / fac), 2) - L  # !A 以上の最小の fac の倍数
    while start < len(isPrime):
        isPrime[start] = False
        start += fac
print(sum(isPrime))
