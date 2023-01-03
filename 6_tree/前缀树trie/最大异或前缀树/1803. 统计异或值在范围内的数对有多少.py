# 1803. 统计异或值在范围内的数对有多少
from typing import List
from ATCBinaryTrie import BinaryTrie


class Solution:
    def countPairs(self, nums: List[int], low: int, high: int) -> int:
        n = len(nums)
        max_log = max(nums).bit_length()
        bt = BinaryTrie(add_query_limit=n, max_log=max_log, allow_multiple_elements=True)
        for num in nums:
            bt.add(num)
        res = 0
        for num in nums:
            bt.xor_all(num)
            res += bt.bisect_right(high) - bt.bisect_left(low)
            bt.xor_all(num)
        return res // 2


assert Solution().countPairs(nums=[1, 4, 2, 7], low=2, high=6) == 6
