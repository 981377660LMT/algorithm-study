# 给定一个无向图，两个独立集 S,T，(独立集就是集合内部的点没有被边连接)，
# S = range(1, s + 1)
# T = range(s + 1, s + t + 1)
# !有 m 条无向边，请找到一个四元环，若没有，输出-1。
# S<=3e5 T<=3e3
# !O(S+m+T^2)


# !四元环只能是S-T-S-T构成
# 观察数据量
# !枚举S 再枚举与Si相连的两个不同的点 看之前这两个点是否连过S中的点
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    s, t, m = map(int, input().split())
    adjMap = [set() for _ in range(s + t)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - s - 1
        adjMap[u].add(v)  # 无向边定向

    connect = [[-1] * (t + 10) for _ in range(t + 10)]  # !不能开字典 会TLE
    for v1 in range(s):
        for v2 in adjMap[v1]:
            for v3 in adjMap[v1]:
                if v2 >= v3:
                    continue
                if connect[v2][v3] == -1:
                    connect[v2][v3] = v1
                else:
                    print(v2 + s + 1, v3 + s + 1, connect[v2][v3] + 1, v1 + 1)
                    exit(0)
    print(-1)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
