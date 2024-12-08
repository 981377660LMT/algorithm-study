import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    times = []
    volumes = []
    for _ in range(n):
        t, v = map(int, input().split())
        times.append(t)
        volumes.append(v)

    currentTime = 0
    water = 0

    for i in range(n):
        elapsedTime = times[i] - currentTime
        water = max(0, water - elapsedTime)
        water += volumes[i]
        currentTime = times[i]

    print(water)
