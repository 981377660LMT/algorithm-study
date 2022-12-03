# 给定一个n个点m条边的有向图
# !对于边(u,v)权值为w,从u到v的距离为w,从v到u的距离为-w
# 给定q个查询 (start,end) 求从start到end的最大得分
# 如果无法到达输出nan，如果可以到达且可以无限增加输出inf
# 所有数据范围为1e5

# 到根节点距离的并查集(带权并查集)
# !nan:不连通
# !inf:存在矛盾(正环)

from typing import List, Tuple
from UnionFindMapWithDist import UnionFindArrayWithDist

INF = int(1e18)


def payOrReceive(
    n: int, edges: List[Tuple[int, int, int]], queries: List[Tuple[int, int]]
) -> List[int]:
    uf = UnionFindArrayWithDist(n)
    cycleGroup = set()  # 带有正环(距离存在矛盾)的组
    for u, v, w in edges:
        if uf.isConnected(u, v):
            if uf.getDist(u, v) != w:  # 矛盾,组内存在正环
                cycleGroup.add(uf.find(u))
        else:
            root1, root2 = uf.find(u), uf.find(v)
            if root1 in cycleGroup or root2 in cycleGroup:  # 合并组内正环
                cycleGroup |= {root1, root2}
            uf.union(u, v, w)

    res = [0] * len(queries)  # !res[i] => 非正数:最大得分,1:无法到达,2:存在矛盾(正环)
    for i, (u, v) in enumerate(queries):
        if not uf.isConnected(u, v):
            res[i] = 1
            continue
        root = uf.find(u)
        res[i] = 2 if root in cycleGroup else uf.getDist(u, v)
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m, q = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))
    queries = []
    for _ in range(q):
        start, end = map(int, input().split())
        start, end = start - 1, end - 1
        queries.append((start, end))

    res = payOrReceive(n, edges, queries)
    for num in res:
        if num == 1:
            print("nan")
        elif num == 2:
            print("inf")
        else:
            print(num)
