import re
from typing import List
from functools import lru_cache
from operator import add, mul, sub, truediv

OPT = {"+": add, "-": sub, "*": mul, "/": truediv}

# 给你一个由数字和运算符组成的字符串 expression ，
# 按不同优先级组合数字和运算符，计算并返回所有可能组合的结果


class Solution:
    def diffWaysToCompute(self, expression: str) -> List[int]:
        """时间复杂度O(2^n)"""

        @lru_cache(None)
        def dfs(left: int, right: int) -> List[int]:
            if left > right:
                return []
            if left == right:
                return [int(arr[left])]

            res: List[int] = []
            for i in range(left + 1, right, 2):
                opt = OPT[arr[i]]
                for leftRes in dfs(left, i - 1):
                    for rightRes in dfs(i + 1, right):
                        cand = opt(leftRes, rightRes)
                        res.append(int(cand))
            return res

        arr = re.split(r"(\D)", expression)  # ()表示保留不匹配的内容
        return dfs(0, len(arr) - 1)


print(Solution().diffWaysToCompute(expression="2*3-4*5"))
print(Solution().diffWaysToCompute(expression="11"))
