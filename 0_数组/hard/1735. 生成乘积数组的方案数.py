from typing import List


class Solution:
    def waysToFillArray(self, queries: List[List[int]]) -> List[int]:
        ...


print(Solution().waysToFillArray(queries=[[2, 6], [5, 1], [73, 660]]))
# 输出：[4,1,50734910]
# 解释：每个查询之间彼此独立。
# [2,6]：总共有 4 种方案得到长度为 2 且乘积为 6 的数组：[1,6]，[2,3]，[3,2]，[6,1]。
# [5,1]：总共有 1 种方案得到长度为 5 且乘积为 1 的数组：[1,1,1,1,1]。
# [73,660]：总共有 1050734917 种方案得到长度为 73 且乘积为 660 的数组。1050734917 对 109 + 7 取余得到 50734910 。


# todo...
