
class Solution(object):
    def subarrayBitwiseORs(self, A):
        """或运算单调不降 所以ndp里最多是32种不同的值
        
        O(nlogA)
        """
        res = set()
        dp = {0}
        for x in A:
            ndp = {x | y for y in dp} | {x}
            res |= ndp
            dp = ndp
        return len(res)

