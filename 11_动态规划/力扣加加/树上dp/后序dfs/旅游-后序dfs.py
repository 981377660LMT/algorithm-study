# Cwbc和XHRlyb生活在 s 市，这天他们打算一起出去旅游。
# 旅行地图上有 n 个城市，它们之间通过 n-1 条道路联通。
# Cwbc和XHRlyb第一天会在 s 市住宿，并游览与它距离不超过 1 的所有城市，之后的每天会选择一个城市住宿，然后游览与它距离不超过 1 的所有城市。
# 他们不想住在一个已经浏览过的城市，又想尽可能多的延长旅行时间。
# XHRlyb想知道她与Cwbc最多能度过多少天的时光呢？

# 1 ≤ n ≤ 500000

# 树,每次访问相邻节点
from collections import defaultdict


n, start = list(map(int, input().split()))

adjMap = defaultdict(list)
for _ in range(n - 1):
    u, v = list(map(int, input().split()))
    adjMap[u].append(v)
    adjMap[v].append(u)

# 每个点,访问还是不访问
dp = [[0, 1] for _ in range(n + 1)]


def dfs(cur: int, pre: int) -> None:
    for next in adjMap[cur]:
        if next == pre:
            continue
        dfs(next, cur)
        dp[cur][1] += dp[next][0]
        dp[cur][0] += max(dp[next][0], dp[next][1])


dfs(start, -1)
print(dp[start][1])
