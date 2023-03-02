# 求组合数2
# C(n,k) 模 mod , n和k都很大
# !0<=k<=n<=1e18


class BinomialCoefficient:
    # https://github.com/strangerxxxx/kyopro

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

    def __call__(self, n, k):
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


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    T, MOD = map(int, input().split())
    C = BinomialCoefficient(MOD)
    for _ in range(T):
        n, k = map(int, input().split())
        print(C(n, k))
