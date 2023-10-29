# 3_树上所有路径的位运算异或的异或和
# 这里的路径至少有两个点
# !方法是考虑`每个点出现在多少条路径上`，若数目为奇数则对答案有贡献
# 路径分两种情况，一种是没有父节点参与的，树形 DP 一下就行了；另一种是父节点参与的，个数就是 子树*(n-子树)

from typing import List, Tuple


def countPath(n: int, adjList: List[List[int]]) -> List[int]:
    """求包含每个点的路径数.这里的路径至少有两个点."""
    res = [0] * n

    def dfs(cur: int, pre: int) -> int:
        count = 0
        size = 1
        for next_ in adjList[cur]:
            if next_ != pre:
                subSize = dfs(next_, cur)
                count += size * subSize
                size += subSize
        count += size * (n - size)
        res[cur] = count
        return size

    dfs(0, -1)
    return res


def xorPathXorSum(n: int, edges: List[Tuple[int, int]], values: List[int]) -> int:
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    pathCount = countPath(n, adjList)
    res = 0
    for i, c in enumerate(pathCount):
        if c & 1:
            res ^= values[i]
    return res


if __name__ == "__main__":
    # https://ac.nowcoder.com/acm/contest/272/B
    import sys

    sys.setrecursionlimit(int(1e6))
    input = sys.stdin.readline
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    values = list(map(int, input().split()))
    print(xorPathXorSum(n, edges, values))
