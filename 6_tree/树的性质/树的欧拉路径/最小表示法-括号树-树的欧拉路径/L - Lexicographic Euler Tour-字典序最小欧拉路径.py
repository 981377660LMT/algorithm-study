# https://atcoder.jp/contests/nadafes2022_day1/tasks/nadafes2022_day1_l?lang=ja
# L - Lexicographic Euler Tour-字典序最小欧拉路径
# 这里记录的是深度,起点深度为0

# 4
# 1 2
# 2 3
# 1 4
# 0-1-0-1-2-1-0


from collections import deque
from typing import Deque, List


def lexicoGraphicEulerTour(adjList: List[List[int]]) -> List[int]:
    def dfs(cur: int, pre: int, dep: int) -> Deque[int]:
        sub = sorted([dfs(next, cur, dep + 1) for next in adjList[cur] if next != pre])
        res = deque([dep])  # !进入
        for d in sub:
            if len(res) > len(d):
                while d:
                    res.append(d.popleft())
            else:
                res, d = d, res
                while d:
                    res.appendleft(d.pop())
            res.append(dep)  # !回溯
        return res

    res = dfs(0, -1, 0)
    return list(res)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        a, b = map(int, input().split())
        adjList[a - 1].append(b - 1)
        adjList[b - 1].append(a - 1)
    print(*lexicoGraphicEulerTour(adjList), end=" ")
