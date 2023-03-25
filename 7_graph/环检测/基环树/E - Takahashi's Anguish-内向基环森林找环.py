# 发糖果
# 每个人有一个讨厌的人pi(不是自己)
# 如果讨厌的人在他之前得到了糖果，那么他的生气值为ai
# 最小化生气值之和
# n<=2e5

# !https://atcoder.jp/contests/abc256/editorial/4135
# ! n个顶点n条边 Namori Graph(基环树森林)
# ! 内向基环树森林找环 每个环中必定有一个点作为起点 怒气值最小
# ! 只有在环里才能产生怒气 如果没环直接按照拓扑序发饼干 怒气值为0
import sys
from collections import defaultdict, deque
from typing import DefaultDict, List, Set, Tuple

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def findCycleAndCalDepth1(
    n: int, adjMap: DefaultDict[int, Set[int]], indeg: List[int]
) -> Tuple[List[List[int]], List[int]]:
    """内向基环树找环上的点，并记录每个点在拓扑排序中的最大距离,最外层的点深度为0"""
    depth = [0] * n
    queue = deque([(i, 0) for i in range(n) if indeg[i] == 0])
    visited = [False] * n
    while queue:
        cur, dist = queue.popleft()
        visited[cur] = True
        for next in adjMap[cur]:
            depth[next] = max(depth[next], dist + 1)
            indeg[next] -= 1
            if indeg[next] == 0:
                queue.append((next, dist + 1))

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in adjMap[cur]:
            dfs(next, path)

    cycleGroup = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        cycleGroup.append(path)

    return cycleGroup, depth


def main() -> None:
    n = int(input())
    hates = [int(num) - 1 for num in input().split()]
    scores = list(map(int, input().split()))

    adjMap = defaultdict(set)
    indeg = [0] * n
    for i in range(n):
        adjMap[i].add(hates[i])
        indeg[hates[i]] += 1

    cycleGroup, _ = findCycleAndCalDepth1(n, adjMap, indeg)
    res = 0
    for group in cycleGroup:
        res += min(scores[i] for i in group)
    print(res)


if sys.argv[0].startswith(r"e:\test"):
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()
