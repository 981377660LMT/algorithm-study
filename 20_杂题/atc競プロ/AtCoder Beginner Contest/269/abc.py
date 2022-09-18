import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    x = int(input())
    bin_ = bin(x)[2:]
    res = set([0])
    g1, g2 = x, 0
    while g1:
        res.add(g1)
        g1 = (g1 - 1) & x
        g2 = x ^ g1
    res = sorted(res)
    print(*res, sep="\n")
