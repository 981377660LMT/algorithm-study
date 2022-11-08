from typing import List, Tuple


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
        Tuple[List[List[int]], List[int],List[bool]]:
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

    #  由于每个强连通分量都是在它的所有后继强连通分量被求出之后求得的
    # 上面得到的 scc 是拓扑序的逆序
    groups.reverse()
    sccId = [0] * n
    for i, group in enumerate(groups):
        for v in group:
            sccId[v] = i

    return groups, sccId


# EXTRA: 缩点: 将边 v-w 转换成 sid[v]-sid[w]
# 缩点后得到了一张 DAG，点的编号范围为 [0,len(scc)-1]
# 模板题 点权 https://www.luogu.com.cn/problem/P3387
#  边权 https://codeforces.com/contest/894/problem/E
# 检测路径是否可达/唯一/无穷 https://codeforces.com/problemset/problem/1547/G
def toTree(n: int, graph: List[List[int]]) -> List[List[int]]:
    """
    # !scc 缩点成DAG

    Args:
        n (int): 图的顶点数
        graph (List[List[int]]):  邻接表

    Returns:
        List[List[int]]: 缩点成DAG后的邻接表
    """

    groups, sccId = findSCC(n, graph)
    m = len(groups)
    adjList = [[] for _ in range(m)]  # !注意这样可能会产生重边，不能有重边时可以用 adjMap 去重
    deg = [0] * m
    for i, nexts in enumerate(graph):
        u = sccId[i]
        for next in nexts:
            v = sccId[next]
            if u != v:
                adjList[u].append(v)
                deg[v] += 1
            else:
                # 这里可以记录自环（指 len(scc) == 1 但是有自环）、汇合同一个 SCC 的权值等 ...
                pass

    return adjList


print(findSCC(5, [[1, 2], [2, 3], [3, 1], [4], []]))
print(toTree(5, [[1, 2], [2, 3], [3, 1], [4], []]))
