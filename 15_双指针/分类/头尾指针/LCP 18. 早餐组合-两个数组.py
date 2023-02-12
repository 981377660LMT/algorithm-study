"""
staple 中记录了每种主食的价格，
drinks 中记录了每种饮料的价格。
!小扣的计划选择一份主食和一款饮料，且花费不超过 x 元。请返回小扣共有多少种购买方案。
注意：答案需要以 1e9 + 7 (1000000007) 为底取余
"""

from typing import List

MOD = int(1e9 + 7)


class Solution:
    def breakfastNumber(self, staple: List[int], drinks: List[int], x: int) -> int:
        staple.sort()
        drinks.sort()
        n1, n2 = len(staple), len(drinks)
        res = 0
        right = n2 - 1
        for i in range(n1):
            while right >= 0 and staple[i] + drinks[right] > x:
                right -= 1
            res += right + 1
            res %= MOD
        return res


assert Solution().breakfastNumber([10, 20, 5], [5, 5, 2], 15) == 6
