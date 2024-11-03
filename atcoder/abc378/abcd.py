from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    q, r = [], []
    for _ in range(N):
        a, b = map(int, input().split())
        q.append(a)
        r.append(b)
    Q = int(input())
    for _ in range(Q):
        t, d = map(int, input().split())
        qi = q[t - 1]
        ri = r[t - 1]
        k = (d - ri + qi - 1) // qi
        D = qi * k + ri
        print(D)
