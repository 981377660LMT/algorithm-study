# inspect.getsource(customfunction.__class__) 查看api

import inspect
from typing import List


class CustomFunction:
    # Returns f(x, y) for any given positive integers x and y.
    # Note that f(x, y) is increasing with respect to both x and y.
    # i.e. f(x, y) < f(x + 1, y), f(x, y) < f(x, y + 1)
    def f(self, x, y) -> int: ...


# 题目保证 f(x, y) == z 的解处于 1 <= x, y <= 1000 的范围内。
# 尽管函数的具体式子未知，但它是单调递增函数
# 请你计算方程 f(x,y) == z 所有可能的正整数 数对 x 和 y
#
# 总结：双指针
class Solution:
    def findSolution(self, customfunction: "CustomFunction", z: int) -> List[List[int]]:
        # print(inspect.getsource(customfunction.__class__))

        res = []
        x, y = 1, 1000
        while x <= 1000 and y:
            cur = customfunction.f(x, y)
            if cur < z:
                x += 1
            elif cur > z:
                y -= 1
            else:
                res.append([x, y])
                x += 1
                y -= 1

        return res
