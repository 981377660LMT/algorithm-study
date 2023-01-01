# 求欧拉函数的前缀和模998244353


def eulerPreSum(upper: int) -> int:
    K = int((upper - 1) ** (2 / 3)) + 1
    L = int((upper - 1) ** (1 / 3)) + 1

    # 1. calculate \phi for all i in [1, K]
    # 1.1 prime sieve
    lpf = [*range(K + 1)]
    primes = []
    for i in range(2, K + 1):
        if lpf[i] == i:
            primes.append(i)
        for p in primes:
            if i * p > K or lpf[i] < p:
                break
            lpf[i * p] = p
        p = lpf[i]
        if lpf[i // p] % p == 0:
            lpf[i] = lpf[i // p] * p

    # 1.2 calculate \phi (s.t. \zeta \ast \phi = \zeta_1)
    phi = [*range(K + 1)]
    Phi = [0] * (K + 1)

    # 1.2.1 calculate p-part of \phi
    for p in primes:
        q = p
        while q <= K:
            i = 1
            while i < q:
                phi[q] -= phi[i]
                i *= p
            q *= p

    # 1.2.2 restore \phi & calculate \Phi
    for i in range(1, K + 1):
        phi[i] = phi[lpf[i]] * phi[i // lpf[i]]
        Phi[i] = phi[i] + Phi[i - 1]

    # 2. calculate \Phi for all n = floor(N / i) where i in [1, L]
    dp = [0] * (L + 1)
    for i in range(L, 0, -1):
        n = upper // i
        m = int(n ** (1 / 2))
        res = n * (n + 1) // 2
        for j in range(1, m + 1):
            if j > 1:
                res -= Phi[n // j] if n // j <= K else dp[i * j]
            res -= (n // j - m) * phi[j]
        dp[i] = res

    ans = dp[1]
    return ans % MOD


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    print(eulerPreSum(n))
