# n<=1e5
# !最大不相交区间数

from typing import List


class Solution:
    def findLongestChain(self, pairs: List[List[int]]) -> int:
        """
        当且仅当 b < c 时，数对(c, d) 才可以跟在 (a, b) 后面
        给定一个数对集合，找出能够形成的最长数对链的长度。
        """
        pairs.sort(key=lambda x: x[1])
        res, preEnd = 0, -int(1e20)
        for start, end in pairs:
            if start > preEnd:
                res += 1
                preEnd = end
        return res
