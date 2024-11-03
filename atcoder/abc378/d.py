from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    H, W, K = map(int, input().split())
    G = [input() for _ in range(H)]
    available_cells = []
    idx_map = {}
    ptr = 0
    for i in range(H):
        for j in range(W):
            if G[i][j] == ".":
                available_cells.append((i, j))
                idx_map[(i, j)] = ptr
                ptr += 1

    adj_list = [[] for _ in range(len(available_cells))]
    for ptr, (i, j) in enumerate(available_cells):
        for dx, dy in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
            ni, nj = i + dx, j + dy
            if 0 <= ni < H and 0 <= nj < W and G[ni][nj] == ".":
                n_idx = idx_map[(ni, nj)]
                adj_list[ptr].append(n_idx)

    total_paths = 0

    for start_idx in range(len(available_cells)):
        visited = 1 << start_idx
        stack = [(start_idx, 0, visited)]
        while stack:
            current_idx, steps, visited_mask = stack.pop()
            if steps == K:
                total_paths += 1
                continue
            for neighbor_idx in adj_list[current_idx]:
                if not (visited_mask & (1 << neighbor_idx)):
                    stack.append((neighbor_idx, steps + 1, visited_mask | (1 << neighbor_idx)))

    print(total_paths)
