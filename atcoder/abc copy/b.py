import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, r = map(int, input().split())
    for _ in range(n):
        d, a = map(int, input().split())
        if d == 1 and 1600 <= r <= 2799:
            r += a
        elif d == 2 and 1200 <= r <= 2399:
            r += a
    print(r)
