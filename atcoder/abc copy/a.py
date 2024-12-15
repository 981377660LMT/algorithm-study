import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, c1, c2 = input().split()
    n = int(n)
    s = input()

    result = "".join(ch if ch == c1 else c2 for ch in s)

    print(result)
