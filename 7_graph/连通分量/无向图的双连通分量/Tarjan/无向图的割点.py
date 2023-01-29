from typing import List

#  无向图的割点（割顶） cut vertices / articulation points
#  https://codeforces.com/blog/entry/68138
#  https://oi-wiki.org/graph/cut/#_1
#  low(v): 在不经过 v 父亲的前提下能到达的最小的时间戳
#  模板题 https://www.luogu.com.cn/problem/P3388
#  LC928 https://leetcode-cn.com/problems/minimize-malware-spread-ii/
# !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L2739


def findCutVertices(n: int, graph: List[List[int]]) -> List[bool]:
    """Tarjan 算法求无向图的割点

    Args:
        n (int): 顶点数
        graph (List[List[int]]): 邻接表

    Returns:
        List[bool]: 每个点是否是割点
    """

    def dfs(cur: int, pre: int) -> int:
        nonlocal dfsId
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        childCount = 0
        for next in graph[cur]:
            if dfsOrder[next] == 0:
                childCount += 1
                nextLow = dfs(next, cur)
                if nextLow >= dfsOrder[cur]:
                    isCut[cur] = True
                if nextLow < curLow:
                    curLow = nextLow
            elif next != pre and dfsOrder[next] < curLow:
                curLow = dfsOrder[next]
        if pre == -1 and childCount == 1:  # 特判：只有一个儿子的树根，删除后并没有增加连通分量的个数，这种情况下不是割顶
            isCut[cur] = False
        return curLow

    isCut = [False] * n
    dfsOrder = [0] * n  # 值从 1 开始
    dfsId = 0
    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i, -1)

    return isCut
