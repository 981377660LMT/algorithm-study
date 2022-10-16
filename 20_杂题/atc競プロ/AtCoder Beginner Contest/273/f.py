import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, x = map(int, input().split())
    wall = list(map(int, input().split()))
    hammer = list(map(int, input().split()))

    # 拓扑排序?
