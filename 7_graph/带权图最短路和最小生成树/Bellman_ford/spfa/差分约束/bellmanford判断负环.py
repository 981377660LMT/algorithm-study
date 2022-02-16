# https://oi-wiki.org/graph/diff-constraints/

# 用途:求一组不等式组的可行解
# 差分约束就是将一些不等式转化为图中的带权边，然后求解最短路或最长路的方法
# 用差分约束，>x转化为>=x+1。=转化为>=且<=。

# xi-xj<=ck可变形为 xi<=xj+ck，这与单源最短路中的三角形不等式
# dist[xi]<=dist[xj]+ck非常相似
# 因此我们可以把每个变量看成图的一个节点，
# 对每个约束条件xi-xj<=ck，从结点j向i连接一条长度为ck的有向边

# 一般使用Bellman-ford或者spfa判断`图中是否存在负环`，最坏时间复杂度为O(n*m)


# 求解差分约束系统，有m条约束条件,判断该差分约束系统有没有解
# 如果要求xi/xj<=ck，只需取对数即可
from collections import defaultdict
from typing import List


Edge = List[int]


def bellman_ford(edges: List[Edge], start: int) -> bool:
    """bellman_ford判断负权环"""
    n = len(edges)
    dist = defaultdict(lambda: int(1e20))  # 起点s到各个点的距离
    dist[start] = 0

    # 松弛i次:其中第i次(i>=1)的内涵为此时至少优化过了过了i-1个`中转点`，最后一次优化了n-1个中转点(即所有点都经过了)
    for _ in range(n):
        for u, v, w in edges:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w

    for u, v, w in edges:
        if dist[u] + w < dist[v]:
            return False  # 存在负权边

    return True


if __name__ == '__main__':
    #  x0-x1>=1    x0-1>=x1
    #  x2-x3<=2    x3+2>=x2
    #  x0=x2      x0+0>=x2  x2+0>=x0
    edges = []
    edges.append([0, 1, -1])
    edges.append([3, 2, 2])
    edges.append([0, 2, 0])
    edges.append([2, 0, 0])
    print(bellman_ford(edges, 0))

