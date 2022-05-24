# 274. 移动服务
# 一个公司有三个移动服务员，最初分别在位置 1，2，3 处。
# 公司必须按顺序依次满足所有请求，且过程中不能去其他额外的位置，目标是最小化公司花费，请你帮忙计算这个最小花费。
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
n, q = map(int, input().split())  # 1-n个位置 和 q个请求
costs = []  # 从 p 到 q 移动一个员工，需要花费 c(p,q)。
for _ in range(n):
    costs.append([int(num) for num in input().split()])
queries = [int(num) - 1 for num in input().split()]


@lru_cache(None)
def dfs(index: int, p2: int, p3: int) -> int:
    """已经处理完前i个请求，且三个服务员分别在p[i], p2, p3的所有方案的集合；
    
    必有一个员工在queries[i]上了。所以dp状态设计的时候，可以省一个维度
    时间复杂度O(q*n^2)
    """
    if index == q - 1:
        return 0
    # 某一时刻只有一个员工能移动，且不允许在同样的位置出现两个员工。
    if queries[index] == p2 or queries[index] == p3 or p2 == p3:
        return int(1e20)
    p1, next = queries[index], queries[index + 1]
    res1 = costs[p1][next] + dfs(index + 1, *sorted((p2, p3)))
    res2 = costs[p2][next] + dfs(index + 1, *sorted((p1, p3)))
    res3 = costs[p3][next] + dfs(index + 1, *sorted((p1, p2)))

    return min(res1, res2, res3)


res = min(
    costs[0][queries[0]] + dfs(0, 1, 2),
    costs[1][queries[0]] + dfs(0, 0, 2),
    costs[2][queries[0]] + dfs(0, 0, 1),
)
dfs.cache_clear()
print(res)
