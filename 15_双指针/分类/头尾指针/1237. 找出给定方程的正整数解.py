from typing import List


class CustomFunction:
    # Returns f(x, y) for any given positive integers x and y.
    # Note that f(x, y) is increasing with respect to both x and y.
    # i.e. f(x, y) < f(x + 1, y), f(x, y) < f(x, y + 1)
    def f(self, x, y) -> int:
        ...


# 题目保证 f(x, y) == z 的解处于 1 <= x, y <= 1000 的范围内。
# 尽管函数的具体式子未知，但它是单调递增函数
# 请你计算方程 f(x,y) == z 所有可能的正整数 数对 x 和 y

# 总结：双指针
class Solution:
    def findSolution(self, customfunction: "CustomFunction", z: int) -> List[List[int]]:
        res = []
        right = 1000
        for left in range(1, 1001):
            while right > 1 and customfunction.f(left, right) > z:
                right -= 1
            if customfunction.f(left, right) == z:
                res.append([left, right])
        return res
