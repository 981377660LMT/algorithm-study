from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def calculateTax(self, brackets: List[List[int]], income: int) -> float:
        """注意要记录pre"""
        brackets.sort(key=lambda x: x[0])
        res = 0
        pre = 0
        for upper, rate in brackets:
            if upper >= income:
                res += (income - pre) * rate / 100
                break
            res += (upper - pre) * rate / 100
            pre = upper
        return res


print(Solution().calculateTax(brackets=[[3, 50], [7, 10], [12, 25]], income=10))
