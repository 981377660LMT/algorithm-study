# https://www.luogu.com.cn/problem/P6808
# 你需要修改序列中的一个数 P 为 Q，使得尽可能多的整数能够被表示出来。
# 如果有多种方案，则输出的 P 尽可能小。
# P 最小时如有多种方案，则输出的 Q 尽可能小。
# n <= 100, nums[i] <= 7000

# 求P:
# 记去掉每个数后,组成i的方案数为dp[i],可用撤销背包求出.
# !将这个数修改为一个很大的数后，可以表达数一定是2*dp[i]+1.那么只需要考虑最大的dp[i]对应的P.
# 求Q:
# !考虑不可行解的性质:去除P后，01背包里存在两个和sum1,sum2使得sum1+Q=sum2
# !等价于问:去除P后，所有元素及其相反数的任一子集能否表示出Q.
# !直接 01 背包即可求出这个最小的不能被表示的 Q，
# !s注意要先把所有正数加上去再加所有负数（否则可能某些方案在中间因为减出了负数而被舍弃)
# O(n*sum_)

from typing import List, Tuple
from Knapsack01Removable import Knapsack01Removable


MOD = int(1e9 + 7)  # 大素数


def candies(nums: List[int]) -> Tuple[int, int]:
    sum_ = sum(nums)
    dp = Knapsack01Removable(sum_, MOD)
    for v in nums:
        dp.add(v)

    maxCount, bestI = -1, -1
    for i, v in enumerate(nums):
        dp.remove(v)
        count = sum(dp.query(i) != 0 for i in range(1, sum_ + 1))
        if count > maxCount or (count == maxCount and v < nums[bestI]):
            maxCount = count
            bestI = i
        dp.add(v)

    dp.remove(nums[bestI])
    bitset = 1 << sum_  # 所有数加一个偏移sum_，避免出现负数
    for i, v in enumerate(nums):
        if i == bestI:
            continue
        bitset |= (bitset << v) | (bitset >> v)

    resQ = -1
    for i in range(1, sum_ + 2):
        if not (bitset & (1 << (i + sum_))):
            resQ = i
            break

    return nums[bestI], resQ


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline
    n = map(int, input().split())
    nums = list(map(int, input().split()))
    p, q = candies(nums)
    print(p, q)
