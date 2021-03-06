from collections import defaultdict
from typing import DefaultDict, Set

# https://blog.csdn.net/baidu_23955875/article/details/46573537
# todo
# 简化版的匈牙利学习 （不判断二分图）
# 有两台机器 A，B 以及 K 个任务。
# 机器 A 有 N 种不同的模式（模式 0∼N−1），机器 B 有 M 种不同的模式（模式 0∼M−1）。
# 两台机器最开始都处于模式 0。
# 每个任务既可以在 A 上执行，也可以在 B 上执行。
# 对于每个任务 i，给定两个整数 a[i] 和 b[i]，表示如果该任务在 A 上执行，需要设置模式为 a[i]，如果在 B 上执行，需要模式为 b[i]。
# 任务可以以任意顺序被执行，但每台机器转换一次模式就要重启一次。
# 求怎样分配任务并合理安排顺序，能使机器重启次数最少。

# (最小覆盖点==最大匹配边-匈牙利算法)


def hungarian(adjMap: DefaultDict[int, Set[int]]) -> int:
    def getColor(adjMap: DefaultDict[int, Set[int]]) -> DefaultDict[int, int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in adjMap[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception('不是二分图')

        colors = defaultdict(lambda: -1)
        for i in range(n):
            if colors[i] == -1:
                dfs(i, 0)
        return colors

    def dfs(boy: int) -> bool:
        """寻找增广路"""
        nonlocal visited
        if boy in visited:
            return False
        visited.add(boy)

        for girl in adjMap[boy]:
            if matching[girl] == -1 or dfs(matching[girl]):
                matching[boy] = girl
                matching[girl] = boy
                return True
        return False

    n = len(adjMap)
    maxMatching = 0
    matching = defaultdict(lambda: -1)
    colors = getColor(adjMap)
    visited = set()
    for i in range(n):
        visited = set()
        if colors[i] == 0 and matching[i] == -1:
            if dfs(i):
                maxMatching += 1

    return maxMatching


n, m, k = map(int, input().split())
adjMap = defaultdict(set)

for _ in range(k):
    _, a, b = map(int, input().split())

    # 原来的两个点已经覆盖了
    if a == 0 or b == 0:
        continue
    adjMap[a].add(n + b)
    adjMap[n + b].add(a)
print(hungarian(adjMap))
