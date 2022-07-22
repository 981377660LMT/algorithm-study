# 给定一张 N 个点 M 条边的有向无环图，分别统计从每个点出发能够到达的点的数量。

# DAG上的dp 有向无环图没有循环依赖 dp无后效性 可以做拓扑序dp
# 注意点要去重 比较适合用Bitset


from collections import deque


n, m = map(int, input().split())
adjList = [set() for _ in range(n)]
indeg = [0] * n
visitedPair = set()
for _ in range(m):
    cur, next = map(int, input().split())
    next, cur = next - 1, cur - 1
    if (next, cur) not in visitedPair:
        visitedPair.add((next, cur))
        adjList[next].add(cur)
        indeg[cur] += 1

# set TLE了 改用二进制 还是在30000 TLE
queue = deque([i for i in range(n) if indeg[i] == 0])
res = [(1 << i) for i in range(n)]
while queue:
    cur = queue.popleft()
    for next in adjList[cur]:
        indeg[next] -= 1
        res[next] |= res[cur]
        if indeg[next] == 0:
            queue.append(next)

for i in range(n):
    print(bin(res[i])[2:].count("1"))
    # print(res[i].bit_count())
