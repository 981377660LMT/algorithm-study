# 给定 n 个正整数 ai，请你输出这些数的乘积的约数个数，答案对 109+7 取模。


# 约数个数：(a1+1)(a2+1)...(ak+1)  考虑每个质因子的贡献即可

from collections import Counter
from math import floor

MOD = int(1e9 + 7)


def getPrimeFactors(n: int) -> "Counter[int]":
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


def countOfFactors(counter: "Counter[int]") -> int:
    """返回所有约数个数, counter 为这个数的所有质数因子分解"""
    res = 1
    for count in counter.values():
        res *= count + 1
        res %= MOD
    return res


if __name__ == "__main__":
    n = int(input())
    counter = Counter()
    for _ in range(n):
        counter += getPrimeFactors(int(input()))
    print(countOfFactors(counter))
