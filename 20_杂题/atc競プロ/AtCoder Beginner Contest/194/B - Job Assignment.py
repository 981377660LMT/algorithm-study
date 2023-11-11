import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 工作任务分配
# !将两个工作分给至多两个员工，求最小的完成时间
# n<=1e5
# O(n)解法


def jobAssignment1(costs: List[Tuple[int, int]]) -> int:
    cost1, cost2 = [list(row) for row in zip(*costs)]
    i, j = cost1.index(min(cost1)), cost2.index(min(cost2))
    if i != j:  # 两个最小值不在同一个位置
        return max(cost1[i], cost2[j])
    min1, min2 = cost1[i], cost2[j]
    cost1[i], cost2[j] = INF, INF
    return min(min1 + min2, max(min1, min(cost2)), max(min2, min(cost1)))


# def jobAssignment2(costs: List[Tuple[int, int]]) -> int:
#     n = len(costs)
#     res = INF
#     for i in range(n):
#         for j in range(n):
#             if i == j:
#                 res = min(res, costs[i][0] + costs[j][1])
#             else:
#                 res = min(res, max(costs[i][0], costs[j][1]))
#     return res


if __name__ == "__main__":
    n = int(input())
    costs = [tuple(map(int, input().split())) for _ in range(n)]  # (cost1,cost2)
    print(jobAssignment1(costs))
