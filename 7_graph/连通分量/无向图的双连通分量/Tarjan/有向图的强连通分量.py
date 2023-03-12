from collections import deque
from typing import List, Set, Tuple


#  !SCC Tarjan (Tarjan 求有向图的强联通分量，缩点成拓扑图)
#  常数比 Kosaraju 略小（在 AtCoder 上的测试显示，5e5 的数据下比 Kosaraju 快了约 100ms）
#  https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm
#  https://oi-wiki.org/graph/scc/#tarjan
#  https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/TarjanSCC.java.html
#  https://stackoverflow.com/questions/32750511/does-tarjans-scc-algorithm-give-a-topological-sort-of-the-scc
#  与最小割结合 https://www.luogu.com.cn/problem/P4126
# !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739
def findSCC(n: int, graph: List[List[int]]) -> Tuple[List[List[int]], List[int]]:
    """
    # !Tarjan 算法求有向图的 scc

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        Tuple[List[List[int]], List[int]]:
        每个 scc 组里包含哪些点，每个点所在 scc 的编号(从0开始)
    """

    def dfs(cur: int) -> int:
        nonlocal dfsId
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        stack.append(cur)
        inStack[cur] = True
        for next in graph[cur]:
            if dfsOrder[next] == 0:
                nextLow = dfs(next)
                if nextLow < curLow:
                    curLow = nextLow
            elif inStack[next] and dfsOrder[next] < curLow:
                curLow = dfsOrder[next]

        if dfsOrder[cur] == curLow:
            group = []
            while True:
                top = stack.pop()
                inStack[top] = False
                group.append(top)
                if top == cur:
                    break
            groups.append(group)

        return curLow

    dfsOrder = [0] * n
    dfsId = 0
    stack = []
    inStack = [False] * n
    groups = []

    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i)

    # 由于每个强连通分量都是在它的所有后继强连通分量被求出之后求得的
    # 上面得到的 scc 是拓扑序的逆序
    groups.reverse()
    sccId = [0] * n
    for i, group in enumerate(groups):
        for v in group:
            sccId[v] = i

    return groups, sccId  # !groups按照拓扑序输出


# EXTRA: 缩点: 将边 v-w 转换成 sid[v]-sid[w]
# 缩点后得到了一张 DAG，点的编号范围为 [0,len(scc)-1]
# 模板题 点权 https://www.luogu.com.cn/problem/P3387
#  边权 https://codeforces.com/contest/894/problem/E
# 检测路径是否可达/唯一/无穷 https://codeforces.com/problemset/problem/1547/G
def toDAG(graph: List[List[int]], groups: List[List[int]], sccId: List[int]) -> List[Set[int]]:
    """
    # !scc 缩点成DAG

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        List[List[int]]: 缩点成DAG后的邻接表
    """

    m = len(groups)
    adjList = [set() for _ in range(m)]  # !set去重
    deg = [0] * m
    for i, nexts in enumerate(graph):
        u = sccId[i]
        for next in nexts:
            v = sccId[next]
            if u != v:
                adjList[u].add(v)
                deg[v] += 1
            else:
                # 这里可以记录自环（指 len(scc) == 1 但是有自环）、汇合同一个 SCC 的权值等 ...
                pass

    return adjList


if __name__ == "__main__":
    # assert findSCC(5, [[1, 2], [2, 3], [3, 1], [4], []]) == (
    #     [[0], [2, 1], [3], [4]],
    #     [0, 1, 1, 2, 3],
    # )

    # https://www.luogu.com.cn/problem/P3387
    #   给定一个 n 个点 m 条边有向图，每个点有一个权值
    #   求一条路径，使路径经过的点权值之和最大
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    values = list(map(int, input().split()))
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList[u - 1].append(v - 1)

    groups, sccId = findSCC(n, adjList)
    dag = toDAG(adjList, groups, sccId)
    weights = [0] * len(groups)
    for i, group in enumerate(groups):
        weights[i] = sum(values[v] for v in group)

    deg = [0] * len(groups)
    dp = [0] * len(groups)
    queue = deque()
    for i, nexts in enumerate(dag):
        for next in nexts:
            deg[next] += 1
    for i, d in enumerate(deg):
        if d == 0:
            queue.append(i)
            dp[i] = weights[i]

    while queue:
        cur = queue.popleft()
        for next in dag[cur]:
            dp[next] = max(dp[next], dp[cur] + weights[next])
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)

    print(max(dp))
