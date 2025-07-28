# 齿轮转动时，任意相邻的两个齿轮转动齿数相同，且方向相反。


from typing import List


class Solution:
    def spinGears(self, ratio: List[int], cnt: int, degree: int) -> List[int]:
        base = ratio[cnt] * degree
        res = []
        for i, r in enumerate(ratio):
            sign = -1 if (i - cnt) & 1 else 1
            res.append(sign * base // r)
        return res
