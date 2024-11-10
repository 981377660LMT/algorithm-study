import sys
import threading


def main():
    sys.setrecursionlimit(1 << 25)
    mod = 998244353

    H, W = map(int, sys.stdin.readline().split())
    grid = [list(sys.stdin.readline().strip()) for _ in range(H)]

    N = H * W
    cell_index = lambda x, y: x * W + y

    adj = [[] for _ in range(N)]
    pre_colors = [0] * N  # 0: unassigned, 1,2,3: colors

    for i in range(H):
        for j in range(W):
            idx = cell_index(i, j)
            if grid[i][j] in "123":
                pre_colors[idx] = int(grid[i][j])
            for dx, dy in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
                ni, nj = i + dx, j + dy
                if 0 <= ni < H and 0 <= nj < W:
                    nidx = cell_index(ni, nj)
                    adj[idx].append(nidx)

    visited = [False] * N
    total_answer = 1

    for v in range(N):
        if visited[v]:
            continue

        # BFS to check bipartiteness and collect pre-assigned colors
        queue = [v]
        visited[v] = True
        partition = [0] * N  # 0 or 1
        partition[v] = 0
        partition_colors = [set(), set()]
        is_bipartite = True

        if pre_colors[v]:
            partition_colors[partition[v]].add(pre_colors[v])

        while queue and is_bipartite:
            u = queue.pop(0)
            for neighbor in adj[u]:
                if not visited[neighbor]:
                    visited[neighbor] = True
                    partition[neighbor] = 1 - partition[u]
                    if pre_colors[neighbor]:
                        partition_colors[partition[neighbor]].add(pre_colors[neighbor])
                    queue.append(neighbor)
                else:
                    if partition[neighbor] == partition[u]:
                        # Not bipartite
                        is_bipartite = False
                        break
                    else:
                        # Check for color conflicts
                        if pre_colors[neighbor] and pre_colors[u]:
                            if pre_colors[neighbor] == pre_colors[u]:
                                is_bipartite = False
                                break
        if not is_bipartite:
            total_answer = 0
            break

        # Calculate the number of valid colorings for this component
        count = 1
        colors_used = [set(), set()]
        for p in [0, 1]:
            if len(partition_colors[p]) == 0:
                colors_used[p] = set([1, 2, 3])
            elif len(partition_colors[p]) == 1:
                colors_used[p] = partition_colors[p]
            else:
                # Conflict in pre-assigned colors
                total_answer = 0
                break
        if total_answer == 0:
            break
        # Remove colors used in the other partition
        for p in [0, 1]:
            other_p = 1 - p
            colors_used[p] = colors_used[p] - colors_used[other_p]
            if len(colors_used[p]) == 0:
                total_answer = 0
                break
        if total_answer == 0:
            break
        # Multiply the number of choices
        count = 1
        for p in [0, 1]:
            count = (count * len(colors_used[p])) % mod
        total_answer = (total_answer * count) % mod

    print(total_answer)


main()
