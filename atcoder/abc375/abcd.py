from math import dist
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    newGrid = [[0] * N for _ in range(N)]
    for i in range(N):
        s = input()
        newGrid[i] = [1 if s[j] == "#" else 0 for j in range(N)]

    for i in range(N // 2 + 1):
        for j in range(N // 2 + 1):
            if (N - i + N - 1 - j) % 2 == 1:
                x1, y1 = i, j
                x2, y2 = N - 1 - i, N - 1 - j
                newGrid[x1][y1], newGrid[x2][y2] = newGrid[x2][y2], newGrid[x1][y1]

    for i in range(N):
        print("".join(["#" if newGrid[i][j] == 1 else "." for j in range(N)]))
