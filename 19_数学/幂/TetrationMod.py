# TetrationMod
# 幂的幂(迭代幂次)的模
# A↑↑B(modM)
# !O(sqrt(mod))

from typing import List


def _getPrimes(m=10**9) -> List[int]:
    upper = int(m**0.5 + 1)
    visited = [False] * (upper + 1)
    i = 2
    while i * i <= upper:
        if not visited[i]:
            for j in range(i * i, upper + 1, i):
                visited[j] = True
        i += 1
    return [i for i in range(2, upper + 1) if not visited[i]]


P = _getPrimes(int(1e9 + 10))


def tetrationMod(a: int, b: int, mod: int) -> int:
    """求a^(a^(a^...^a)) (b个a)的值模mod.
    0<=a,b<=1e9
    1<=mod<=1e9
    """

    def totient(n):
        res = n
        for p in P:
            if p * p > n:
                break
            if n % p == 0:
                res = res // p * (p - 1)
                while n % p == 0:
                    n //= p
        if n != 1:
            res = res // n * (n - 1)
        return res

    def mpow(a, p, m):
        res, flg = 1 % m, 1
        while p:
            if p & 1:
                res *= a
                if res >= m:
                    flg = 0
                    res %= m
            if p == 1:
                break
            a *= a
            if a >= m:
                flg = 0
                a %= m
            p >>= 1
        return res, flg

    def calc(rec, a, b, m):
        if a == 0:
            return (b & 1) ^ 1, 1
        if a == 1:
            return 1, 1
        if m == 1:
            return 0, 0
        if b == 0:
            return 1, 1
        if b == 1:
            return a % m, int(a < m)
        phi_m = totient(m)
        pre, flg1 = rec(rec, a, b - 1, phi_m)
        if flg1:
            res, flg2 = mpow(a % m, pre, m)
        else:
            res, flg2 = mpow(a % m, pre + phi_m, m)
        return res, flg1 & flg2

    return calc(calc, a, b, mod)[0] % mod


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        a, b, m = map(int, input().split())
        print(tetrationMod(a, b, m))
