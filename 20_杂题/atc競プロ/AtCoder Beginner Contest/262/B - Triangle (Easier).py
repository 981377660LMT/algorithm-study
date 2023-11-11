from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def main() -> None:
    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = sorted([u - 1, v - 1])
        adjMap[u].add(v)

    res = 0
    for v1 in range(n):
        for v2 in adjMap[v1]:
            for v3 in adjMap[v2]:
                if v3 in adjMap[v1]:
                    res += 1
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
