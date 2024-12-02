# AtCoder Beginner Contest 259 - SGColin的文章 - 知乎 (删边、删除树边、树上删边)
# https://zhuanlan.zhihu.com/p/539701972
# https://atcoder.jp/contests/abc259/tasks/abc259_f
#
# !给定一棵树，每条边有一个边权 w ，每个点有一个限制di。
# !选一个边集，使得每个点与至多 di 个其他节点有边直接相连。
# 最大化选中边的权值和。
#
# !设f[i][0/1]表示节点i及其子树内，是否要选i到父亲的边(0/1)，能得到的最大价值。
# 如果选到父亲的边:就是最多把d-1个儿子的贡献从 f[son][0]改为f[son][1] + w[u][son] ;特殊的如果d=0 则f[i][1] = -inf
# 如果不选到父亲的边:就是最多把d;个儿子的贡献从f[son][0]改为f[son][1] + w[u][son]，挑能贡献最多的选（修改后较修改前差值最大的d个)
# !对贡献的差值排序 选边
#
#    parent
#      !|*
#     node
#    /    \
#  child   child
#
#
# # *を選ぶ場合のnode以下の部分木に関する最適値
# dp1 = [0] * n
# # *を選ばない場合のnode以下の部分木に関する最適値
# dp0 = [0] * n


from typing import List, Tuple
from heapq import nlargest

INF = int(1e18)


def selectEdges(n: int, edges: List[Tuple[int, int, int]], limits: List[int]) -> int:
    """
    从树中选出一个边集，使得每个点与至多 di 个其他节点有边直接相连.
    最大化选中边的权值和.

    https://atcoder.jp/contests/abc259/tasks/abc259_f
    """

    def dfs(cur: int, pre: int) -> Tuple[int, int]:
        """
        返回[选择连接父亲的边时子树最大价值, 不选择连接父亲的边时子树最大价值].

        选择父亲,那么子树里最多选limit-1条边;
        不选择父亲,那么子树里最多选limit条边;
        考虑两种策略的差值，排序.
        """
        res1, res2 = 0, 0
        diff = []
        for next, weight in adjList[cur]:
            if next == pre:
                continue
            select, skip = dfs(next, cur)
            res1, res2 = res1 + skip, res2 + skip  # !先默认不选择
            tmp = select + weight - skip
            diff.append(tmp)

        limit = limits[cur]
        if limit == 0:  # !选不了特判
            return -INF, res2

        diff = nlargest(limit, diff)
        for i in range(len(diff)):
            if diff[i] <= 0:  # !选比不选还差就不选
                break
            if i < limit - 1:
                res1 += diff[i]
            if i < limit:
                res2 += diff[i]
        return res1, res2

    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    res = dfs(0, -1)
    return res[1]  # !不选虚拟根节点


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))

    n = int(input())
    limits = list(map(int, input().split()))
    edges = []
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    print(selectEdges(n, edges, limits))
