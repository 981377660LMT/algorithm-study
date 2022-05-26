# 无向图中求1到N的所有可能路径中第K+1大的边的最小值是多少
# (二分+dijkstra)  O(nlog(n))
import collections
import heapq

INF = int(1e20)
adjMap = collections.defaultdict(list)

# N 座通信基站，P 条 双向 电缆 求第k+1大的边长最小值
n, p, k = map(int, input().split())
for _ in range(p):
    x, y, cost = map(int, input().split())
    x -= 1  # 结点编号，统一成从0开始
    y -= 1
    adjMap[x].append((y, cost))
    adjMap[y].append((x, cost))


# ----dijkstra算法
def check(mid: int) -> bool:
    """
    从 1→n 的路径中应该存在一条路，使得这条路上最多有 k 条大于mid的边

    即最短路上最多有k条大于mid的边 置为01即可
    """
    dist = [INF] * n
    dist[0] = 0
    pq = [(0, 0)]
    while pq:
        _, cur = heapq.heappop(pq)
        for next, cost in adjMap[cur]:
            cost = 1 if cost > mid else 0  # 方便统计有多少条边的长度超过了mid
            if dist[cur] + cost < dist[next]:
                dist[next] = dist[cur] + cost
                heapq.heappush(pq, (dist[next], next))
    return dist[n - 1] <= k


# ----二分--最左端
left = 0
right = int(1e10)
ok = False
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        ok = True
        right = mid - 1
    else:
        left = mid + 1

print(left if ok else -1)

