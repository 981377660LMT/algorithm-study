# E - Unique Color
# !求树中所有的好点，一个点是好点当且仅当它的颜色在从根到它的路径上只出现了一次。

# 显然从 1 到结点 x 的路径只有一条，
# 我们直接从 1 开始DFS所有点即可，我们可以统计一下每种颜色的出现次数，
# 如果当前搜到的的颜色只出现了一次，那么他就是一个好点。


from collections import defaultdict
from typing import DefaultDict, List, Tuple


def uniqueColor(n: int, edges: List[Tuple[int, int]], values: List[int]) -> List[int]:
    def dfs(cur: int, pre: int, counter: DefaultDict[int, int]) -> None:
        counter[values[cur]] += 1
        if counter[values[cur]] == 1:
            res[cur] = True
        for next in adjList[cur]:
            if next == pre:
                continue
            dfs(next, cur, counter)
        counter[values[cur]] -= 1

    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    res = [False] * n
    dfs(0, -1, defaultdict(int))
    return [i for i in range(n) if res[i]]


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    values = list(map(int, input().split()))
    edges = []
    for i in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))

    res = uniqueColor(n, edges, values)
    print(*[num + 1 for num in res], sep="\n")
