# 2N个人 每次从队列中选两个人出队 这两个人身高差越大 点踩的人越多
# 问如何选择出队顺序 使得点踩人数最少
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


n = int(input())
nums = list(map(int, input().split()))


@lru_cache(None)
def dfs(left: int, right: int) -> int:
    """[left, right] 间的点踩人数
    
    1. 两个端点全取
    2. 否则划分区间取
    """
    if left >= right:
        return int(1e20)
    if left + 1 == right:
        return abs(nums[left] - nums[right])

    res = dfs(left + 1, right - 1) + abs(nums[left] - nums[right])
    for i in range(left + 1, right, 2):
        res = min(res, dfs(left, i) + dfs(i + 1, right))
    return res


print(dfs(0, 2 * n - 1))

