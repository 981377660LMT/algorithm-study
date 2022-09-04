# n 和 k (1 <= n <= 2000, 0 <= k <= 2000)，分别表示单词数量和所需的相似对数。
# 有多少种不同的方案数，可以从 n 个单词中选择一个子集（可以是空集），而且该子集中恰好有 k 对相似的单词。
# 输入保证所有给定的单词都是不同的。
import sys
from collections import Counter
from functools import lru_cache
from math import comb

sys.setrecursionlimit(int(1e9))
comb = lru_cache(comb)


@lru_cache(None)
def dfs(index: int, remain: int) -> int:
    if remain < 0:
        return 0
    if index == len(counts):
        return 1 if remain == 0 else 0

    res = 0
    for select in range(counts[index] + 1):
        # 有多少对
        pair = comb(select, 2)
        if pair > remain:
            break
        res += dfs(index + 1, remain - pair) * comb(counts[index], select)
        res %= MOD
    return res


MOD = int(1e9 + 7)

n, k = map(int, input().split())
counter = Counter()

for _ in range(n):
    key = tuple(sorted(input()))
    counter[key] += 1

counts = list(counter.values())

res = dfs(0, k)
dfs.cache_clear()
print(res)

