# https://zhuanlan.zhihu.com/p/496282947
# 题意
# 有n个人来自不同的学校ai,擅长不同的学科bi,每个人有一个能力值ci
# 要求组建一支i个人的梦之队最大化队员的能力值,并且满足队伍中所有人来自的学校和擅长的学科都不同.
# n<=3e4
# !ai,bi<=150 (暗示作为顶点的数据量)
# ci<=1e9

# 分析
# 把学校和学科看作点,把一个人看成匹配边,能力值看作边权,其实就是求有i条匹配边的最优匹配.可以用费用流解决.
# 此外题目要求输出匹配数为1,2,…k个匹配时的最优匹配.
# !在spfa费用流算法中一次spfa只会找到一条费用最小的增广流,
# !所以每次增广之后得到的费用就对应匹配数为1,2,…k个匹配时的答案.

# !能力值对应 流的费用
# !队伍里每个人学校和学科都不同:学校学科为虚拟源汇点，容量为1，这样就不会取到重复的学生了

from heapq import heappop, heappush
import sys
from typing import Optional


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(1e18)


# https://atcoder.jp/contests/abc247/submissions/30874572
class ATCMinCostFlow:
    def __init__(self, n: int):
        self._n = n
        self.graph = [[] for _ in range(n)]

    def addEdge(self, fromV: int, toV: int, cap: int, cost: int) -> None:
        forward = [toV, cap, cost, None]
        backward = forward[3] = [fromV, 0, -cost, forward]  # type: ignore
        self.graph[fromV].append(forward)
        self.graph[toV].append(backward)

    def flow(self, start: int, end: int, flow: int) -> Optional[int]:
        n = self._n
        g = self.graph

        res = 0
        head = [0] * n
        preV = [0] * n
        preE = [None] * n

        d0 = [INF] * n
        dist = [INF] * n

        while flow:
            dist[:] = d0
            dist[start] = 0
            que = [(0, start)]

            while que:
                c, v = heappop(que)
                if dist[v] < c:
                    continue
                r0 = dist[v] + head[v]
                for e in g[v]:
                    w, cap, cost, _ = e
                    if cap > 0 and r0 + cost - head[w] < dist[w]:
                        dist[w] = r = r0 + cost - head[w]
                        preV[w] = v
                        preE[w] = e
                        heappush(que, (r, w))
            if dist[end] == INF:
                return None

            for i in range(n):
                head[i] += dist[i]

            d = flow
            v = end
            while v != start:
                d = min(d, preE[v][1])
                v = preV[v]
            flow -= d
            res += d * head[end]
            v = end
            while v != start:
                e = preE[v]
                e[1] -= d
                e[3][1] += d
                v = preV[v]
        return res


#####################################

n = int(input())
v = 150
START, END, OFFSET = 2 * v, 2 * v + 1, v
mcmf = ATCMinCostFlow(2 * v + 2)

for i in range(v):
    mcmf.addEdge(START, i, 1, 0)
    mcmf.addEdge(i + OFFSET, END, 1, 0)

for _ in range(n):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    mcmf.addEdge(u, v + OFFSET, 1, -w)  # !要求最大费用流 所以(BIG-w)费用


res = []
cur = 0
while True:
    delta = mcmf.flow(START, END, 1)  # !流量为1
    if delta is None:
        break
    cur += delta
    res.append(cur)


print(len(res))
for i, cost in enumerate(res, 1):
    print(-cost)

# TODO 不一定是k限流
# 4
# 1 1 1
# 1 2 1
# 2 3 1
# 3 3 1


# 2
# 1
# 2

# !输出为 1 2 少了一个匹配的情况
