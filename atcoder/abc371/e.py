from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))

    indexMap = defaultdict(list)
    for i, char in enumerate(A):
        indexMap[char].append(i)

    res = 0
    for indexes in indexMap.values():
        res += N * (N + 1) // 2

        indexes = [-1] + indexes + [N]
        for pre, cur in zip(indexes, indexes[1:]):
            count = cur - pre
            res -= count * (count - 1) // 2

    print(res)
