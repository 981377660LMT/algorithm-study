#  e-BCC：删除无向图中所有的割边后，剩下的每一个 CC 都是 e-BCC
#  缩点后形成一颗 bridge tree
#  模板题 https://codeforces.com/problemset/problem/1000/E
#  较为综合的一道题 http://codeforces.com/problemset/problem/732/F
# !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739
from typing import List, Tuple


def findEBCC(n: int, edges: List[Tuple[int, int]]) -> Tuple[List[List[int]], List[int]]:
    """
    !Tarjan 算法求无向图的 e-BCC

    Args:
        n (int): 顶点数
        edges (List[Tuple[int, int]]): 边的列表，每个元素为 (v1, v2)

    Returns:
        Tuple[List[List[int]], List[int]]:
        每个 e-BCC 组里包含哪些点，每个点所在 e-BCC 的编号(从0开始)

    Notes:
        - e-BCC 形成的子图中没有桥
        - e-BCC 内的任意两点之间都至少有两条不同的路径相连
    """

    def findBridges(n: int, edges: List[Tuple[int, int]]) -> List[bool]:
        """Tarjan 算法求无向图的割边(桥)

        Args:
            n (int): 顶点数
            edges (List[Tuple[int, int]]): 边的列表，每个元素为 (v1, v2)

        Returns:
            List[bool]: 每条边是否是桥
        """

        def dfs(cur: int, preId: int) -> int:  # 使用 preId 而不是 pre，可以兼容重边的情况
            nonlocal dfsId
            dfsId += 1
            dfsOrder[cur] = dfsId
            curLow = dfsId
            for next, ei in graph[cur]:
                if dfsOrder[next] == 0:
                    nextLow = dfs(next, ei)
                    if nextLow > dfsOrder[cur]:
                        isBridge[ei] = True
                    if nextLow < curLow:
                        curLow = nextLow
                elif ei != preId:
                    if dfsOrder[next] < curLow:
                        curLow = dfsOrder[next]
            return curLow

        m = len(edges)
        isBridge = [False] * m
        dfsId = 0
        dfsOrder = [0] * n  # 值从 1 开始

        for i, order in enumerate(dfsOrder):
            if order == 0:
                dfs(i, -1)

        return isBridge

    graph = [[] for _ in range(n)]  # (next,weight)
    for ei, (u, v) in enumerate(edges):
        graph[u].append((v, ei))
        graph[v].append((u, ei))

    def dfs2(cur: int) -> None:
        nonlocal idCount
        ebccId[cur] = idCount
        group.append(cur)
        for next, ei in graph[cur]:
            if not isBridge[ei] and ebccId[next] == 0:
                dfs2(next)

    isBridge = findBridges(n, edges)
    # 求出原图中每个点的 e-BCC ID
    ebccId = [0] * n
    idCount = -1
    groups = []

    for i, order in enumerate(ebccId):
        if order == 0:
            idCount += 1
            group = []
            dfs2(i)
            groups.append(group)

    return groups, ebccId


def toTree(n: int, edges: List[Tuple[int, int]]) -> List[List[int]]:
    """
    # !e-BCC 缩点成树

    Args:
        n (int): 顶点数
        edges (List[Tuple[int, int]]): 边的列表，每个元素为 (v1, v2)

    Returns:
        List[List[int]]: 树的邻接表

    ## !各个 e-BCC 的连接成一棵树，每个 e-BCC 为树的一个节点
    ## !ebcc1 - ebcc2 - ebcc3 - ...
    """

    groups, ebccId = findEBCC(n, edges)
    idCount = len(groups)
    tree = [[] for _ in range(idCount)]

    # 遍历 edges，若两端点的 bccId 不同（割边）则建边
    # 也可以遍历 isBridge，割边两端点 bccIDs 一定不同
    for a, b in edges:
        u, v = ebccId[a], ebccId[b]
        if u != v:
            tree[u].append(v)
            tree[v].append(u)

    return tree


if __name__ == "__main__":

    assert toTree(5, [(0, 1), (1, 2), (2, 3), (3, 4)]) == [[1], [0, 2], [1, 3], [2, 4], [3]]
    assert findEBCC(5, [(0, 1), (1, 2), (2, 3), (3, 4)]) == (
        [[0], [1], [2], [3], [4]],
        [0, 1, 2, 3, 4],
    )
