from math import gcd
from typing import List

# 统计每个数字和k的最大公约数，只要最大公约数之乘积可以被k整除就行
class Solution:
    def coutPairs(self, nums: List[int], k: int) -> int:
        ss = set([k])
        prime = 2
        while prime * prime <= k:
            if k % prime == 0:
                ss.add(prime)
                ss.add(k // prime)
            prime += 1

        mapping = dict([(c, 0) for c in ss])
        res = 0
        for prime in range(len(nums)):
            x = gcd(k, nums[prime])
            y = k // x
            if y == 1:
                res += prime
            elif y in mapping:
                res += mapping[y]
            for num in ss:
                if nums[prime] % num == 0:
                    mapping[num] += 1
        return res
