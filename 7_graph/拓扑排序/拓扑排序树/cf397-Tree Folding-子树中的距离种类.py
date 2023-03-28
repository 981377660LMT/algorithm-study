# 给你一棵树，可以把树上父亲相同的两条长度相同的链合并。
# 问你最后能不能变成一条链，
# 能的话求链的最短长度。 （n<2*10^5）
# https://www.luogu.com.cn/problem/CF765E
# https://blog.csdn.net/baidu_35520981/article/details/55261333

# 两种情况:
# 1. 有一个结点子树中距离种类为2,其余结点子树距离种类为1
#    (此时答案为子树种距离种类为2的两种距离的和)
#    !这个结点也就是拓扑排序最后一个结点
# 2. 所有结点子树距离种类为1
#    (此时答案为所有子树距离最大值)
# !注意答案为偶数时,需要除以2直到奇数为止

from collections import deque
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    edges = []
    adjList = [[] for _ in range(n)]
    deg = [0] * n
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
        adjList[u].append(v)
        adjList[v].append(u)
        deg[u] += 1
        deg[v] += 1

    queue = deque()
    dp = [set() for _ in range(n)]  # dp[i]: 子树i中距离种类
    for i in range(n):
        if deg[i] == 1:
            queue.append(i)
            dp[i].add(0)

    visited = [False] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next_ in adjList[cur]:
            if visited[next_]:
                continue
            deg[next_] -= 1
            dp[next_] |= {pre + 1 for pre in dp[cur]}
            if deg[next_] == 1 and len(dp[next_]) == 1:  # 如果超过1,就是非法的
                queue.append(next_)

    # 非法
    dpSize = [len(s) for s in dp]
    if any(v > 2 or v == 0 for v in dpSize):
        print(-1)
        exit(0)
    if dpSize.count(2) > 1:
        print(-1)
        exit(0)

    size2 = next((i for i, v in enumerate(dpSize) if v == 2), -1)
    res = INF
    if size2 == -1:
        res = max(sum(s) for s in dp)
    else:
        res = sum(dp[size2])

    while res % 2 == 0:
        res //= 2
    print(res)
