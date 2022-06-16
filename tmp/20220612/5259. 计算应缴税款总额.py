from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# !把浮点运算放在最后做。
class Solution:
    def calculateTax(self, brackets: List[List[int]], income: int) -> float:
        """注意要记录pre"""
        brackets.sort(key=lambda x: x[0])
        res = 0
        pre = 0
        for upper, rate in brackets:
            if upper >= income:
                res += (income - pre) * rate
                break
            res += (upper - pre) * rate
            pre = upper
        return res / 100


print(Solution().calculateTax(brackets=[[3, 50], [7, 10], [12, 25]], income=10))
