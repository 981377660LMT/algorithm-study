# 有一棵二叉苹果树，如果树枝有分叉，一定是分两叉，即没有只有一个儿子的节点。
# 这棵树共 N 个节点，编号为 1 至 N，树根编号一定为 1。
# 我们用一根树枝两端连接的节点编号描述一根树枝的位置。
# 一棵苹果树的树枝太多了，需要剪枝。但是一些树枝上长有苹果，给定需要保留的树枝数量，求最多能留住多少苹果。
# 这里的保留是指最终与1号点连通。
# 第一行包含两个整数 N 和 Q，分别表示树的节点数以及要保留的树枝数量。
# 接下来 N−1 行描述树枝信息，每行三个整数，前两个是它连接的节点的编号，第三个数是这根树枝上苹果数量。
# 输出仅一行，表示最多能留住的苹果的数量。
'''
在树上选择Q条边，等价于选择连在一起的Q+1个节点，每个节点的权值就是和其父节点连接的边的权值，
根节点的权值为0，问题就变成了一个有依赖的背包问题，首先根节点必选，如果一个节点被选择了，其
父节点必定被选择

dp(i, j) 表示以i为根的子树上选择j个节点，所有节点的权值最大和是多少
在进行决策时候其实就是两个分组进行背包，物品开销是子树选择节点数jj，物品的价值就是dp(子节点id, jj)

'''
from collections import defaultdict
from functools import lru_cache


n, k = map(int, input().split())

adjMap = defaultdict(list)

for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u].append((v, w))
    adjMap[v].append((u, w))

# 每个节点的节点数和权重
subtreeCounts = [0] * (n + 1)
weights = [0] * (n + 1)


def dfs1(cur: int, pre: int) -> int:
    """统计信息"""
    curCount = 1
    for next, weight in adjMap[cur]:
        if next == pre:
            continue
        weights[next] = weight
        curCount += dfs1(next, cur)

    subtreeCounts[cur] = curCount
    # 注意要删除父节点，因为之后dfs2里要找两个子节点
    adjMap[cur] = [item for item in adjMap[cur] if item[0] != pre]
    return curCount


dfs1(1, -1)


@lru_cache(None)
def dfs2(root: int, select: int) -> int:
    """以root为根的子树上选择select个节点，所有节点的权值最大和是多少"""
    if select == 0:
        return 0
    if select > subtreeCounts[root]:  # impossible
        return -int(1e20)
    if select == 1:
        return weights[root]

    res = 0
    left, right = adjMap[root][0][0], adjMap[root][1][0]
    for leftSelect in range(select):
        rightSelect = select - leftSelect - 1
        res = max(res, dfs2(left, leftSelect) + dfs2(right, rightSelect) + weights[root])
    return res


print(dfs2(1, k + 1))
