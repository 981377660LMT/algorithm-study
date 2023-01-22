# 统计好子序列的数目
# 给定一个仅由小写字母组成的字符串 s ，求 s 的所有非空子序列中好子序列的数目。
# !好子序列:子序列中所有字符的频率相等

from collections import Counter


MOD = int(1e9 + 7)
fac = [1]
ifac = [1]
for i in range(1, int(2e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


class Solution:
    def countGoodSubsequences(self, s: str) -> int:
        n = len(s)
        counter = Counter(s)
        res = 0
        for same in range(1, n + 1):  # 每种字母选same个或者不选
            cur = 1
            for char in counter:
                cur *= C(counter[char], same) + 1
                cur %= MOD
            res += cur - 1  # 减去都不选的情况(空集)
            res %= MOD
        return res


assert Solution().countGoodSubsequences(s="aabb") == 11
