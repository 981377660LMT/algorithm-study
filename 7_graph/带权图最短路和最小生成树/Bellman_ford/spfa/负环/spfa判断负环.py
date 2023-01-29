# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次


from collections import defaultdict, deque
from typing import DefaultDict


n, m = map(int, input().split())
adjMap = defaultdict(list)
for _ in range(m):
    # 编号从1开始
    u, v, w = list(map(int, input().split()))
    u, v = u - 1, v - 1
    adjMap[u].append((v, w))


def spfa(n: int, adjMap: DefaultDict) -> bool:
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
print("Yes" if res else "No")
