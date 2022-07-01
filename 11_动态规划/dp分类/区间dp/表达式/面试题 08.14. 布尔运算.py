from functools import lru_cache
from operator import and_, or_, xor
import re
from typing import List

# 给定一个布尔表达式和一个期望的布尔结果 result，布尔表达式由 0 (false)、1 (true)、& (AND)、 | (OR) 和 ^ (XOR) 符号组成。
# 你可以决定运算顺序(即添加括号)
# and or xor
OPT = {"&": and_, "|": or_, "^": xor}


class Solution:
    def countEval(self, s: str, result: int) -> int:
        """实现一个函数，算出有几种可使该表达式得出 result 值的括号方法。"""

        @lru_cache(None)
        def dfs(left: int, right: int) -> List[int]:
            """返回:计算出[0,1]的方法数"""
            if left > right:
                return [0, 0]
            if left == right:
                return [int(s[left] == "0"), int(s[left] == "1")]

            res = [0, 0]
            for i in range(left + 1, right, 2):
                opt = OPT[s[i]]
                for leftValue, leftCount in enumerate(dfs(left, i - 1)):
                    for rightValue, rightCount in enumerate(dfs(i + 1, right)):
                        cand = opt(leftValue, rightValue)
                        res[cand] += leftCount * rightCount
            return res

        arr = re.split(r"(\D)", s)
        return dfs(0, len(arr) - 1)[result]


print(Solution().countEval(s="1^0|0|1", result=1))
