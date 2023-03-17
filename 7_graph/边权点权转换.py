# !图的边权/点权转换
# 边权转点权: 是在每条边上插入一个点, 新点的编号为 n 到 n+len(edges)
# 点权转边权: 将每个点拆成入点和出点, 入点到出点的边权为原点的点权


from typing import List, Tuple


def convertGraph1(
    n: int, edges: List[Tuple[int, int, int]]
) -> Tuple[List[Tuple[int, int]], List[int]]:
    """将图的边权转换为点权, 返回新图的边和各个新点的点权.
    在每条边上插一个点, 新点的编号为 n 到 n+len(edges).
    """
    values = [e[-1] for e in edges]
    newEdges = []
    for u, v, _ in edges:
        newEdges.append((u, n))
        newEdges.append((n, v))
        n += 1
    return newEdges, values


def convertGraph2(
    edges: List[Tuple[int, int]], values: List[int]
) -> Tuple[int, List[Tuple[int, int, int]]]:
    """将图的点权转换为边权, 返回新图的点数和边.
    将每个点拆成入点和出点, 入点到出点的边权为原点的点权.
    入点的编号为 0 到 n-1, 出点的编号为 n 到 2n-1.
    """
    n = len(values)
    newEdges = []
    for u, v in edges:
        newEdges.append((u, n + u, values[u]))
        newEdges.append((n + u, v, 0))
        newEdges.append((v, n + v, values[v]))
    return 2 * n, newEdges


# cycle
print(convertGraph1(5, [(0, 1, 1), (0, 2, 2), (0, 3, 3), (0, 4, 4)]))
print(convertGraph1(6, [(0, 1, 1), (1, 2, 2), (2, 3, 3), (3, 4, 4), (4, 0, 5)]))
print(convertGraph2([(0, 1), (0, 2), (0, 3), (0, 4)], [1, 2, 3, 4, 5]))
