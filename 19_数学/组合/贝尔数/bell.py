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
    if n == 0:
        return 1
    k = min(k, n)
    res = 0
    pre = [0] * (k + 1)
    for i in range(1, k + 1):
        if i & 1:
            pre[i] = pre[i - 1] - ifac[i]
        else:
            pre[i] = pre[i - 1] + ifac[i]
        pre[i] %= MOD
    for i in range(1, k + 1):
        res += (pow(i, n, MOD) * ifac[i] % MOD) * pre[k - i]
        res %= MOD
    return res


print(bell(5, 1))
