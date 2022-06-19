from typing import List


class Solution:
    def calculateTax(self, brackets: List[List[int]], income: int) -> float:
        """表示第 i 个税级的上限是 upperi ，征收的税率为 percenti
        
        最后一个税级的上限大于等于 income
        """
        res = 0
        for i, (upper, precent) in enumerate(brackets):
            pre = brackets[i - 1][0] if i > 0 else 0
            if upper >= income:
                diff = income - pre
                res += diff * precent
                break
            diff = upper - pre
            res += diff * precent
        return res / 100  # !最后处理浮点数


print(Solution().calculateTax(brackets=[[3, 50], [7, 10], [12, 25]], income=10))

