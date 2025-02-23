BINOM = [
    [1],
    [1, 1],
    [1, 2, 1],
    [1, 3, 3, 1],
    [1, 4, 1, 4, 1],
]


class Solution:
    def hasSameDigits(self, s: str) -> bool:
        m = len(s)
        n = m - 2

        digits = list(map(int, s))

        def b5(nv: int, kv: int) -> int:
            res = 1
            while nv or kv:
                nd = nv % 5
                kd = kv % 5
                if kd > nd:
                    return 0
                res = (res * BINOM[nd][kd]) % 5
                nv //= 5
                kv //= 5
            return res

        def add10(c2: int, c5: int) -> int:
            if (c5 % 2) == c2:
                return c5
            return c5 + 5

        s1, s2 = 0, 0
        for j in range(n + 1):
            c2 = 1 if (n & j) == j else 0
            c5 = b5(n, j)
            coeff = add10(c2, c5)
            s1 = (s1 + coeff * digits[j]) % 10
            s2 = (s2 + coeff * digits[j + 1]) % 10
        return s1 == s2
