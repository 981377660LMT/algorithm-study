from typing import List, Tuple

# 无向图的双连通分量 Biconnected Components (BCC)          也叫重连通图
# v-BCC：任意割点都是至少两个不同 v-BCC 的公共点              广义圆方树
# https://oi-wiki.org/graph/bcc/
# https://www.csie.ntu.edu.tw/~hsinmu/courses/_media/dsa_13spring/horowitz_306_311_biconnected.pdf
# 好题 https://codeforces.com/problemset/problem/962/F
# https://leetcode-cn.com/problems/s5kipK/
# 结合树链剖分 https://codeforces.com/problemset/problem/487/E
# !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739


def findVBCC(n: int, graph: List[List[int]]) -> Tuple[List[List[int]], List[int], List[bool]]:
    """
    !Tarjan 算法求无向图的 v-BCC

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        Tuple[List[List[int]], List[int], List[bool]]:
        每个 v-BCC 组里包含哪些点，每个点所在 v-BCC 的编号(从0开始)，每个顶点是否为割点(便于缩点成树)

    Notes:
        - 原图的割点`至少`在两个不同的 v-BCC 中
        - 原图不是割点的点都`只存在`于一个 v-BCC 中
        - v-BCC 形成的子图内没有割点
    """

    def dfs(cur: int, pre: int) -> int:
        nonlocal dfsId, idCount
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        childCount = 0
        for _ei, next in enumerate(graph[cur]):
            # edge = (cur, next, ei)
            edge = (cur, next)
            if dfsOrder[next] == 0:
                stack.append(edge)
                childCount += 1
                nextLow = dfs(next, cur)
                if nextLow >= dfsOrder[cur]:
                    isCut[cur] = True
                    idCount += 1
                    group = []
                    # eids = []
                    while True:
                        topEdge = stack.pop()
                        v1, v2 = topEdge[0], topEdge[1]
                        if vbccId[v1] != idCount:
                            vbccId[v1] = idCount
                            group.append(v1)
                        if vbccId[v2] != idCount:
                            vbccId[v2] = idCount
                            group.append(v2)
                        # eids.append(topEdge[2])
                        if v1 == cur and v2 == next:
                            break
                    # 点数和边数相同，说明该 v-BCC 是一个简单环，且环上所有的边只属于一个简单环
                    # if len(comp) == len(eids):
                    #     for eid in eids:
                    #         onSimpleCycle[eid] = True
                    groups.append(group)
                if nextLow < curLow:
                    curLow = nextLow
            elif next != pre and dfsOrder[next] < dfsOrder[cur]:
                stack.append(edge)
                if dfsOrder[next] < curLow:
                    curLow = dfsOrder[next]
        if pre == -1 and childCount == 1:
            isCut[cur] = False
        return curLow

    dfsId = 0
    dfsOrder = [0] * n
    vbccId = [0] * n
    idCount = 0
    isCut = [False] * n
    stack = []  # (u, v, eid)
    groups = []

    for i, order in enumerate(dfsOrder):
        if order == 0:
            if len(graph[i]) == 0:  # 零度，即孤立点（isolated vertex）
                idCount += 1
                vbccId[i] = idCount
                groups.append([i])
                continue
            dfs(i, -1)

    return groups, [v - 1 for v in vbccId], isCut


def toTree(graph: List[List[int]], groups: List[List[int]], isCut: List[bool]) -> List[List[int]]:
    """
    # !v-BCC 缩点成树

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        List[List[int]]: 缩点成树后的邻接表

    ## !BCC 和割点作为新图中的节点，并在每个割点与包含它的所有 BCC 之间连边
    ## !bcc1 - 割点1 - bcc2 - 割点2 - ...
    """
    n = len(graph)
    idCount = len(groups)
    cutId = [0] * n
    for i, v in enumerate(isCut):
        if v:
            cutId[i] = idCount
            idCount += 1

    tree = [[] for _ in range(idCount)]
    for cur, group in enumerate(groups):
        for g in group:
            if isCut[g]:
                next = cutId[g]
                tree[next].append(cur)
                tree[cur].append(next)

    return tree


if __name__ == "__main__":
    assert findVBCC(5, [[1, 2], [0, 2], [0, 1, 3, 4], [2], [2]]) == (
        [[2, 3], [2, 4], [2, 0, 1]],
        [2, 2, 2, 0, 1],
        [False, False, True, False, False],
    )

    # https://leetcode.cn/circle/discuss/1AjM9B/
    edges = [(0, 1), (1, 2), (0, 5), (1, 4), (2, 3), (3, 4), (4, 5)]
    graph = [[] for _ in range(6)]
    for u, v in edges:
        graph[u].append(v)
        graph[v].append(u)
    print(findVBCC(6, graph=graph))
