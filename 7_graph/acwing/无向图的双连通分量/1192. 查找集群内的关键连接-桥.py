from typing import List, Tuple

# !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739


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

    graph = [[] for _ in range(n)]  # (next,weight)
    for ei, (u, v) in enumerate(edges):
        graph[u].append((v, ei))
        graph[v].append((u, ei))

    m = len(edges)
    isBridge = [False] * m
    dfsId = 0
    dfsOrder = [0] * n  # 值从 1 开始

    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i, -1)

    return isBridge


if __name__ == "__main__":
    print(findBridges(5, [(0, 1), (1, 2), (2, 3), (3, 4)]))

    # https://leetcode.cn/problems/critical-connections-in-a-network/submissions/
    class Solution:
        def criticalConnections(self, n: int, connections: List[List[int]]) -> List[List[int]]:
            isBridge = findBridges(n, connections)  # type: ignore
            return [connections[i] for i, v in enumerate(isBridge) if v]

    assert Solution().criticalConnections(4, [[0, 1], [1, 2], [2, 0], [1, 3]]) == [[1, 3]]
