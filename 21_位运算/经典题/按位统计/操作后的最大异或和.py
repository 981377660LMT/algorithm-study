from functools import reduce
from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumXOR(self, nums: List[int]) -> int:
        """按位统计 我们可以通过操作将每一位有1的位保留到最后"""
        return reduce(lambda x, y: x | y, nums)
        counter = [sum((num >> i) & 1 for num in nums) for i in range(32)]
        return sum(1 << i for i in range(32) if counter[i])


print(Solution().maximumXOR(nums=[1, 2, 3, 9, 2]))
print(Solution().maximumXOR(nums=[3, 2, 4, 6]))
