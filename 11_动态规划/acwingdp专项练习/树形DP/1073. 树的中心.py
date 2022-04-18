# 给定一棵树，树中包含 n 个结点（编号1~n）和 n−1 条无向边，每条边都有一个权值。
# 请你在树中找到一个点，使得该点到树中其他结点的最远距离最近。
# 输出一个整数，表示所求点到树中其他结点的最远距离。
# 310. 最小高度树-换根dp

# 树的中心
from collections import defaultdict

n = int(input())
adjMap = defaultdict(set)
for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u].add((v, w))
    adjMap[v].add((u, w))


# 分别记录向下的最大值和次大值
down1, down2 = [0] * (n + 1), [0] * (n + 1)
# 向下取最大值时必须经过的结点
downMaxNeedRoot = [0] * (n + 1)
# 记录节点向上的最大距离
up = [0] * (n + 1)


# 子结点更新父结点向下的最远距离
def dfs1(cur: int, parent: int) -> int:
    """后序dfs，数组记录一下每个结点往下的最远距离、次远距离，返回每个root处的最大路径长度"""
    for next, weight in adjMap[cur]:
        if next == parent:
            continue
        maxCand = dfs1(next, cur) + weight
        if maxCand >= down1[cur]:
            down2[cur], down1[cur] = down1[cur], maxCand
            downMaxNeedRoot[cur] = next
        elif maxCand >= down2[cur]:
            down2[cur] = maxCand
    return down1[cur]


# 父结点更新子结点向上的最远距离
def dfs2(cur: int, parent: int) -> None:
    """前序dfs，利用父结点来更新子结点"""
    """若最远距离 d1[u] 是经过当前子节点 vv 才得到的，那么就只能退而求其次，
    在 up[u]up[u] 和 d2[u]d2[u] 中取最大值，作为当前子节点 vv 往上走的最远距离 up[v]，反之在 up[u] 和 d1[u]中取最大值"""
    # 每个点处最长和次长
    for next, weight in adjMap[cur]:
        if next == parent:
            continue
        if downMaxNeedRoot[cur] == next:
            # 另一条次长路
            up[next] = max(up[cur], down2[cur]) + weight
        else:
            up[next] = max(up[cur], down1[cur]) + weight
        dfs2(next, cur)


dfs1(1, -1)
dfs2(1, -1)
print(down1, down2, up)
res = min(max(up, down) for down, up in zip(down1[1:], up[1:]))
print(res)
