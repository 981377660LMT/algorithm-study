# https://atcoder.jp/contests/abc140/tasks/abc140_e

# 输入 n(≤1e5) 和一个 1~n 的排列 p。
# !输出 p 中所有长度至少为 2 的子数组的第二大元素的和。

# !贡献法，对每个 p[i]，求上上个、上个、下个、下下个更大元素的位置，分别记作 LL L  R RR。
# !那么 p[i] 对答案的贡献为 p[i] * ((L-LL)*(R-i) + (RR-R)*(i-L))。 (起点*终点个数)

from typing import List
from 对每个数寻找右侧第k个比自己大的数 import kthGreaterElement


def solve(nums: List[int]) -> int:
    """求nums所有子数组第二大元素的和"""
    n = len(nums)
    right1 = kthGreaterElement(nums, 1)  # !不存在为n
    right2 = kthGreaterElement(nums, 2)
    left1 = kthGreaterElement(nums[::-1], 1)[::-1]
    left1 = [n - 1 - pos for pos in left1]  # !不存在为-1
    left2 = kthGreaterElement(nums[::-1], 2)[::-1]
    left2 = [n - 1 - pos for pos in left2]
    res = 0
    for i in range(len(nums)):
        res += nums[i] * (
            (i - left1[i]) * (right2[i] - right1[i]) + (right1[i] - i) * (left1[i] - left2[i])
        )
    return res


n = int(input())
nums = list(map(int, input().split()))
print(solve(nums))
