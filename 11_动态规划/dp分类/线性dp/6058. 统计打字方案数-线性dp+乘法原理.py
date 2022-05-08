from functools import lru_cache
from itertools import groupby

MOD = int(1e9 + 7)

MAPPING = {
    '2': 'abc',
    '3': 'def',
    '4': 'ghi',
    '5': 'jkl',
    '6': 'mno',
    '7': 'pqrs',
    '8': 'tuv',
    '9': 'wxyz',
}


@lru_cache(None)
def cal(maxStep: int, groupLength: int) -> int:
    @lru_cache(None)
    def dfs(index: int) -> int:
        if index >= groupLength:
            return int(index == groupLength)

        res = 0
        for select in range(1, maxStep + 1):
            res += dfs(index + select)
            res %= MOD
        return res

    res = dfs(0)
    return res


# 线性dp+乘法原理
class Solution:
    def countTexts(self, pressedKeys: str) -> int:
        groups = [[char, len(list(group))] for char, group in groupby(pressedKeys)]
        res = 1
        for char, length in groups:
            res *= cal(len(MAPPING[char]), length)
            res %= MOD
        return res


print(Solution().countTexts(pressedKeys="22233"))

