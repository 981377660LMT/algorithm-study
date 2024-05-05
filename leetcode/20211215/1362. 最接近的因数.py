from math import ceil
from typing import List

# 给你一个整数 num，请你找出同时满足下面全部要求的两个整数：

# 两数乘积等于  num + 1 或 num + 2
# 以绝对差进行度量，两数大小最接近

# 从sqrt(num)递减找因数


class Solution:
    def closestDivisors(self, num: int) -> List[int]:
        # 比谁快
        for fac in range(ceil(num ** 0.5), 0, -1):
            if (num + 1) % fac == 0:
                return [fac, (num + 1) // fac]
            if (num + 2) % fac == 0:
                return [fac, (num + 2) // fac]

