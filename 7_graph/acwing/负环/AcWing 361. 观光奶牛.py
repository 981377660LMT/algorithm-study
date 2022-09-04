# 给定一张 n 个点、m 条边的有向图，每个点都有一个权值 f[i]，每条边都有一个权值 t[i]。
# 求图中的一个环，使“环上各点的权值之和”除以“环上各边的权值之和”最大。
# 输出这个最大值。
# 注意：数据保证至少存在一个环。
# n<=1000
# m<=5000

# 01分数规划 二分法
# https://www.acwing.com/solution/content/40640/
# !求 ∑f/∑t 的最大值
# ! ∑f/∑t > mid 即 ∑(mid*ti - fi) < 0  即存在负环
# ! mid∗ti−fi 为边权
# TODO

from collections import defaultdict, deque
import sys
from typing import Mapping

EPS = 1e-3


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m = map(int, input().split())
points = [int(input()) for _ in range(n)]  # 点权
adjMap = defaultdict(lambda: defaultdict(lambda: INF))
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u][v] = min(adjMap[u][v], w)
    adjMap[v][u] = min(adjMap[v][u], w)


def check(mid: float) -> int:
    """mid时是否存在负环"""

    def spfa2(n: int, adjMap: Mapping[int, Mapping[int, int]]) -> bool:
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
                weight = adjMap[cur][next]  # 注意这里的边权 TODO
                cand = dist[cur] + weight
                if cand < dist[next]:
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

    return spfa2(n, adjMap)  # 存在负环


left, right = 0, 1
while left <= right:
    mid = (left + right) / 2
    if check(mid):
        left = mid + EPS
    else:
        right = mid - EPS

print(format(right, ".2f"))
