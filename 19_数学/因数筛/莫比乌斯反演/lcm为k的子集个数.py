# F - Subsequence LCM
# https://atcoder.jp/contests/abc349/tasks/abc349_f
# N<=2e5,A[i]<=1e16
#
# 思路：
# 对k质因数分解，得到p1^e1*p2^e2*...*pm^em
# 每个元素可以状态压缩(是否包含pi).
# !问题变成了有多少种选择方案，使得选出来元素对应的二进制状态的按位或等于 2^m-1
# 记得特判 k=1的情况。

from typing import List, Tuple

MOD = 998244353


def countSubsetLcm(nums: List[int], lcm_: int) -> int:
    """lcm为k的子集个数(模mod)."""
    if not nums:
        return 0

    pfs = getPrimeFactors(lcm_)
    m = len(pfs)
    counter = [0] * (1 << m)
    for num in nums:
        if lcm_ % num != 0:
            continue
        mask = 0
        for i, (p, e) in enumerate(pfs):
            if num % (p**e) == 0:
                mask |= 1 << i
        counter[mask] += 1

    # divisorZeta
    for i in range(m):
        for pre in range(1 << m):
            if (pre >> i) & 1 == 0:
                counter[pre | (1 << i)] += counter[pre]

    dp = [pow(2, counter[i], MOD) - 1 for i in range(1 << m)]

    # divisorMobius
    for i in range(m):
        for pre in range(1 << m):
            if (pre >> i) & 1 == 0:
                dp[pre | (1 << i)] -= dp[pre]

    res = dp[-1]
    return res % MOD


def getPrimeFactors(n: int) -> List[Tuple[int, int]]:
    """质因数分解.

    >>> getPrimeFactors(100)
    [(2, 2), (5, 2)]
    """
    res = []
    upper = n
    i = 2
    while i * i <= upper:
        if upper % i == 0:
            c = 0
            while upper % i == 0:
                c += 1
                upper //= i
            res.append((i, c))
        i += 1
    if upper != 1:
        res.append((upper, 1))
    return res


if __name__ == "__main__":
    N, K = map(int, input().split())
    A = list(map(int, input().split()))
    print(countSubsetLcm(A, K))
