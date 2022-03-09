from collections import defaultdict, deque

# 最坏O(VE)
f = int(input())

# 使用酸辣粉优化后	4280 ms => 3570 ms
for _ in range(f):
    n, M, W = map(int, input().split())
    adjMap = defaultdict(list)
    for _ in range(M):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u].append((v, w))
        adjMap[v].append((u, w))

    for _ in range(W):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u].append((v, -w))

    # 现在农夫约翰希望能够从农场中的某片田地出发，经过一些路径和虫洞回到过去，并在他的出发时刻之前赶到他的出发地。
    def spfa(n: int, adjMap: defaultdict) -> bool:
        """判断负环要以所有点为起点"""
        dist = [0] * (n)

        queue = deque()
        count = [0] * (n)  # 边数
        isInqueue = [False] * (n)  # 在队列里的点
        for i in range(n):
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
                        # 酸辣粉优化：队列不为空，且当前元素距离小于队头，则加入队头，否则加入队尾
                        if queue and dist[next] < dist[queue[0]]:
                            queue.appendleft(next)
                        else:
                            queue.append(next)

        return False

    res = spfa(n, adjMap)
    print("YES" if res else "NO")

