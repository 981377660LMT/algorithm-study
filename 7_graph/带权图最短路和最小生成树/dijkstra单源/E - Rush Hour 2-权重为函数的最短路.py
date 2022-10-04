"""权重为函数的最短路"""

# 一个国家有N座城市和M条路。城市的编号是1至N,路的编号则为1至M。
# 第i条路双向连接着城市A和B.
# !在这个国家中,如果你在时间点t通过公路i,所需时间为 `Ci + floor(Di/(t+1))`
# 一个人想从城市1到达城市n。他在每个城市可以停留任意自然数的时间。
# 求这个人最早到达城市N的时间。如果无法到达城市n,输出-1 。
# N,M<=1e5

# 假设t时刻离开Ai前往Bi
# !那么t时刻到达Bi的时刻为f(t)=t+Ci+floor(Di/(t+1))
# !这是一个双钩函数
# !结论:最好的整数出发时间为 `floor(sqrt(D)) - 1` 可以在这个范围上枚举几个值

from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def calMinTime(startTime: int, C: int, D: int) -> int:
        """计算最早的到达时间"""
        res = startTime + C + D // (startTime + 1)  # !从startTime时刻准时出发
        sqrt_ = int(D**0.5)
        for cand in range(sqrt_ - 3, sqrt_ + 3):  # !不知道精准的解就在范围内枚举几个值
            if cand >= startTime:
                res = min(res, cand + C + D // (cand + 1))
        return res

    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        a, b, c, d = map(int, input().split())
        a, b = a - 1, b - 1
        if a != b:
            adjList[a].append((b, c, d))
            adjList[b].append((a, c, d))

    dist = [INF] * n
    dist[0] = 0
    pq = [(0, 0)]
    while pq:
        curTime, cur = heappop(pq)
        if curTime > dist[cur]:
            continue
        if cur == n - 1:
            print(curTime)
            exit(0)

        for next, C, D in adjList[cur]:
            cand = calMinTime(curTime, C, D)
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (cand, next))

    print(-1)
