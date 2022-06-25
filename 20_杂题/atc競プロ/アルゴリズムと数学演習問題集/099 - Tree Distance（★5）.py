# 求树中所有点对的距离之和
# n<=1e5
# !计算每条边的贡献 上面连了几个点*下面连了几个点

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline


n = int(input())
root = 1
adjMap = defaultdict(set)
for _ in range(n - 1):
    u, v = map(int, input().split())
    adjMap[u].add(v)
    adjMap[v].add(u)
sub = defaultdict(int)


def dfs(cur: int, pre: int) -> None:
    sub[cur] = 1
    for next in adjMap[cur]:
        if next == pre:
            continue
        dfs(next, cur)
        sub[cur] += sub[next]


dfs(root, -1)

res = 0
for i in range(1, n + 1):
    res += (sub[i]) * (n - sub[i])
print(res)
