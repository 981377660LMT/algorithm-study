# + 给出一颗n个点的树（1e5），每条边有一个权值c，
# + 每次操作可以令某个权值+1，
# + 求最少操作次数令根节点到每个叶节点路径上的权值和相等

# 思路（构造）：
# + 容易发现，越靠近根节点的，调整代价越小。
# + 我们可以把节点深度类比成距离，题目即为求把所有叶子节点调整到同一高度，每次优先调整靠近根部的。
# + 先dfs一遍更新到最远叶节点的距离dis[x]，再循环一遍更新调整其余子节点的距离跟它一样。
from collections import defaultdict

# 构造 时态同步
def dfs(cur: int, parent: int) -> None:
    global res
    for next, weight in adjMap[cur].items():
        if next == parent:
            continue
        dfs(next, cur)
        # 当前结点到最远叶子的距离
        maxDist[cur] = max(maxDist[cur], maxDist[next] + weight)

    for next, weight in adjMap[cur].items():
        if next == parent:
            continue
        res += maxDist[cur] - (maxDist[next] + weight)


n = int(input())
root = int(input())
adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u][v] = w
    adjMap[v][u] = w

maxDist = [0] * (n + 1)
res = 0
dfs(root, -1)


print(res)

