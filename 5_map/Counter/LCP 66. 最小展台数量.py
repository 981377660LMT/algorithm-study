from functools import reduce
from operator import or_
from typing import List
from collections import Counter


#  demand[i][j] 表示第 i 天展览时第 j 个展台的类型。
# 在满足每一天展台需求的基础上，请返回后勤部需要准备的 最小 展台数量。


class Solution:
    def minNumBooths(self, demand: List[str]) -> int:
        return reduce(or_, map(Counter, demand)).total()


print(Solution().minNumBooths(["A", "B", "C"]))


# !Counter的或运算操作表示取最大
print(Counter("A") | Counter("B"))
