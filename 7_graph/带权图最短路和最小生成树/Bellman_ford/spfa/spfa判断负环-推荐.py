# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次


from collections import defaultdict, deque


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
