from decimal import Decimal
from math import isqrt
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 ranks ，表示一些机械工的 能力值 。ranksi 是第 i 位机械工的能力值。能力值为 r 的机械工可以在 r * n2 分钟内修好 n 辆车。

# 同时给你一个整数 cars ，表示总共需要修理的汽车数目。

# 请你返回修理所有汽车 最少 需要多少时间。

# 注意：所有机械工可以同时修理汽车。
# 1 <= ranks.length <= 105
# 1 <= ranks[i] <= 100
# 1 <= cars <= 106

# !代码中对整数开方，只要整数转浮点没有丢失精度（在 2^53-1内），开方出来的整数部分就是正确的


class Solution:
    def repairCars(self, ranks: List[int], cars: int) -> int:
        def check(mid: int) -> bool:
            res = 0
            for rank in ranks:
                res += int((Decimal(mid) / Decimal(rank)).sqrt())
                if res >= cars:
                    return True
            return False

        left, right = 0, int(1e18)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
