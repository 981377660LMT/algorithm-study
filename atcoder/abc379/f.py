import sys
import threading


def main():
    import sys
    from collections import deque, defaultdict

    sys.setrecursionlimit(1 << 25)
    MOD = 998244353

    H, W = map(int, sys.stdin.readline().split())
    S = [list(sys.stdin.readline().strip()) for _ in range(H)]

    # 二分图划分
    partition = [[-1 for _ in range(W)] for _ in range(H)]
    dirs = [(-1, 0), (1, 0), (0, -1), (0, 1)]
    components = []

    for i in range(H):
        for j in range(W):
            if partition[i][j] == -1:
                queue = deque()
                queue.append((i, j))
                partition[i][j] = 0
                nodes_A = []
                nodes_B = []
                nodes_A.append((i, j))
                valid = True
                while queue and valid:
                    x, y = queue.popleft()
                    for dx, dy in dirs:
                        nx, ny = x + dx, y + dy
                        if 0 <= nx < H and 0 <= ny < W:
                            if partition[nx][ny] == -1:
                                partition[nx][ny] = 1 - partition[x][y]
                                if partition[nx][ny] == 0:
                                    nodes_A.append((nx, ny))
                                else:
                                    nodes_B.append((nx, ny))
                                queue.append((nx, ny))
                            elif partition[nx][ny] == partition[x][y]:
                                print(0)
                                return
                components.append((nodes_A, nodes_B))

    ans = 1

    for nodes_A, nodes_B in components:
        # 记录A部分固定的颜色
        fixed_A = {}
        for x, y in nodes_A:
            if S[x][y] != "?":
                fixed_A[(x, y)] = int(S[x][y])

        # 记录B部分固定的颜色
        fixed_B = {}
        for x, y in nodes_B:
            if S[x][y] != "?":
                fixed_B[(x, y)] = int(S[x][y])

        # 列出A部分的节点
        list_A = nodes_A
        list_B = nodes_B

        # 准备B节点的邻接A节点
        B_neighbors = []
        for x, y in list_B:
            neighbors = []
            for dx, dy in dirs:
                nx, ny = x + dx, y + dy
                if 0 <= nx < H and 0 <= ny < W:
                    if partition[nx][ny] == 0:
                        neighbors.append((nx, ny))
            B_neighbors.append(neighbors)

        memo = {}
        n_A = len(list_A)

        def backtrack(idx, color_assign):
            if idx == n_A:
                # 计算B部分的可选颜色
                ways = 1
                for b_idx, (x, y) in enumerate(list_B):
                    forbidden = set()
                    for nx, ny in B_neighbors[b_idx]:
                        if (nx, ny) in color_assign:
                            forbidden.add(color_assign[(nx, ny)])
                    if (x, y) in fixed_B:
                        if fixed_B[(x, y)] in forbidden:
                            return 0
                        else:
                            continue
                    else:
                        available = 3 - len(forbidden)
                        if available <= 0:
                            return 0
                        ways = (ways * available) % MOD
                return ways

            key = tuple(sorted(color_assign.items()))
            if (idx, tuple(sorted(color_assign.items()))) in memo:
                return memo[(idx, tuple(sorted(color_assign.items())))]

            total = 0
            x, y = list_A[idx]
            if (x, y) in fixed_A:
                color = fixed_A[(x, y)]
                conflict = False
                # 检查与相邻B节点是否冲突
                # 暂不需要，因为B节点在后续处理
                color_assign[(x, y)] = color
                total = (total + backtrack(idx + 1, color_assign)) % MOD
                del color_assign[(x, y)]
            else:
                for color in [1, 2, 3]:
                    color_assign[(x, y)] = color
                    total = (total + backtrack(idx + 1, color_assign)) % MOD
                    del color_assign[(x, y)]
            memo[(idx, tuple(sorted(color_assign.items())))] = total
            return total

        count = backtrack(0, {})
        ans = (ans * count) % MOD
        if ans == 0:
            break

    print(ans)


main()
