from typing import List, Optional


def find_cycle(g: List[List[int]]) -> Optional[List[int]]:
    n = len(g)
    used = [0] * n  # 0:not yet 1: visiting 2: visited
    for v in range(n):  # 各点でDFS
        if used[v] == 2:
            continue
        # 初期化
        stack = [v]
        hist = []  # 履歴
        while stack:
            v = stack[-1]
            if used[v] == 1:
                used[v] = 2  # 帰りがけの状態に
                stack.pop()
                hist.pop()
                continue
            hist.append(v)
            used[v] = 1  # 行きがけの状態に
            for c in g[v]:
                if used[c] == 2:
                    continue
                elif used[c] == 1:  # cを始点とするサイクル発見！
                    return hist[hist.index(c) :]
                else:
                    stack.append(c)
    return None


if __name__ == "__main__":

    def abc_456_e():
        T = int(input())
        for _ in range(T):
            n, m = map(int, input().split())
            edges = [tuple(map(lambda x: int(x) - 1, input().split())) for _ in range(m)]
            w = int(input())
            states = [input() for _ in range(n)]

            adjList = [[] for _ in range(n * w)]
            for u, v in edges:
                for j in range(w):
                    if states[v][(j + 1) % w] == "o":
                        adjList[u * w + j].append(v * w + (j + 1) % w)
                    if states[u][(j + 1) % w] == "o":
                        adjList[v * w + j].append(u * w + (j + 1) % w)
            for i in range(n):
                for j in range(w):
                    if states[i][(j + 1) % w] == "o":
                        adjList[i * w + j].append(i * w + (j + 1) % w)

            cycle = find_cycle(adjList)
            print("Yes" if cycle else "No")

    abc_456_e()
