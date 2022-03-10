# 给出一个`有向无环`的连通图，起点为 1，终点为 N，每条边都有一个长度。
# 数据保证从起点出发能够到达图中所有的点，图中所有的点也都能够到达终点。
# 绿豆蛙从起点出发，走向终点。
# 到达每一个顶点时，如果有 K 条离开该点的道路，绿豆蛙可以选择任意一条道路离开该点，并且走向每条路的概率为 1/K。
# 现在绿豆蛙想知道，从起点走到终点所经过的路径总长度的期望是多少？
# 输出从起点到终点路径总长度的期望值，结果四舍五入保留两位小数。
from collections import defaultdict
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))


@lru_cache(None)
def dfs(cur: int) -> float:
    """从cur到达终点的路径总长度的期望"""
    if cur == n:
        return 0
    res = 0
    for next in adjMap[cur]:
        res += (dfs(next) + adjMap[cur][next]) / len(adjMap[cur])
    return res


n, m = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(int))
for _ in range(m):
    u, v, w = map(int, input().split())
    adjMap[u][v] = w


res = dfs(1)
print('{:.2f}'.format(res))

# 此题也可拓扑序dp
