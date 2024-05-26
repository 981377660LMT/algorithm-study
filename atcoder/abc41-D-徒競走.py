# D - 徒競走(竞走比赛)
# https://atcoder.jp/contests/abc041/tasks/abc041_d
# 给定n个人，m个关系，求有多少种排列方式使得所有关系都满足
# 每个关系形如(a, b)，表示a在b前面
# n<=16

# !按照顺序安排位置即可

from functools import lru_cache


N, M = map(int, input().split())
edges = []
for _ in range(M):
    a, b = map(int, input().split())
    edges.append((a - 1, b - 1))

mask = [0] * N
for a, b in edges:
    mask[b] |= 1 << a


@lru_cache(None)
def dfs(remain: int) -> int:
    if remain == 0:
        return 1
    res = 0
    for i in range(N):
        bit = 1 << i
        if (remain & bit) and not (remain & mask[i]):
            res += dfs(remain ^ bit)
    return res


target = (1 << N) - 1
res = dfs(target)
dfs.cache_clear()
print(res)
