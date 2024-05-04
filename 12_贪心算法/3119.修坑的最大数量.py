# 3119. Maximum Number of Potholes That Can Be Fixed
# https://leetcode.cn/problems/maximum-number-of-potholes-that-can-be-fixed/description/
# 给定一个01数组road，其中0表示平整的路面，1表示有坑的路面，给定一个整数budget，表示修坑的预算.
# 每修一个坑的花费为len+1，其中len表示坑的长度，求最多能修多少个坑。

from itertools import groupby
from typing import Counter


class Solution:
    def maxPotholes(self, road: str, budget: int) -> int:
        groups = [(char, len(list(group))) for char, group in groupby(road)]
        onesCounter = Counter(v for k, v in groups if k == "x")
        res, remain = 0, budget
        for len_, count in sorted(onesCounter.items(), reverse=True):
            cost = len_ + 1
            # 这个长度可以全部修完
            if cost * count <= remain:
                res += count * len_
                remain -= cost * count
            else:
                div = remain // cost
                res += div * len_
                remain = remain - cost * div
                if remain >= 1:
                    res += remain - 1
                break

        return res


print(Solution().maxPotholes(road="..xxxxx", budget=4))  # 3
print(Solution().maxPotholes(road="x.xx", budget=3))  # 2
