# 算法学习笔记(40): 原根 - Pecco的文章 - 知乎
# 算法学习笔记(40): 原根 - Pecco的文章 - 知乎
# https://zhuanlan.zhihu.com/p/166043237

# a在模m下的阶是同余方程 a**x = 1 (mod m) 的最小正整数解,记为ord(m,a)
# ord(m,a)一定是phi(m)的因子 ,特别地,当 ord(m,a) = phi(m)时,称a为模m的原根

# 原根的应用:NTT

# 原根
# 任何质数都有原根
from random import randrange
from math import gcd

# 求出一个原根
# https://judge.yosupo.jp/problem/primitive_root
# p<=1e18
def primitive_root_prime(p: int) -> int:
    if p == 2:
        return 1
    pf = factorize(p - 1)
    res = 2
    count = 0
    while True:
        for pi, _ in pf.items():
            if 1 != pow(res, (p - 1) // pi, p):
                count += 1
            else:
                count = 0
                break
        if count == len(pf):
            break
        res = randrange(3, p - 1)

    return res


def lcm(a, b):
    return (a // gcd(a, b)) * b


def miller_rabin(N, bases):
    d, s = N - 1, 0
    while d % 2 == 0:
        d >>= 1
        s += 1
    for a in bases:
        if N <= a:
            return True
        a = pow(a, d, N)
        if a == 1:
            continue
        r = 1
        while a != N - 1:
            if r == s:
                return False
            a = a * a % N
            r += 1
    return True


def is_prime32(N):
    return miller_rabin(N, [2, 7, 61])


def is_prime64(N):
    return miller_rabin(N, [2, 325, 9375, 28178, 450775, 9780504, 1795265022])


def is_prime(N):
    if N <= 1:
        return False
    if N == 2 or N == 3 or N == 5 or N == 7:
        return True
    if N % 2 == 0 or N % 3 == 0 or N % 5 == 0 or N % 7 == 0:
        return False
    if N < 121:
        return True
    if N < 4759123141:
        return is_prime32(N)
    return is_prime64(N)


def find_prime_factor(n):
    m = max(1, int(n**0.125))

    while True:
        c = randrange(n)
        y = k = 0
        g = q = r = 1
        while g == 1:
            x = y
            mr = 3 * r // 4
            while k < mr:
                y = (pow(y, 2, n) + c) % n
                k += 1
            while k < r and g == 1:
                ys = y
                for _ in range(min(m, r - k)):
                    y = (pow(y, 2, n) + c) % n
                    q = q * abs(x - y) % n
                g = gcd(q, n)
                k += m
            k = r
            r <<= 1
        if g == n:
            g = 1
            y = ys
            while g == 1:
                y = (pow(y, 2, n) + c) % n
                g = gcd(abs(x - y), n)
        if g == n:
            continue
        if is_prime(g):
            return g
        elif is_prime(n // g):
            return n // g
        else:
            return find_prime_factor(g)


def factorize(n):
    res = {}
    for p in range(2, 1000):
        if p * p > n:
            break
        if n % p:
            continue
        s = 0
        while n % p == 0:
            n //= p
            s += 1
        res[p] = s

    while not is_prime(n) and n > 1:
        p = find_prime_factor(n)
        s = 0
        while n % p == 0:
            n //= p
            s += 1
        res[p] = s
    if n > 1:
        res[n] = 1

    return res


def carmichael_lambda(n):
    if n == 1:
        return 1
    if n == 2:
        return 1
    if n == 4:
        return 2
    bin_n = bin(n)
    if bin_n.count("1") == 1:
        return 1 << (len(bin_n) - 5)

    pf = factorize(n)
    res = 1
    if len(pf) == 1:
        for k, v in pf.items():
            res = (k - 1) * pow(k, v - 1)
        return res
    else:
        for k, v in pf.items():
            res = lcm(res, carmichael_lambda(pow(k, v)))
        return res


if __name__ == "__main__":
    # 求素数的原根
    # 2<=p<=1e18
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)
    q = int(input())
    for _ in range(q):
        p = int(input())
        print(primitive_root_prime(p))
