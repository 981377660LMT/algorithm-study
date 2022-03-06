from collections import defaultdict, deque


# 请你求出 1 号点到 n 号点的最短距离，如果无法从 1 号点走到 n 号点，则输出 impossible。
# 边权可能为负数。

# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能

n, m = map(int, input().split())
adjMap = defaultdict(list)
for _ in range(m):
    u, v, w = list(map(int, input().split()))
    adjMap[u].append((v, w))


def spfa(n: int, adjMap: defaultdict) -> bool:
    """判断负环要以所有点为起点，无须初始化dist"""
    dist = [0] * (n + 1)

    queue = deque()
    count = [0] * (n + 1)  # 边数
    isInqueue = [False] * (n + 1)  # 在队列里的点
    for i in range(1, n + 1):
        isInqueue[i] = True
        queue.append(i)

    while queue:
        cur = queue.popleft()
        isInqueue[cur] = False  # 点从队列出来了

        # 更新过谁，就拿谁去更新别人
        for next, weight in adjMap[cur]:
            if dist[cur] + weight < dist[next]:
                count[next] = count[cur] + 1
                if count[next] >= n:
                    return True
                dist[next] = dist[cur] + weight
                if not isInqueue[next]:
                    isInqueue[next] = True
                    queue.append(next)

    return False


res = spfa(n, adjMap)
print("Yes" if res else "No")
