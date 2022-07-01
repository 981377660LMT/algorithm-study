# 给你一个数字 n ，你需要使用最少的操作次数，在记事本上输出 恰好 n 个 'A' 。
# 返回能够打印出 n 个 'A' 的最少操作次数。
from collections import Counter
from functools import lru_cache
from math import floor


@lru_cache(None)
def getPrimeFactors(n: int) -> Counter[int]:
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


class Solution:
    def minSteps(self, n: int) -> int:
        """n^3/2"""

        @lru_cache(None)
        def dfs(cur: int, paste: int) -> int:
            """记事本上有 cur 个字符，粘贴板上有 paste 个字符"""
            if cur >= n:
                return 0 if cur == n else int(1e20)
            if paste != 0 and (n - cur) % paste != 0:
                return int(1e20)
            res = int(1e20)
            if paste != 0:
                res = dfs(cur + paste, paste) + 1
            if paste != cur:
                res = min(res, dfs(cur, cur) + 1)
            return res

        return dfs(1, 0)

    # 在达到最小的总代价时，我们将 n 拆分成了若干个素数的乘积

    def minSteps2(self, n: int) -> int:
        """n^1/2"""
        factors = getPrimeFactors(n)
        return sum(v * c for v, c in factors.items())
