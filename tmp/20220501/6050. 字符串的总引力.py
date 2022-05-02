from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)


# 枚举贡献
class Solution:
    def appealSum(self, s: str) -> int:
        n = len(s)
        indexMap = defaultdict(list)
        for i, char in enumerate(s):
            indexMap[char].append(i)

        res = 0
        # 我们枚举每个字符, 找出不含该字符的区间, 用M减去不含该字符的区间所能产生的子字符串就是单个字符对答案的贡献
        for indexes in indexMap.values():
            res += n * (n + 1) // 2  # 字符串子串数
            # 减去不含该字符的区间
            indexes = [-1] + indexes + [n]
            for pre, cur in zip(indexes, indexes[1:]):
                count = cur - pre
                res -= count * (count - 1) // 2

        return res


print(Solution().appealSum(s="abbca"))
print(Solution().appealSum(s="code"))
