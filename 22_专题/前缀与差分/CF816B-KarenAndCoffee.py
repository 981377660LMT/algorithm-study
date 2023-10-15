"""
# KarenAndCoffee
# 九条可怜和咖啡

给定n本食谱,每本食谱有两个整数l和r,表示这本食谱推荐[l, r]内的温度煮咖啡。
karen认为一个温度t是合适的,当且仅当存在>=threshold本食谱推荐t。
给定m个询问,每个询问有两个整数a和b,表示询问[a, b]内有多少个合适的温度。

所有数<=2e5

!等价于给定一个数组，查询区间[a, b]内有多少个数>=threshold.
!用前缀和查询.
"""


from itertools import accumulate
from typing import List, Tuple


def karenAndCoffee(
    recipes: List[Tuple[int, int]], queries: List[Tuple[int, int]], threshold: int
) -> List[int]:
    N = int(2e5 + 10)
    diff = [0] * N
    for left, right in recipes:
        diff[left] += 1
        diff[right + 1] -= 1
    nums = list(accumulate(diff))
    preSum = [0] * (len(nums) + 1)
    for i in range(len(nums)):
        preSum[i + 1] = preSum[i] + (nums[i] >= threshold)
    return [preSum[right + 1] - preSum[left] for left, right in queries]


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, k, q = map(int, input().split())
    recipes = [tuple(map(int, input().split())) for _ in range(n)]
    queries = [tuple(map(int, input().split())) for _ in range(q)]
    print(*karenAndCoffee(recipes, queries, k))
