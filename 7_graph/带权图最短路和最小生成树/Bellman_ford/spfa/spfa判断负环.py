# https://www.acwing.com/activity/content/problem/content/921/

# !只是找负环的话，初始时将所有点入队即可
# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次


import sys
from collections import defaultdict, deque
from typing import Mapping

sys.setrecursionlimit(int(1e9))
INF = int(1e18)


def spfa1(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> bool:
    """spfa判断负环 存在负环返回True 否则返回False

    dfs 加速
    https://www.acwing.com/blog/content/22651/
    """

    def dfs(cur: int) -> None:
        nonlocal hashMinusCycle
        states[cur] = 1
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            if weight + dist[cur] < dist[next]:
                dist[next] = weight + dist[cur]
                if states[next] == 1:
                    hashMinusCycle = True
                    return
                dfs(next)
                if hashMinusCycle:
                    return
        states[cur] = 2

    dist = [0] * n
    states = [0] * n
    hashMinusCycle = False
    for i in range(n):
        if states[i]:
            continue
        if hashMinusCycle:
            return True
        dfs(i)
    return hashMinusCycle


def spfa2(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> bool:
    """spfa判断负环 存在负环返回True 否则返回False

    dfs 加速
    https://www.acwing.com/solution/content/87368/
    """

    def detectCycle() -> bool:
        stack = []
        inStack = [False] * n
        visited = [False] * n
        for cur in range(n):
            if not visited[cur]:
                next = cur
                while True:
                    if next == -1:
                        break
                    if not visited[next]:
                        visited[next] = True
                        stack.append(next)
                        inStack[next] = True
                    else:
                        if inStack[next]:
                            return True
                        break
                    next = pre[next]

                for num in stack:
                    inStack[num] = False
                stack = []
        return False

    dist = [0] * n
    pre = [-1] * n
    queue = deque(range(n))
    inQueue = [True] * n
    idx = 0  # 计数器
    while queue:
        cur = queue.popleft()
        inQueue[cur] = False
        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                pre[next] = cur
                idx += 1
                if idx == n:
                    idx = 0
                    if detectCycle():
                        return True
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:
                        queue.appendleft(next)
                    else:
                        queue.append(next)
    return detectCycle()


def spfa3(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> bool:
    """spfa判断负环 存在负环返回True 否则返回False

    在原图的基础上新建一个虚拟源点,
    从该点向其他所有点连一条权值为0的有向边。
    那么原图有负环等价于新图有负环
    也等价于开始时将所有点加入队列
    """
    dist = [0] * n
    queue = deque(range(n))
    inQueue = [True] * n
    count = [0] * n

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False

        for next in adjMap[cur]:
            weight = adjMap[cur][next]
            cand = dist[cur] + weight
            if cand < dist[next]:  # 如果要最长路这里需要改成 >
                dist[next] = cand
                count[next] = count[cur] + 1
                if count[next] >= n:
                    return True
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化 如果要最长路这里需要改成 >
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return False


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))  # 最短路
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u][v] = min(adjMap[u][v], w)

    res = spfa1(n, adjMap)
    print("Yes" if res else "No")
