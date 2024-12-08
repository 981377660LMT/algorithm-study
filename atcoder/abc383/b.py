import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    h, w, d = map(int, input().split())
    grid = [input() for _ in range(h)]
    floors = []

    for i in range(h):
        for j in range(w):
            if grid[i][j] == ".":
                floors.append((i, j))

    maxCount = 0

    for i in range(len(floors)):
        for j in range(i + 1, len(floors)):
            humidified = set()

            for x, y in [floors[i], floors[j]]:
                for a, b in floors:
                    if abs(x - a) + abs(y - b) <= d:
                        humidified.add((a, b))

            maxCount = max(maxCount, len(humidified))

    print(maxCount)
