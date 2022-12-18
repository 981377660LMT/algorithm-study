# n<=2000 nums[i]<=10^9
# !每次移除数组中的两个数x,y 然后将gcd(x,y)或者min(x,y)加入数组中，问最后剩下的数有多少种可能

# !最后答案是若干个数的gcd 并且这个gcd需要小于等于所有数的最小值
# => 多个数的gcd 对每个数分解因数

from collections import defaultdict
from typing import List


def gcdOrMin(nums: List[int]) -> int:
    mp = defaultdict(list)
    for num in nums:
        for factor in getFactors(num):
            mp[factor].append(num)
    groupGcd = defaultdict(int)
    for factor, group in mp.items():
        g = tuple(group)
        groupGcd[g] = max(groupGcd[g], factor)
    min_ = min(nums)
    return sum(v <= min_ for v in groupGcd.values())


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


n = int(input())
nums = list(map(int, input().split()))
print(gcdOrMin(nums))
