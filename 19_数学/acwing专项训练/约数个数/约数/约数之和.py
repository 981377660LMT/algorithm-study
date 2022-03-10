# 约数之和
from collections import Counter
from math import floor


def getPrimeFactors(n: int) -> Counter:
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n ** 0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


n = int(input())
counter = Counter()
for _ in range(n):
    counter += getPrimeFactors(int(input()))

res = 1
for prime, count in counter.items():
    cur = 1
    for _ in range(count):
        cur = (cur * prime) + 1
        cur %= int(1e9 + 7)
    res *= cur
    res %= int(1e9 + 7)

print(res)
