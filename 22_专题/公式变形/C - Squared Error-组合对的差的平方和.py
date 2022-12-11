# C - Squared Error
# !求所有组合对的差的平方和
# !n<=3e5 nums[i]<=200

from collections import defaultdict
from typing import List


def squaredError1(nums: List[int]) -> int:
    """O(n)"""
    n = len(nums)
    return n * sum(num * num for num in nums) - sum(nums) ** 2


def squaredError2(nums: List[int]) -> int:
    """O(n+k^2) 小值域搜答案/计算贡献"""
    counter = defaultdict(int)
    res = 0
    for num in nums:
        for pre, count in counter.items():
            res += (num - pre) * (num - pre) * count
        counter[num] += 1
    return res


n = int(input())
nums = list(map(int, input().split()))
# print(squaredError1(nums))
print(squaredError2(nums))
