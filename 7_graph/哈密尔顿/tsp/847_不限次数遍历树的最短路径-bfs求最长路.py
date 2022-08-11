# 给定一张包含 N 个点、 N-1 条边的无向连通图
# 每条边的长度均为 1 。假设你从 1 号节点出发并打算遍历所有节点，那么总路程至少是多少？


# 所以要遍历所有点，最少得走法是
# 所有边走两边再减去图中得最长路径(最大层数)

from collections import defaultdict, deque


adjMap = defaultdict(list)
n = int(input())
for _ in range(n - 1):
    u, v = map(int, input().split())
    adjMap[u].append(v)
    adjMap[v].append(u)

res = 0
queue = deque([1])
visited = set([1])
while queue:
    len_ = len(queue)
    # 注意要层序遍历
    for _ in range(len_):
        cur = queue.popleft()
        for next in adjMap[cur]:
            if next not in visited:
                visited.add(next)
                queue.append(next)
    res += 1

print(2 * (n - 1) - (res - 1))
