# 虫洞非常奇特，它可以看作是一条 单向 路径，
# 通过它可以使你回到过去的某个时刻（相对于你进入虫洞之前）。
# 现在农夫约翰希望能够从农场中的某片田地出发，经过一些路径和虫洞回到过去，并在他的出发时刻之前赶到他的出发地。


from spfa判断负环 import spfa


n = int(input())

# 使用酸辣粉优化后	4280 ms => 3570 ms
for _ in range(n):
    N, M, W = map(int, input().split())
    adjList = [[] for _ in range(N)]
    for _ in range(M):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, w))
        adjList[v].append((u, w))

    for _ in range(W):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, -w))

    res = spfa(N, adjList)
    print("YES" if res else "NO")
