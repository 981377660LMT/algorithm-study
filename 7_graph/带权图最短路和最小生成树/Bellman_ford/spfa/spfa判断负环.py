# https://www.acwing.com/activity/content/problem/content/921/

# !只是找负环的话，初始时将所有点入队即可
# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次


from collections import deque
from typing import List, Tuple


def spfa(n: int, adjList: List[List[Tuple[int, int]]]) -> bool:
    """spfa判断负环 存在负环返回True 否则返回False

    在原图的基础上新建一个虚拟源点,
    从该点向其他所有点连一条权值为0的有向边。
    那么原图有负环等价于新图有负环
    也等价于开始时将所有点加入队列
    """
    dist = [0] * n
    queue = deque(range(n))
    inQueue = [True] * n
    count = [1] * n

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:  # !如果要正环这里需要改成 >
                dist[next] = cand
                if not inQueue[next]:
                    count[next] += 1
                    if count[next] >= n + 1:  # +1是虚拟源点
                        return True
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化 如果要正环这里需要改成 >
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return False


if __name__ == "__main__":
    # 有向图中是否存在负环
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]  # 最短路
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, w))

    hasCycle = spfa(n, adjList)
    print("Yes" if hasCycle else "No")
