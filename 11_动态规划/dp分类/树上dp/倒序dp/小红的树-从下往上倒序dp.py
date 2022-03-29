# 现在小红想给一部分点染成红色。之后她有 q 次询问，每次询问某点的`子树`红色节点的个数。
# 根节点为1号结点

n = int(input())
parents = [int(i) for i in input().split()]
colors = input()

# 建图
adjList = [[] for _ in range(n + 1)]
for i in range(n - 1):
    adjList[parents[i]].append(i + 2)

queryLen = int(input())
queries = []
for _ in range(queryLen):
    queries.append(int(input()))

res = [0] * (n + 1)

# 从最底层向上更新
for root in range(n, 0, -1):
    if colors[root - 1] == 'R':
        res[root] += 1
    for next in adjList[root]:
        res[root] += res[next]
for i in queries:
    print(res[i])
