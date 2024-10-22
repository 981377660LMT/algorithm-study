# TODO: O(nlogn)解法 斯特林数+FFT优化
# https://leetcode.cn/problems/find-the-number-of-possible-ways-for-an-event/solutions/2948582/oxlogxjie-fa-si-te-lin-shu-fftyou-hua-by-vnus/


class 写像十二相:
    """https://qiita.com/drken/items/f2ea4b58b0d21621bd51"""

    __slots__ = ("_fac", "_ifac", "_inv", "_mod")

    def __init__(self, size: int, mod: int) -> None:
        self._mod = mod
        self._fac = [1]
        self._ifac = [1]
        self._inv = [1]
        self._expand(size)

    def query(
        self,
        n: int,
        k: int,
        *,
        isBallDistinct: bool,
        isBoxDistinct: bool,
        atMostOneBallPerBox=False,
        noLimitWithBox=False,
        atLeastOneBallPerBox=False,
    ) -> int:
        """n个球放入k个盒子的方案数.

        Args:
            isBallDistinct (bool): 球是否有区别.
            isBoxDistinct (bool): 盒子是否有区别.
            atMostOneBalPerBox (bool, optional): 每个盒子最多放一个球.
            noLimitWithBox (bool, optional): 每个盒子可以放任意个球.
            atLeastOneBallPerBox (bool, optional): 每个盒子至少放一个球.
        """
        limits = (atMostOneBallPerBox, noLimitWithBox, atLeastOneBallPerBox)
        assert limits.count(True) == 1, "Must have one limit and only one limit with box."
        if isBallDistinct and isBoxDistinct:
            if atMostOneBallPerBox:
                return self._solve1(n, k)
            if noLimitWithBox:
                return self._solve2(n, k)
            if atLeastOneBallPerBox:
                return self._solve3(n, k)
        if not isBallDistinct and isBoxDistinct:
            if atMostOneBallPerBox:
                return self._solve4(n, k)
            if noLimitWithBox:
                return self._solve5(n, k)
            if atLeastOneBallPerBox:
                return self._solve6(n, k)
        if isBallDistinct and not isBoxDistinct:
            if atMostOneBallPerBox:
                return self._solve7(n, k)
            if noLimitWithBox:
                return self._solve8(n, k)
            if atLeastOneBallPerBox:
                return self._solve9(n, k)
        if not isBallDistinct and not isBoxDistinct:
            if atMostOneBallPerBox:
                return self._solve10(n, k)
            if noLimitWithBox:
                return self._solve11(n, k)
            if atLeastOneBallPerBox:
                return self._solve12(n, k)

        raise Exception("Unreachable code.")

    def _solve1(self, n: int, k: int) -> int:
        """有区别的球放入有区别的盒子(每个盒子最多放一个球)."""
        return self.P(n, k)

    def _solve2(self, n: int, k: int) -> int:
        """有区别的球放入有区别的盒子(每个盒子可以放任意个球)."""
        return pow(k, n, self._mod)

    def _solve3(self, n: int, k: int) -> int:
        """有区别的球放入有区别的盒子(每个盒子至少放一个球).
        容斥原理:用总方案数减去不合法的方案数.
        O(k*logn)
        """
        mod = self._mod
        res = 0
        for i in range(k + 1):
            if (k - i) & 1:
                res -= self.C(k, i) * pow(i, n, mod)
            else:
                res += self.C(k, i) * pow(i, n, mod)
            res %= mod
        return res

    def _solve4(self, n: int, k: int) -> int:
        """无区别的球放入有区别的盒子(每个盒子最多放一个球)."""
        return self.C(n, k)

    def _solve5(self, n: int, k: int) -> int:
        """无区别的球放入有区别的盒子(每个盒子可以放任意个球)."""
        return self.C(n + k - 1, n)

    def _solve6(self, n: int, k: int) -> int:
        """无区别的球放入有区别的盒子(每个盒子至少放一个球)."""
        return self.C(n - 1, k - 1)

    def _solve7(self, n: int, k: int) -> int:
        """有区别的球放入无区别的盒子(每个盒子最多放一个球)."""
        return 0 if n > k else 1

    def _solve8(self, n: int, k: int) -> int:
        """有区别的球放入无区别的盒子(每个盒子可以放任意个球).
        贝尔数B(n,k).
        O(min(n,k)*logn).
        """
        return self.bell(n, k)

    def _solve9(self, n: int, k: int) -> int:
        """有区别的球放入无区别的盒子(每个盒子至少放一个球).
        第二类斯特林数S(n,k).
        O(k*logn).
        """
        return self.stirling2(n, k)

    def _solve10(self, n: int, k: int) -> int:
        """无区别的球放入无区别的盒子(每个盒子最多放一个球)."""
        return 0 if n > k else 1

    def _solve11(self, n: int, k: int) -> int:
        """无区别的球放入无区别的盒子(每个盒子可以放任意个球).
        分割数P(n,k).
        """
        return self.partition(n, k)

    def _solve12(self, n: int, k: int) -> int:
        """无区别的球放入无区别的盒子(每个盒子至少放一个球).
        分割数P(n-k,k).
        """
        if n < k:
            return 0
        return self.partition(n - k, k)

    def fac(self, k: int) -> int:
        self._expand(k)
        return self._fac[k]

    def ifac(self, k: int) -> int:
        self._expand(k)
        return self._ifac[k]

    def inv(self, k: int) -> int:
        self._expand(k)
        return self._inv[k]

    def C(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(k) % mod * self.ifac(n - k) % mod

    def P(self, n: int, k: int) -> int:
        if n < 0 or k < 0 or n < k:
            return 0
        mod = self._mod
        return self.fac(n) * self.ifac(n - k) % mod

    def H(self, n: int, k: int) -> int:
        """可重复选取元素的组合数"""
        if n == 0:
            return 1 if k == 0 else 0
        return self.C(n + k - 1, k)

    def put(self, n: int, k: int) -> int:
        """n个相同的球放入k个不同的盒子(盒子可放任意个球)的方案数."""
        return self.C(n + k - 1, n)

    def partition(self, n: int, k: int) -> int:
        """O(n*k)"""
        dp = [[0] * (k + 1) for _ in range(n + 1)]
        dp[0][0] = 1
        for i in range(n + 1):
            for j in range(1, k + 1):
                if i >= j:
                    dp[i][j] = dp[i][j - 1] + dp[i - j][j]
                else:
                    dp[i][j] = dp[i][j - 1]
        return dp[n][k]

    def bell(self, n: int, k: int) -> int:
        """O(min(n,k)*logn)"""
        if k > n:
            k = n
        mod = self._mod
        jsum = [0] * (k + 2)
        for j in range(k + 1):
            add = self.ifac(j)
            if j & 1:
                jsum[j + 1] = (jsum[j] - add) % mod
            else:
                jsum[j + 1] = (jsum[j] + add) % mod
        res = 0
        for i in range(k + 1):
            res += pow(i, n, mod) * self.ifac(i) % MOD * jsum[k - i + 1]
            res %= mod
        return res

    def stirling2(self, n: int, k: int) -> int:
        """O(k*logn)"""
        mod = self._mod
        res = 0
        for i in range(k + 1):
            if (k - i) & 1:
                res -= self.C(k, i) * pow(i, n, mod)
            else:
                res += self.C(k, i) * pow(i, n, mod)
            res %= mod
        return res * self.ifac(k) % mod

    def _expand(self, size: int) -> None:
        if len(self._fac) < size + 1:
            mod = self._mod
            preSize = len(self._fac)
            diff = size + 1 - preSize
            self._fac += [1] * diff
            self._ifac += [1] * diff
            self._inv += [1] * diff
            for i in range(preSize, size + 1):
                self._fac[i] = self._fac[i - 1] * i % mod
            self._ifac[size] = pow(self._fac[size], mod - 2, mod)  # !modInv
            for i in range(size - 1, preSize - 1, -1):
                self._ifac[i] = self._ifac[i + 1] * (i + 1) % mod
            for i in range(preSize, size + 1):
                self._inv[i] = self._ifac[i] * self._fac[i - 1] % mod


MOD = int(1e9 + 7)
E = 写像十二相(int(1e5), MOD)

if __name__ == "__main__":

    class Solution:
        # 3317. 安排活动的方案数
        # https://leetcode.cn/problems/find-the-number-of-possible-ways-for-an-event/description/
        # 一个活动总共有 n 位表演者。每一位表演者会 被安排 到 x 个节目之一，有可能有节目 没有 任何表演者。
        # 所有节目都安排完毕后，评委会给每一个 有表演者的 节目打分，分数是一个 [1, y] 之间的整数。
        # 请你返回 总 的活动方案数。
        def numberOfWays(self, n: int, x: int, y: int) -> int:
            res = 0
            for k in range(1, x + 1):
                v1 = E.query(
                    n, k, isBallDistinct=True, isBoxDistinct=True, atLeastOneBallPerBox=True
                ) * E.C(x, k)
                v2 = pow(y, k, MOD)
                res = (res + v1 * v2) % MOD
            return res

    assert (E.partition(10, 5)) == 30
    assert (E.stirling2(8, 2)) == 127
    assert (E.bell(5, 5)) == 52
    n, k = map(int, input().split())
    print(E.bell(n, k))
