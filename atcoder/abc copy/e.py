from collections import deque
import sys
import heapq

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    H, W, X = map(int, input().split())
    P, Q = map(int, input().split())
    P -= 1
    Q -= 1
    S = []
    for _ in range(H):
        S.append(list(map(int, input().split())))
    strength = S[P][Q]
    absorbed = [[False] * W for _ in range(H)]
    absorbed[P][Q] = True
    heap = []
    directions = [(-1, 0), (1, 0), (0, -1), (0, 1)]
    for di, dj in directions:
        ni, nj = P + di, Q + dj
        if 0 <= ni < H and 0 <= nj < W and not absorbed[ni][nj]:
            heapq.heappush(heap, (S[ni][nj], ni, nj))
    while heap:
        s, i, j = heapq.heappop(heap)
        if absorbed[i][j]:
            continue
        if s * X < strength:
            strength += s
            absorbed[i][j] = True
            for di, dj in directions:
                ni, nj = i + di, j + dj
                if 0 <= ni < H and 0 <= nj < W and not absorbed[ni][nj]:
                    heapq.heappush(heap, (S[ni][nj], ni, nj))
        else:
            continue
    print(strength)
