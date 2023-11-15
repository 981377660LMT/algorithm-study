# F - Intervals on Tree-森林的连通分量个数
# !森林中树的数量 = 总顶点数 − 总边数
# https://blog.csdn.net/Ratina/article/details/107163248


# https://atcoder.jp/contests/abc173/tasks/abc173_f
# 输入 n (2≤n≤2e5) 和一棵树的 n-1 条边（节点编号从 1 开始）。
# 定义 f(L,R) 表示用节点编号在 [L,R] 内的点组成的连通块的个数（边的两个端点必须都在 [L,R] 内）。
# 输出满足 1≤L≤R≤n 的所有 f(L,R) 的和。
# !总顶点数为每个顶点在所有区间出现的次数(左*右)
# !总边数可计算每条边的贡献
from typing import List


def intervalsOnTree(n: int, edges: List[List[int]]) -> int:
    vertexCount, edgeCount = 0, 0
    for i in range(n):
        vertexCount += (i + 1) * (n - i)
    for edge in edges:
        min_, max_ = edge
        if min_ > max_:
            min_, max_ = max_, min_
        edgeCount += (min_ + 1) * (n - max_)
    return vertexCount - edgeCount


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append([u, v])
    print(intervalsOnTree(n, edges))
