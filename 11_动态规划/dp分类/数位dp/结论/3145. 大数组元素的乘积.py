# 3145. 大数组元素的乘积
# https://leetcode.cn/problems/find-products-of-elements-of-big-array/description/
# nums    0 | 1 | 2  3    | 4  5     6     7       | 8  9     10    11       ...
# bignums _ | 0 | 1  0  1 | 2  0  2  1  2  0  1  2 | 3  0  3  1  3  0  1  3  ...


from typing import List


class Solution:
    def findProductsOfElements(self, queries: List[List[int]]) -> List[int]:
        ...
