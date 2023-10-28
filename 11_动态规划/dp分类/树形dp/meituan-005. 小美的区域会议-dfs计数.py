from collections import defaultdict, Counter
from typing import List

MOD = int(1e9 + 7)

# 注意不要重复计数:
# level相等时，只看nextId<rootId
# level不等时，只看nextLevel<rootLevel


def dfs(cur: int, visited: List[bool], root: int) -> int:
    res = 1
    for next in adjMap[cur]:
        if visited[next]:
            continue
        if (
            levels[next] == levels[root]
            and next < root
            or levels[next] < levels[root]
            and levels[root] - levels[next] <= k
        ):
            visited[next] = True
            res *= (1 + dfs(next, visited, root)) % MOD
    return res % MOD


n, k = list(map(int, input().split()))
adjMap = defaultdict(set)

for _ in range(n - 1):
    u, v = list(map(int, input().split()))
    u, v = u - 1, v - 1
    adjMap[u].add(v)
    adjMap[v].add(u)

levels = list(map(int, input().split()))

# 从每个点开始dfs，相邻不超过k
res = 0
for i in range(n):
    res += dfs(i, [False] * n, i) % MOD
print(res % MOD)
