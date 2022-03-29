from functools import lru_cache
from operator import add, mul, sub
import re
from typing import List

OPTIONS = {
    '+': add,
    '-': sub,
    '*': mul,
}


class Solution:
    def diffWaysToCompute(self, expression: str) -> List[int]:
        @lru_cache(None)
        def dfs(left: int, right: int) -> List[int]:
            if left + 1 >= right:
                return [int(arr[left])]

            res = []
            for mid in range(left, right):
                if arr[mid] not in OPTIONS:
                    continue
                opt = OPTIONS[arr[mid]]
                for leftRes in dfs(left, mid):
                    for rightRes in dfs(mid + 1, right):
                        cand = opt(leftRes, rightRes)
                        res.append(int(cand))
            return res

        arr = list(re.split(r'(\D)', expression))
        return dfs(0, len(arr) - 1)


print(Solution().diffWaysToCompute(expression="2*3-4*5"))
print(Solution().diffWaysToCompute(expression="11"))
