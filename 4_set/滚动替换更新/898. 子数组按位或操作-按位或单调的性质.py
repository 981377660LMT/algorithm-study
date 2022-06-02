class Solution(object):
    def subarrayBitwiseORs(self, A):
        """或运算单调不降 所以ndp里最多是32种不同的值
        
        O(nlogA)
        """
        res = set()
        dp = set()
        for cur in A:
            ndp = {cur | pre for pre in (dp | {0})}
            res |= ndp
            dp = ndp
        return len(res)

