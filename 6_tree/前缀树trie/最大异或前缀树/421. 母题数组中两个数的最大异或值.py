from typing import List
from ATCBinaryTrie import BinaryTrie


class Solution:
    def findMaximumXOR(self, nums: List[int]) -> int:
        n = len(nums)
        max_log = max(nums).bit_length()
        bt = BinaryTrie(add_query_limit=n, max_log=max_log, allow_multiple_elements=False)
        for num in nums:
            bt.add(num)
        res = 0
        for num in nums:  # 查询最大异或值
            bt.xor_all(num)
            res = max(res, bt.maximum())
            bt.xor_all(num)
        return res


assert Solution().findMaximumXOR([3, 10, 5, 25, 2, 8]) == 28
