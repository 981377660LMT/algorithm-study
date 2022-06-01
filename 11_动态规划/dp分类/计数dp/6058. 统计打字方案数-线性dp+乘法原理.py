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


# 预处理
MOD = int(1e9 + 7)
dp3 = [1, 1, 2, 4]
dp4 = [1, 1, 2, 4]
for _ in range(int(1e5 + 10)):
    dp3.append((dp3[-1] + dp3[-2] + dp3[-3]) % MOD)
    dp4.append((dp4[-1] + dp4[-2] + dp4[-3] + dp4[-4]) % MOD)


class Solution2:
    def countTexts(self, pressedKeys: str) -> int:
        groups = [[char, len(list(group))] for char, group in groupby(pressedKeys)]
        res = 1
        for char, length in groups:
            res = res * (dp4[length] if char in "79" else dp3[length]) % MOD
        return res

