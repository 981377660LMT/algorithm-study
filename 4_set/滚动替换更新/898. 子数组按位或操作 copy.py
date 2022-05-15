# 本质上是子序列dp
class Solution(object):
    def subarrayBitwiseORs(self, A):
        res = set()
        dp = {0}
        for x in A:
            ndp = {x | y for y in dp} | {x}
            res |= ndp
            dp = ndp
        return len(res)

