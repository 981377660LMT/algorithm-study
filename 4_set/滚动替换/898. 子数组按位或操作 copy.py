class Solution(object):
    def subarrayBitwiseORs(self, A):
        res = set()
        cur = {0}
        for x in A:
            cur = {x | y for y in cur} | {x}
            res |= cur
        return len(res)

