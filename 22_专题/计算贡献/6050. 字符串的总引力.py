# 枚举贡献
# 2262. 字符串的总引力
# https://leetcode.cn/problems/total-appeal-of-a-string/
# 子数组不同元素个数之和.

from collections import defaultdict
from itertools import pairwise


class Solution:
    def appealSum(self, s: str) -> int:
        """所有子串的字符种类之和"""
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
            for pre, cur in pairwise(indexes):
                count = cur - pre
                res -= count * (count - 1) // 2

        return res


print(Solution().appealSum(s="abbca"))
print(Solution().appealSum(s="code"))
