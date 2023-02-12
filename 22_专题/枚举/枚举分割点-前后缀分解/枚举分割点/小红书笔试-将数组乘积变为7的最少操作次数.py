# 小红书-前后缀分解

# 给定n个整数a1, a2, a3 … an。每次操作可以选择其中一个数，
# 并将这个数加上1或者减去1。小红非常喜欢7这个数，
# 他想知道至少需要多少次操作可以使这n个数的乘积为7？
# n<=2e4 -1e9<=ai<=1e9
# !最后只能有1 -1 7 -7
# !两种:最大的数变成7 或者最小的数变成-7 其余变成1/-1
# dp也可以。dp[i][v]表示[0,i]区间的乘积为v的最小代价

from typing import List

INF = int(1e18)

n = int(input())
nums = list(map(int, input().split()))
if n == 1:
    print(abs(nums[0] - 7))
    exit(0)


def cal(arr: List[int], target: int) -> int:
    """
    乘积变为target(1/-1)需要的最小操作次数
    枚举哪些(哪几个前缀)元素变为-1
    前后缀和
    """
    n = len(arr)
    preSum = [0]  # 到-1的距离的前缀和
    sufSum = [0]  # 到1的距离的后缀和
    for i in range(n):
        preSum.append(preSum[-1] + abs(arr[i] - (-1)))
        sufSum.append(sufSum[-1] + abs(arr[~i] - 1))
    sufSum = sufSum[::-1]

    lower = 0 if target == 1 else 1
    res = INF
    for i in range(lower, n + 1, 2):  # 枚举-1的个数
        cand = preSum[i] + sufSum[i]
        res = min(res, cand)

    return res


nums.sort()
res1 = cal(nums[1:], -1) + abs(nums[0] + 7)
res2 = cal(nums[:-1], 1) + abs(nums[-1] - 7)
print(min(res1, res2))
