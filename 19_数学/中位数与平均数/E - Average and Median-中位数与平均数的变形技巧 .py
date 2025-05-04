# 给定一个长度为N(N ≤1e5)的数组a。
# 你需要从a中选出一些数构成数组b,需要满足数组a中相邻的两个数至少有一个被选上。
# !问构造出的b的平均数和中位数的最大值。
# 长度为x的数组的中位数为x/2向上取整。
# 输出平均值时，只要答案相差1e-3一下就算正确
# !2<=n<=1e5,1<=ai<=1e9


# 看到最大和最小就应该想到是二分答案
# !二分+dp。

# 有关数组平均数的技巧:
# !所有数减去平均数avg，此时如果数组的和>=0，说明其平均值>=avg
# 有关数组中位数的技巧:
# !nums[i]>=k则为1，否则为-1
# !此时如果数组的和>0，说明有`超过`一半的数都>=k，即中位数>=k

from typing import List
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
EPS = 1e-5

n = int(input())
nums = list(map(int, input().split()))


def calMaxSum(nums: List[int]) -> float:
    """给定一个数组，需要满足相邻的两项至少有一个选上，问和的最大值是多少"""
    n = len(nums)
    dp1, dp2 = nums[0], 0  # select jump
    for i in range(1, n):
        cur = nums[i]
        dp1, dp2 = max(dp1 + cur, dp2 + cur), dp1
    return max(dp1, dp2)


def check1(mid: float) -> bool:
    """平均数是否可以达到mid"""
    newNums = [num - mid for num in nums]
    return calMaxSum(newNums) >= 0  # type: ignore


def check2(mid: float) -> bool:
    """中位数是否可以达到mid"""
    newNums = [1 if num >= mid else -1 for num in nums]
    return calMaxSum(newNums) > 0  # !注意这里严格大于0


left, right = 0, int(1e9 + 10)
while left <= right:
    mid = (left + right) / 2  # float
    if check1(mid):
        left = mid + EPS
    else:
        right = mid - EPS
print(right)

left, right = 0, int(1e9 + 10)
while left <= right:
    mid = (left + right) // 2  # int
    if check2(mid):
        left = mid + 1
    else:
        right = mid - 1
print(right)
