import sys

import math
from bisect import bisect_right

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def sieve(n):
    is_prime = [True] * (n + 1)
    is_prime[0] = False
    is_prime[1] = False
    for i in range(2, int(n**0.5) + 1):
        if is_prime[i]:
            for j in range(i * i, n + 1, i):
                is_prime[j] = False
    return [i for i in range(2, n + 1) if is_prime[i]]


N = int(input())


limit_P8 = int(N ** (1 / 8))

limit = int(math.isqrt(N))

primes = sieve(limit)

count_p8 = 0
for p in primes:
    if p > limit_P8:
        break
    if p**8 <= N:
        count_p8 += 1
    else:
        break

count_p2q2 = 0
length = len(primes)
for i, p in enumerate(primes):
    p2 = p * p
    if p2 > N:
        break
    limitQ = int((N // p2) ** 0.5)
    if limitQ < p:
        continue
    idx = bisect_right(primes, limitQ)
    if idx > i + 1:
        count_p2q2 += idx - (i + 1)

answer = count_p8 + count_p2q2
print(answer)
