from collections import defaultdict
from itertools import pairwise

# Unique Characters of Every Substring

# 0 <= s.length <= 10^4
# 对每一个字符i，向前找到相同的字符j，向后找到相同的字符k。当前字符对最终结果的贡献是：（i-j）*(k-i)。
# 枚举start,end 统计贡献


class Solution:
    def uniqueLetterString(self, s: str) -> int:
        indexMap = defaultdict(list)
        for i, char in enumerate(s):
            indexMap[char].append(i)

        res, n = 0, len(s)
        for indexes in indexMap.values():
            indexes = [-1] + indexes + [n]
            res += sum((b - a) * (c - b) for a, b, c in zip(indexes, indexes[1:], indexes[2:]))

        return res % (int(1e9 + 7))


print(Solution().uniqueLetterString(s="ABC"))
# 输出: 10
# 解释: 所有可能的子串为："A","B","C","AB","BC" 和 "ABC"。
#      其中，每一个子串都由独特字符构成。
#      所以其长度总和为：1 + 1 + 1 + 2 + 2 + 3 = 10
