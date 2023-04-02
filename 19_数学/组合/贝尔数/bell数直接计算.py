# https://ei1333.github.io/library/math/combinatorics/bell-number.hpp

# 贝尔数B(n,k):
# n个不同的球放到不超过k个相同的盒子里的方案数
# B(n,n)表示将n个球分成任意组的方案数


MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(2e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def bell(n: int, k: int) -> int:
    if k > n:
        k = n
    mod = MOD
    jsum = [0] * (k + 2)
    for j in range(k + 1):
        add = ifac[j]
        if j & 1:
            jsum[j + 1] = (jsum[j] - add) % mod
        else:
            jsum[j + 1] = (jsum[j] + add) % mod
    res = 0
    for i in range(k + 1):
        res += pow(i, n, mod) * ifac[i] % MOD * jsum[k - i + 1]
        res %= mod
    return res


print(bell(5, 1))
