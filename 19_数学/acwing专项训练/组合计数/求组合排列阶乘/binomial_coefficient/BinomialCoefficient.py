# Lucas 定理用于求解大组合数取模的问题，但模数必须是质数。
# !扩展 Lucas 定理用于求解大组合数取模的问题，模数不一定要求是质数。
#
# P4720 【模板】扩展卢卡斯定理(EXLucas)
# https://www.luogu.com.cn/problem/P4720
# 1≤k≤n≤1e18 ，2≤mod≤1e6 ，不保证 mod 是质数。
#
# https://judge.yosupo.jp/problem/binomial_coefficient


class BinomialCoefficient:
    def __init__(self, mod: int):
        self.MOD = mod
        self.factorization = self._factorize(mod)
        self.facs = []
        self.invs = []
        self.coeffs = []
        self.pows = []
        for p, pe in self.factorization:
            fac = [1] * pe
            for i in range(1, pe):
                fac[i] = fac[i - 1] * (i if i % p else 1) % pe
            inv = [1] * pe
            inv[-1] = fac[-1]
            for i in range(1, pe)[::-1]:
                inv[i - 1] = inv[i] * (i if i % p else 1) % pe
            self.facs.append(fac)
            self.invs.append(inv)
            # coeffs
            c = self._modinv(mod // pe, pe)
            self.coeffs.append(mod // pe * c % mod)
            # pows
            powp = [1]
            while powp[-1] * p != pe:
                powp.append(powp[-1] * p)
            self.pows.append(powp)

    def __call__(self, n: int, k: int):
        if k < 0 or k > n:
            return 0
        if k == 0 or k == n:
            return 1 % self.MOD
        res = 0
        for i, (p, pe) in enumerate(self.factorization):
            res += (
                self._choose_pe(n, k, p, pe, self.facs[i], self.invs[i], self.pows[i])
                * self.coeffs[i]
            )
            res %= self.MOD
        return res

    def _E(self, n, k, r, p):
        res = 0
        while n:
            n //= p
            k //= p
            r //= p
            res += n - k - r
        return res

    def _choose_pe(self, n, k, p, pe, fac, inv, powp):
        r = n - k
        e0 = self._E(n, k, r, p)
        if e0 >= len(powp):
            return 0
        res = powp[e0]
        if (p != 2 or pe == 4) and self._E(n // (pe // p), k // (pe // p), r // (pe // p), p) % 2:
            res = pe - res
        while n:
            res = res * fac[n % pe] % pe * inv[k % pe] % pe * inv[r % pe] % pe
            n //= p
            k //= p
            r //= p
        return res

    def _factorize(self, N):
        factorization = []
        for i in range(2, N + 1):
            if i * i > N:
                break
            if N % i:
                continue
            c = 0
            while N % i == 0:
                N //= i
                c += 1
            factorization.append((i, i**c))
        if N != 1:
            factorization.append((N, N))
        return factorization

    def _modinv(self, a, MOD):
        r0, r1, s0, s1 = a, MOD, 1, 0
        while r1:
            r0, r1, s0, s1 = r1, r0 % r1, s1, s0 - r0 // r1 * s1
        return s0 % MOD


if __name__ == "__main__":
    # https://www.luogu.com.cn/problem/P4720
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, k, mod = map(int, input().split())
    C = BinomialCoefficient(mod)
    print(C(n, k))
