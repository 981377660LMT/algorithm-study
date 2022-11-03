#  pow() 函数的三参数形式在`底数与模数不可约的情况`下允许指数为负值。 随后它会在指数为 -1 时计算底数的模乘逆元，并对其他负指数计算反模的适当幂次


# 想要求出 4258𝑥 + 147𝑦 = 369 的整数解
# 首先应重写为 4258𝑥 ≡ 369 (mod 147) 然后求解：


x = 369 * pow(4258, -1, 147) % 147
y = (4258 * x - 369) // -147
print(x, y)


# 写出 7x+13y+29z=n的解
MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(4e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def A(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return (fac[n] * ifac[n - k]) % MOD


print(C(12345, 123))
