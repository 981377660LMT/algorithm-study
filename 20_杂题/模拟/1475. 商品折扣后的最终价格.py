from typing import List


# !这种一般用 `while remain` 模拟 思路会比较好
class Solution:
    def calculateTax(self, brackets: List[List[int]], income: int) -> float:
        """
        brackets[i] 表示第 i 个税级的上限是 upperi ，征收的税率为 percenti
        最后一个税级的上限大于等于 income
        给你一个整数 income 表示你的总收入。返回你需要缴纳的税款总额
        """
        remain = income
        res, index = 0, 0
        while remain:
            pre = brackets[index - 1][0] if index > 0 else 0
            cur = min(remain, brackets[index][0] - pre)
            res += cur * brackets[index][1]
            remain -= cur
            index += 1

        return res / 100  # !最后处理浮点数


print(Solution().calculateTax(brackets=[[3, 50], [7, 10], [12, 25]], income=10))
