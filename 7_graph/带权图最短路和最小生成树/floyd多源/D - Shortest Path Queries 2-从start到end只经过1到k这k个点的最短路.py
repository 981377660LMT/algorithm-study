# 最短路查询 只经过k个点的最短路
# 定义f(start,end,k)为从start到end只经过1到k这k个点的最短路
# !求所有f(start,end,k) (1<=start<=n,1<=end<=n,1<=k<=n) 的和 (无法到达则f=0)
# n<=400

# !对floyd松弛过程的理解
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    dist = [[INF] * n for _ in range(n)]
    for i in range(n):
        dist[i][i] = 0
    for _ in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        if w < dist[u][v]:
            dist[u][v] = w

    res = 0
    for k in range(n):
        for i in range(n):
            for j in range(n):
                cand = dist[i][k] + dist[k][j]
                if dist[i][j] > cand:
                    dist[i][j] = cand
        for i in range(n):
            for j in range(n):
                if dist[i][j] != INF:
                    res += dist[i][j]
    print(res)
