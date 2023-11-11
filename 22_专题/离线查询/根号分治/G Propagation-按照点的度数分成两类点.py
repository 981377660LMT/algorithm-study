"""传播 Propagation

开始每个点i都有自己的颜色i(1<=i<=n)
给你q个操作,每次操作将点i的颜色扩散给相邻点。
问你q个操作之后每个点的颜色是什么。
"""

# 22_专题/离线查询/083 - Colorful Graph（★6）-分块.py

# !按照点的度数分成两类点 针对不同大小根号分治
# !度小的点实时暴力更新 度大的点打上标记延迟更新
# !n,m<=2e5


from math import sqrt
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, q = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)

    queries = [int(x) - 1 for x in map(int, input().split())]

    SQRT = int(sqrt(m))
    bigNexts = [[] for _ in range(n)]  # !预处理出所有点的邻接点中度大于SQRT的大顶点
    for cur in range(n):
        for next in adjList[cur]:
            if len(adjList[next]) >= SQRT:
                bigNexts[cur].append(next)

    colors = [i + 1 for i in range(n)]  # !每个点的颜色
    lasts = [-1] * n  # !每个点最后一次更新的时间
    history = []  # !每次更新的颜色
    for i in range(q):
        node = queries[i]
        curColor = -1

        # !查询当前点的颜色
        # 大顶点实时查询
        if len(adjList[node]) >= SQRT:
            curColor = colors[node]
        else:
            # 小顶点邻居查询
            preI = lasts[node]
            for next in adjList[node]:
                preI = max(preI, lasts[next])
            if preI == -1:  # 没有被更新过
                curColor = colors[node]
            else:
                curColor = history[preI]

        # !大顶点实时更新
        for big in bigNexts[node]:
            colors[big] = curColor

        lasts[node] = i
        history.append(curColor)

    # !查询每个点最后的颜色
    res = []
    for i in range(n):
        if len(adjList[i]) >= SQRT:
            res.append(colors[i])
        else:
            preI = lasts[i]
            for next in adjList[i]:
                preI = max(preI, lasts[next])
            if preI == -1:
                res.append(colors[i])
            else:
                res.append(history[preI])
    print(*res)
