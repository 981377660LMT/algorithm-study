# 每次只能合并相邻的两堆，合并的代价为这两堆石子的质量之和，
# 合并后与这两堆石子相邻的石子将和新堆相邻，
# 合并时由于选择的顺序不同，合并的总代价也不相同。
# 1≤N≤300

# 石子合并
# 倒着考虑:
# !分解石头成两个部分 代价为原来石头的重量 求分解成n个部分的最小代价

from functools import lru_cache
from itertools import accumulate
import sys

INF = int(1e20)
sys.setrecursionlimit(int(1e9))

n = int(input())
nums = list(map(int, input().split()))

preSum = [0] + list(accumulate(nums))


@lru_cache(maxsize=None)
def dfs(left: int, right: int) -> int:
    """[left,right]这一段合并的代价最小"""
    if left >= right:
        return 0

    res = INF
    for i in range(left, right):
        res = min(res, dfs(left, i) + dfs(i + 1, right) + preSum[right + 1] - preSum[left])
    return res


print(dfs(0, n - 1))
