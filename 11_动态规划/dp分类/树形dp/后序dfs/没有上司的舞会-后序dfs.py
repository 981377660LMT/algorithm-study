# 状态:每个点,访问还是不访问
# dp = [[0, 1] for _ in range(n + 1)]
# 先处理子树，在处理当前root

from collections import defaultdict
import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))


# 后序dfs，当然也可以把结果存在外面dp
def dfs(cur: int) -> List[int]:
    """后序dfs返回[选当前，不选当前]"""
    res = [values[cur], 0]
    for next in adjMap[cur]:
        nextRes = dfs(next)
        res[0] += nextRes[1]
        res[1] += max(nextRes)

    return res


n = int(input())
values = [0] * n
for i in range(n):
    happy = int(input())
    values[i] = happy

adjMap = defaultdict(set)
indegree = [0] * n
for _ in range(n - 1):
    cur, pre = map(int, input().split())
    cur, pre = cur - 1, pre - 1
    adjMap[pre].add(cur)
    indegree[cur] += 1

root = next(i for i in range(n) if indegree[i] == 0)
print(max(dfs(root)))
