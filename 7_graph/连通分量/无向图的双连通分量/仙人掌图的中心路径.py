# http://kmyk.github.io/competitive-programming-library/graph/catapillar_graph.hpp
# get a central path of a catapillar graph
# 仙人掌图的中心路径
# https://upload.wikimedia.org/wikipedia/commons/b/b7/Caterpillar_tree.svg

# 仙人掌图:https://blog.csdn.net/CreationAugust/article/details/48007069
# 不含自环的,一条边最多属于一个简单环的无向连通图.


from typing import List, Optional


def getCentralPathOfCatapillarGraph(g: List[List[int]]) -> Optional[List[int]]:
    n = len(g)

    # construct the tree with non-leaf vertices
    m = 0
    h = [[] for _ in range(n)]
    for i in range(n):
        if len(g[i]) != 1:
            m += 1
            for j in g[i]:
                if len(g[j]) != 1:
                    h[i].append(j)

    # the tree must be a path graph
    for i in range(n):
        if len(h[i]) >= 3:
            return

    # reconstruct the path
    if m == 0:
        if n == 1:
            return [0]
        if n == 2:
            return [0, 1]
        raise AssertionError("Impossible")

    assert n >= 3
    path = []
    i = 0
    while len(g[i]) == 1 or len(h[i]) == 2:
        i += 1
    for j in g[i]:
        if len(g[i]) == 1:
            path.append(j)
            break
    parent = path[-1]
    while True:
        path.append(i)
        found = False
        for j in h[i]:
            if j != parent:
                parent = i
                i = j
                found = True
                break
        if not found:
            break
    for j in g[i]:
        if len(g[i]) == 1 and j != parent:
            path.append(j)
            break
    return path
