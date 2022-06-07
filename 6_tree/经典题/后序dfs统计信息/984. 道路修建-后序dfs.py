# n  个国家修建恰好 n–1 条双向道路 国家从 1 到 n 编号
# 每条道路的修建费用等于道路长度乘以道路两端的国家个数之差的绝对值。
# 求修建所有道路所需要的总费用
from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))


def dfs(cur: int, pre: int) -> None:
    subTreeCount[cur] = 1
    for next in adjMap[cur]:
        if next == pre:
            continue
        dfs(next, cur)
        subTreeCount[cur] += subTreeCount[next]


def getRes(cur: int, pre: int) -> None:
    global res
    for next in adjMap[cur]:
        if next == pre:
            continue
        diff = abs(subTreeCount[next] - (n - subTreeCount[next]))
        res += diff * adjMap[cur][next]
        getRes(next, cur)


n = int(input())
adjMap = defaultdict(lambda: defaultdict(int))
for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u][v] = w
    adjMap[v][u] = w

subTreeCount = [0] * (n + 1)
dfs(1, -1)
res = 0
getRes(1, -1)
print(res)
