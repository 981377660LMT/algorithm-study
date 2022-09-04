# 给定一个长度不超过5000的字符串，
# 问有多少不同的字符串可以是这个字符串的子序列的一个排列。
# !题意翻译就是，给定每个小写字母有多少个，问能拼出多少个不同的字符串。
# O - 文字列-组合计数dp

# 考虑能构成的字符串只与字符集合有关

import sys
from collections import Counter

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

fac = [1]
ifac = [1]
for i in range(1, 5005):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


if __name__ == "__main__":
    s = input()
    counter = Counter(s)
    freq = list(counter.values())
    n, g = len(s), len(freq)
    memo = [-1] * (26 + 1) * (n + 1)

    def dfs(index: int, count: int) -> int:
        """当前在第index组,已经选了count个元素"""
        if index == g:
            return 1

        hash_ = index * n + count
        if ~memo[hash_]:
            return memo[hash_]

        res = 0
        for select in range(freq[index] + 1):
            res += dfs(index + 1, count + select) * C(select + count, count)
            res %= MOD

        memo[hash_] = res
        return res

    print(dfs(0, 0) - 1)  # 减去1:空字符串。
