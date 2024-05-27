from typing import Sequence, Tuple


class ArrayMod:
    __slots__ = "_mod", "_fact", "_invFact", "_one", "_n"

    def __init__(self, n: int, mod: int):
        self._fact = [0] * n
        self._invFact = [0] * n
        self._mod = mod
        self._one = 1 % mod

    def init(self, seq: Sequence[int]):
        n = len(seq)
        self._n = n
        if n == 0:
            return
        mod, fact, invFact = self._mod, self._fact, self._invFact
        for i, v in enumerate(seq):
            fact[i] = v
            if fact[i] == 0:
                fact[i] = 1
            if i > 0:
                fact[i] = (fact[i] * fact[i - 1]) % mod
        invFact[n - 1] = modInv(fact[n - 1], mod)
        for i in range(n - 2, -1, -1):
            invFact[i] = (invFact[i + 1] * seq[i + 1]) % mod

    def prefixMod(self, right: int) -> int:
        if right >= self._n:
            right = self._n - 1
        if right < 0:
            return self._one
        return self._fact[right]

    def prefixModInv(self, right: int) -> int:
        if right >= self._n:
            right = self._n - 1
        if right < 0:
            return self._one
        return self._invFact[right]

    def rangeMod(self, left: int, right: int) -> int:
        if left < 0:
            left = 0
        if right >= self._n:
            right = self._n - 1
        if left > right:
            return self._one
        res = self._fact[right]
        if left > 0:
            res = (res * self._invFact[left - 1]) % self._mod
        return res

    def rangeModInv(self, left: int, right: int) -> int:
        if left < 0:
            left = 0
        if right >= self._n:
            right = self._n - 1
        if left > right:
            return self._one
        res = self._invFact[right]
        if left > 0:
            res = (res * self._fact[left - 1]) % self._mod
        return res


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """
    求a, b最大公约数,同时求出裴蜀定理中的一组系数x, y,
    满足 x*a + y*b = gcd(a, b)

    ax + by = gcd_ 返回 `(gcd_, x, y)`

    """
    if b == 0:
        return a, 1, 0
    gcd_, x, y = exgcd(b, a % b)
    return gcd_, y, x - a // b * y


def modInv(a: int, mod: int) -> int:
    """
    扩展gcd求a在mod下的逆元
    即求出逆元 `inv` 满足 `a*inv ≡ 1 (mod m)`
    """
    gcd_, x, _ = exgcd(a, mod)
    if gcd_ not in (1, -1):
        raise ValueError("No inverse")
    return x % mod


if __name__ == "__main__":

    def preixMod(arr, i: int, mod: int) -> int:
        res = 1
        for j in range(i + 1):
            res *= arr[j]
            res %= mod
        return res % mod

    def preixModInv(arr, i: int, mod: int) -> int:
        res = 1
        for j in range(i + 1):
            res *= modInv(arr[j], mod)
            res %= mod
        return res % mod

    def rangeMod(arr, l: int, r: int, mod: int) -> int:
        res = 1
        for i in range(l, r + 1):
            res *= arr[i]
            res %= mod
        return res % mod

    def rangeModInv(arr, l: int, r: int, mod: int) -> int:
        res = 1
        for i in range(l, r + 1):
            res *= modInv(arr[i], mod)
            res %= mod
        return res % mod

    import random

    n = random.randint(1, 100)
    MOD = 1000000007
    nums = [random.randint(1, MOD - 1) for _ in range(n)]
    arr = ArrayMod(n, MOD)
    arr.init(nums)
    for _ in range(100):
        i = random.randint(0, n - 1)
        j = random.randint(0, n - 1)
        if i > j:
            i, j = j, i
        assert arr.prefixMod(i) == preixMod(nums, i, MOD)
        assert arr.prefixModInv(i) == preixModInv(nums, i, MOD)
        assert arr.rangeMod(i, j) == rangeMod(nums, i, j, MOD)
        assert arr.rangeModInv(i, j) == rangeModInv(nums, i, j, MOD)

    print("PASSED")
