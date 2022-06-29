# 求 a*b*c=k 的三元组(a,b,c)的个数 (a<=b<=c)
# 注意:
# !求约数的复杂度为O(n^1/2) 而n的约数个数近似于O(n^1/3)

from math import floor
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def getFactors(n: int) -> List[int]:
    """n 的所有因数"""
    if n <= 0:
        return []
    small, big = [], []
    upper = floor(n ** 0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


k = int(input())
factors = getFactors(k)

res = 0
for i in range(len(factors)):  # 差不多1e4个
    n1 = factors[i]
    for j in range(i, len(factors)):
        n2 = factors[j]
        if k % (n1 * n2) == 0:
            n3 = k // (n1 * n2)
            if n2 <= n3:
                res += 1
print(res)
