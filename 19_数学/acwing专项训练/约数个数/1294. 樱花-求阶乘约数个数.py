# 给定一个整数 n，求有多少正整数数对 (x,y) 满足 1/x+1/y=1/n!。

"""
本质上是求(n!)^2 的约数个数，对齐进行质因数分解，然
后用乘法原理求解
"""


from typing import List


def getPrimes(upper: int) -> List[int]:
    """筛选出1-upper中的质数"""
    visited = [False] * (upper + 1)
    for num in range(2, upper + 1):
        if visited[num]:
            continue
        for multi in range(num * num, upper + 1, num):
            visited[multi] = True

    return [num for num in range(2, upper + 1) if not visited[num]]


primes = set(getPrimes(int(1e6 + 10)))
counter = [0] * (int(1e6 + 10))
n = int(input())
for p in primes:
    k = 1
    while p**k <= n:
        counter[p] += n // (p**k)
        k += 1

res = 1
MOD = int(1e9 + 7)
for count in counter:
    res *= 2 * count + 1
    res %= MOD

print(res)
