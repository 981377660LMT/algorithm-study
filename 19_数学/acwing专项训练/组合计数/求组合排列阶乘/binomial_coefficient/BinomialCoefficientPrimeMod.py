# 模为素数的组合数计算


class BinomialCoefficientPrimeMod:
    __slots__ = ("fact", "fact_inv", "mod")

    def __init__(self, max_: int, mod: int):
        self.fact = [1] * (max_ + 1)
        self.mod = mod
        tmp = 1
        for i in range(max_):
            tmp *= i + 1
            tmp %= mod
            self.fact[i + 1] = tmp
        self.fact_inv = [1] * (max_ + 1)
        tmp = pow(self.fact[-1], -1, mod)
        self.fact_inv[-1] = tmp
        for i in range(max_ - 1, -1, -1):
            tmp *= i + 1
            tmp %= mod
            self.fact_inv[i] = tmp

    def binomial(self, n: int, r: int):
        if n < 0 or r < 0 or r > n:
            return 0
        return self.fact[n] * self.fact_inv[r] * self.fact_inv[n - r] % self.mod

    def permutation(self, n: int, r: int):
        return self.fact[n] * self.fact_inv[n - r] % self.mod

    def factorial(self, n: int):
        return self.fact[n]

    def factorial_inv(self, n: int):
        return self.fact_inv[n]


if __name__ == "__main__":
    # https://judge.yosupo.jp/problem/binomial_coefficient_prime_mod

    import sys

    input = lambda: sys.stdin.buffer.readline().rstrip()
    T, m = map(int, input().split())
    comb = BinomialCoefficientPrimeMod(min(m, 10**7) - 1, m)
    res = []
    for _ in range(T):
        n, k = map(int, input().split())
        res.append(comb.binomial(n, k))
    print(*res, sep="\n")
