# https://judge.yosupo.jp/submission/98425

# gcd卷积/gcdConvolve
# lcm卷积/lcmConvolve
# n<=1e6
# 0<=ai,bi<998244353

from typing import List


MOD = 998244353


# c[k] = ∑a[i]*b[j] mod 998244353 (gcd(i,j)=k)
def gcd_convolve(a: List[int], b: List[int]) -> List[int]:
    add = lambda x, y: (x + y) % MOD
    inv = lambda x: -x
    mul = lambda x, y: (x * y) % MOD

    a, b = [0] + a, [0] + b
    a = multiple_zeta_transform(a, add)
    b = multiple_zeta_transform(b, add)
    res = [mul(v1, v2) for v1, v2 in zip(a, b)]
    res = multiple_mobius_transform(res, add, inv)
    return res[1:]


# c[k] = ∑a[i]*b[j] mod 998244353 (lcm(i,j)=k)
def lcm_convolve(a: List[int], b: List[int]) -> List[int]:
    add = lambda x, y: (x + y) % MOD
    inv = lambda x: -x
    mul = lambda x, y: (x * y) % MOD

    a, b = [0] + a, [0] + b
    a = divisor_zeta_transform(a, add)
    b = divisor_zeta_transform(b, add)
    res = [mul(v1, v2) for v1, v2 in zip(a, b)]
    res = divisor_mobius_transform(res, add, inv)
    return res[1:]


def multiple_zeta_transform(a, op):
    n = len(a)
    res = a[:]
    prime_table = [1] * n
    for p in range(2, n):
        if not prime_table[p]:
            continue
        i = (n - 1) // p
        while i > 0:
            res[i] = op(res[i], res[i * p])
            prime_table[i * p] = 0
            i -= 1
    return res


def multiple_mobius_transform(a, op, inv):
    n = len(a)
    res = a[:]
    prime_table = [1] * n
    for p in range(2, n):
        if not prime_table[p]:
            continue
        i = 1
        while i * p < n:
            res[i] = op(res[i], inv(res[i * p]))
            prime_table[i * p] = 0
            i += 1
    return res


def divisor_zeta_transform(a, op):
    n = len(a)
    res = a[:]
    prime_table = [1] * n
    for p in range(2, n):
        if not prime_table[p]:
            continue
        i = 1
        while i * p < n:
            res[i * p] = op(res[i], res[i * p])
            prime_table[i * p] = 0
            i += 1
    return res


def divisor_mobius_transform(a, op, inv):
    n = len(a)
    res = a[:]
    prime_table = [1] * n
    for p in range(2, n):
        if not prime_table[p]:
            continue
        i = (n - 1) // p
        while i > 0:
            res[i * p] = op(inv(res[i]), res[i * p])
            prime_table[i * p] = 0
            i -= 1
    return res


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    # print(*gcd_convolve(nums1, nums2))
    print(*lcm_convolve(nums1, nums2))
